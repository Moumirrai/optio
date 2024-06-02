package videoManager

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/moumirrai/goffmpeg/transcoder"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	manager "optio/backend"
	"optio/backend/config"
	"optio/backend/stat"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

type SessionData struct {
	Files      map[string]VideoFileInfo `json:"videoFiles"`
	TotalSize  uint64                   `json:"totalVideoSize"`
	FilesOrder []string
}

// Manager handles collections of Files for conversion.
type VideoManager struct {
	Files         []VideoFileInfo
	Session       SessionData
	ctx           context.Context
	config        *config.Config
	stats         *stat.Stat
	manager       *manager.Manager
	transcoder    *transcoder.Transcoder
	stopRequested bool
}

// NewManager creates a new Manager.
func NewManager(c *config.Config, s *stat.Stat, m *manager.Manager) *VideoManager {
	return &VideoManager{
		config: c,
		stats:  s,
		Session: SessionData{
			Files:      make(map[string]VideoFileInfo),
			FilesOrder: make([]string, 0),
			TotalSize:  0,
		},
		manager:       m,
		stopRequested: false,
		transcoder:    new(transcoder.Transcoder),
	}
}

// WailsInit performs setup when Wails is ready.
func (vm *VideoManager) Startup(ctx context.Context) error {
	print("Manager startup")
	vm.ctx = ctx
	runtime.LogInfo(ctx, "Manager initialized...")
	return nil

}

func (vm *VideoManager) AddFiles() (string, error) {
	lastDir := vm.config.App.LastDir
	if _, err := os.Stat(lastDir); os.IsNotExist(err) {
		lastDir, err = os.UserHomeDir()
		if err != nil {
			runtime.LogError(vm.ctx, "Error getting user home dir")
		}
	}
	runtime.LogInfo(vm.ctx, lastDir)
	files, err := runtime.OpenMultipleFilesDialog(vm.ctx, runtime.OpenDialogOptions{
		Title:            "Pepe",
		DefaultDirectory: lastDir,
		Filters: []runtime.FileFilter{
			{DisplayName: "Videos", Pattern: "*.mp4;*.mov;*.avi;*.mkv;*.wmv;*.flv;*.webm"},
		},
	})
	runtime.EventsEmit(vm.ctx, "video:addingFiles")
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", nil
	}

	//get path of the first file without the filename, so i guess, split by \ and then join all but the last one
	vm.config.SetLastDir(filepath.Dir(files[0]))

	addedFilesCount := 0 // Initialize the counter

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	for _, file := range files {
		fileInfo, err := GetFileInfo(file, ctx)
		if err != nil {
			return "", err
		}
		// Check if file already exists in the map
		if _, exists := vm.Session.Files[file]; !exists {
			vm.Session.Files[file] = fileInfo
			addedFilesCount++ // Increment the counter
			vm.Session.TotalSize += uint64(fileInfo.Size)
			vm.Session.FilesOrder = append(vm.Session.FilesOrder, file)
		}
	}

	jsonString, err := json.Marshal(vm.Session)
	if err != nil {
		return "", err
	}

	//send "added x files" notification
	vm.manager.Notify(fmt.Sprintf("Added %d files", addedFilesCount), manager.Success)

	return string(jsonString), nil
}

// Clear removes the files in the Manager.
func (vm *VideoManager) Clear() {
	if len(vm.Session.Files) == 0 {
		return
	}
	vm.manager.Notify("Cleared files", manager.Info)
	vm.Session.Files = make(map[string]VideoFileInfo)
	vm.Session.FilesOrder = make([]string, 0)
	vm.Session.TotalSize = 0
	debug.FreeOSMemory()
}

func (vm *VideoManager) Debug(source string) (VideoFileInfo, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	return GetFileInfo(source, ctx)
}

func (vm *VideoManager) StartReencoding() {
	vm.stopRequested = false

	for _, oderedFile := range vm.Session.FilesOrder {
		file := vm.Session.Files[oderedFile]
		if vm.stopRequested {
			fmt.Println("Stop requested")
			break
		}

		_config := vm.config.App.VideoOpt

		now := time.Now()
		file.StartTimestamp = &now

		if _, err := os.Stat(file.Path); os.IsNotExist(err) {
			vm.manager.Notify(fmt.Sprintf("Source file does not exist: %s", file.Name), manager.Error)
			continue
		}

		filenameWithoutExt := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))
		dest := path.Join(vm.config.App.OutDir, vm.config.App.Prefix+filenameWithoutExt+vm.config.App.Suffix+".mp4")

		if _, err := os.Stat(dest); err == nil {
			vm.manager.Notify(fmt.Sprintf("File already exists: %s", dest), manager.Error)
			continue
		}

		err := vm.transcoder.Initialize(file.Path, dest)
		if err != nil {
			vm.manager.Notify("Error with encoder", manager.Error)
			panic(err)
		}
		vm.transcoder.MediaFile().SetVideoCodec(_config.Codec)
		if _config.PercentageMode {
			vm.transcoder.MediaFile().SetVideoBitRate(fmt.Sprintf("%dk", file.Bitrate*_config.Percentage/100000))
		} else {
			vm.transcoder.MediaFile().SetVideoBitRate(fmt.Sprintf("%dk", _config.Bitrate))
		}

		done := vm.transcoder.Run(true)
		progress := vm.transcoder.Output()

		for p := range progress {
			file.Progress = p.Progress
			//eta := vm.calculateEta(file)
			sinceStartFloat := float64(time.Since(*file.StartTimestamp).Milliseconds())
			runtime.EventsEmit(vm.ctx, "conversion:video:progress", map[string]interface{}{
				"id":       file.ID,
				"progress": p.Progress,
				"time":     time.Since(*file.StartTimestamp).Milliseconds(),
				"eta":      ((sinceStartFloat / p.Progress) * 100.0) - sinceStartFloat,
			})
		}

		<-done

		runtime.EventsEmit(vm.ctx, "conversion:video:file", map[string]interface{}{
			"id":   file.ID,
			"time": time.Since(*file.StartTimestamp).Milliseconds(),
		})
	}

	vm.stopRequested = false
	runtime.EventsEmit(vm.ctx, "conversion:video:complete", nil)
}

func (vm *VideoManager) StopReencoding() {
	if vm.transcoder != nil {
		vm.stopRequested = true
		vm.transcoder.Stop()
	}
}
