package videoManager

import (
	//video "optio/backend/videoEnc"
	"os"
	"os/exec"
	"path"
	"strings"

	//"bufio"
	//"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime/debug"

	//"encoding/json"
	//"fmt"
	//"image"
	//"sync"

	//"io"
	manager "optio/backend"
	"optio/backend/config"

	//"optio/backend/jpeg"
	//"optio/backend/metadata"
	//"optio/backend/png"
	"optio/backend/stat"
	//"optio/backend/video"
	//"optio/backend/webp"
	//"os"
	//"path"
	//"runtime/debug"
	//"strings"

	"time"

	//"path/filepath"
	//goruntime "runtime"

	//"github.com/dsoprea/go-exif/v3"
	//exifcommon "github.com/dsoprea/go-exif/v3/common"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xfrr/goffmpeg/transcoder"
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
	files, err := runtime.OpenMultipleFilesDialog(vm.ctx, runtime.OpenDialogOptions{
		Title:            "Pepe",
		DefaultDirectory: vm.config.App.LastDir,
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

// TODO: DELETE
func CheckNVEncAvailable() bool {
	cmd := exec.Command("nvidia-smi")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "NVIDIA-SMI")
}

func CheckFFmpegAvailable() bool {
	cmd := exec.Command("ffmpeg")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "ffmpeg")
}

type Options struct {
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Bitrate int    `json:"bitrate"`
	Codec   string `json:"codec"`
}

func (vm *VideoManager) StartReencoding() {
	vm.stopRequested = false

	for _, oderedFile := range vm.Session.FilesOrder {
		file := vm.Session.Files[oderedFile]
		if vm.stopRequested {
			fmt.Println("Stop requested")
			break
		}
		start := time.Now()

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

		err := vm.transcoder.Initialize(file.Path, dest) // Replace "output.mp4" with your desired output file path
		if err != nil {
			vm.manager.Notify("Error with encoder", manager.Error)
			panic(err)
		}
		vm.transcoder.MediaFile().SetVideoCodec("h264_nvenc")
		vm.transcoder.MediaFile().SetVideoBitRate(fmt.Sprintf("%dk", 2500))

		done := vm.transcoder.Run(true)
		progress := vm.transcoder.Output()

		// Handle progress in the main goroutine
		for p := range progress {
			runtime.EventsEmit(vm.ctx, "conversion:video:progress", map[string]interface{}{
				"id":       file.ID,
				"progress": p.Progress,
			})
		}

		// Wait for the transcoder to finish
		<-done

		runtime.EventsEmit(vm.ctx, "conversion:video:file", map[string]interface{}{
			"id":   file.ID,
			"time": time.Since(start).Milliseconds(),
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
