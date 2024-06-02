package imageManager

import (
	"bufio"
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"image"
	"math"
	manager "optio/backend"
	"optio/backend/config"
	"optio/backend/image/jpeg"
	"optio/backend/image/png"
	"optio/backend/image/webp"
	"optio/backend/metadata"
	"optio/backend/stat"
	"os"
	"path"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"sync"
	"time"
)

type ImageFileInfo struct {
	Name          string    `json:"name"`
	ID            string    `json:"id"`
	Size          int64     `json:"size"`
	MimeType      string    `json:"type"`
	DateCreated   time.Time `json:"dateCreated"`
	Path          string    `json:"path"`
	ConvertedPath string    `json:"convertedPath"`
	Converted     bool      `json:"converted"`
	Error         string    `json:"error"`
	ConvertedSize int64     `json:"convertedSize"`
	Ratio         float64   `json:"ratio"`
}

type SessionData struct {
	FilesLookup map[string]bool
	Files       []ImageFileInfo `json:"imageFiles"`
	TotalSize   uint64          `json:"totalImageSize"`
}

// Manager handles collections of Files for conversion.
type ImageManager struct {
	Session       SessionData
	ctx           context.Context
	config        *config.Config
	manager       *manager.Manager
	stats         *stat.Stat
	stopRequested bool
}

// NewManager creates a new Manager.
func NewManager(c *config.Config, s *stat.Stat, m *manager.Manager) *ImageManager {
	return &ImageManager{
		config: c,
		stats:  s,
		Session: SessionData{
			Files:       make([]ImageFileInfo, 0),
			FilesLookup: make(map[string]bool),
			TotalSize:   0,
		},
		manager:       m,
		stopRequested: false,
	}
}

// WailsInit performs setup when Wails is ready.
func (im *ImageManager) Startup(ctx context.Context) error {
	print("Manager startup")
	im.ctx = ctx
	runtime.LogInfo(ctx, "Manager initialized...")
	return nil
}

func (im *ImageManager) Debug() bool {
	return false
}

func (im *ImageManager) AddFiles() (*SessionData, error) {
	lastDir := im.config.App.LastDir
	if _, err := os.Stat(lastDir); os.IsNotExist(err) {
		lastDir, err = os.UserHomeDir()
		if err != nil {
			runtime.LogError(im.ctx, "Error getting user home dir")
		}
	}
	files, err := runtime.OpenMultipleFilesDialog(im.ctx, runtime.OpenDialogOptions{
		Title:            "Pepe",
		DefaultDirectory: lastDir,
		Filters: []runtime.FileFilter{
			{DisplayName: "Images", Pattern: "*.jpg;*.jpeg;*.png;*.webp"},
		},
	})
	runtime.EventsEmit(im.ctx, "image:addingFiles")
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	//get path of the first file without the filename, so i guess, split by \ and then join all but the last one
	im.config.SetLastDir(filepath.Dir(files[0]))

	addedFilesCount := 0 // Initialize the counter

	for _, file := range files {
		if _, exists := im.Session.FilesLookup[file]; !exists {
			fileInfo, err := os.Stat(file)
			if err != nil {
				continue
			}
			im.addFileToSession(file, fileInfo)
			addedFilesCount++
		}
	}

	return &im.Session, nil
}

func (im *ImageManager) addFileToSession(file string, fileInfo os.FileInfo) {
	im.Session.TotalSize += uint64(fileInfo.Size())
	im.Session.Files = append(im.Session.Files, ImageFileInfo{
		Name:          fileInfo.Name(),
		ID:            fmt.Sprintf("%d%d%s", fileInfo.ModTime().Unix(), fileInfo.Size(), fileInfo.Name()),
		Size:          fileInfo.Size(),
		MimeType:      getFileType(filepath.Ext(file)),
		DateCreated:   metadata.GetCTime(fileInfo),
		Path:          file,
		ConvertedPath: "",
		Converted:     false,
		ConvertedSize: 0,
	})
	im.Session.FilesLookup[file] = true
}

// Clear removes the files in the Manager.
func (im *ImageManager) Clear() {
	if len(im.Session.Files) == 0 {
		return
	}
	im.manager.Notify("Cleared files", manager.Info)
	im.Session = SessionData{
		Files:       make([]ImageFileInfo, 0),
		FilesLookup: make(map[string]bool),
		TotalSize:   0,
	}
}

func (im *ImageManager) worker(fileInfo ImageFileInfo, done chan<- time.Duration) {
	start := time.Now()

	// Check if the source file exists
	if _, err := os.Stat(fileInfo.Path); os.IsNotExist(err) {
		im.manager.Notify(fmt.Sprintf("Source file does not exist: %s", fileInfo.Name), manager.Error)
		done <- 0
		return
	}

	// Open the image file
	file, err := os.Open(fileInfo.Path)
	if err != nil {
		// Handle error
		im.manager.Notify(fmt.Sprintf("Failed to open file: %s", fileInfo.Name), manager.Error)
		runtime.LogError(im.ctx, err.Error())
		done <- 0
		return
	}
	defer file.Close()

	// Create a buffered reader from the file
	reader := bufio.NewReader(file)

	var img image.Image

	switch filepath.Ext(fileInfo.Path) {
	case ".jpg", ".jpeg":
		img, err = jpeg.DecodeJPEG(reader)
	case ".png":
		img, err = png.DecodePNG(reader)
	case ".webp":
		img, err = webp.DecodeWebp(reader)
	}

	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to decode file: %s", fileInfo.Name), manager.Error)
		done <- 0
		return
	}

	if img == nil {
		im.manager.Notify(fmt.Sprintf("Decoded image is nil: %s", fileInfo.Name), manager.Error)
		done <- 0
		return
	}

	filenameWithoutExt := strings.TrimSuffix(fileInfo.Name, filepath.Ext(fileInfo.Name))
	dest := path.Join(im.config.App.OutDir, im.config.App.Prefix+filenameWithoutExt+im.config.App.Suffix+"."+im.config.App.Target)

	// Check if the destination file already exists
	if _, err := os.Stat(dest); err == nil {
		im.manager.Notify(fmt.Sprintf("File already exists: %s", dest), manager.Error)
		done <- 0
		return
	}

	destFile, err := os.Create(dest)

	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to create file: %s", dest), manager.Error)
		runtime.LogError(im.ctx, err.Error())
		done <- 0
		return
	}

	defer destFile.Close()

	writer := bufio.NewWriter(destFile)
	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to write to file: %s", dest), manager.Error)
		runtime.LogError(im.ctx, err.Error())
		done <- 0
		return
	}

	switch im.config.App.Target {
	case "jpg":
		err = jpeg.EncodeJPEG(img, writer, im.config.App.ImageOpt.JpegOpt)
	case "png":
		err = png.EncodePNG(img, writer, im.config.App.ImageOpt.PngOpt)
	case "webp":
		err = webp.EncodeWebp(img, writer, im.config.App.ImageOpt.WebpOpt)
	}
	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to encode file: %s", fileInfo.Name), manager.Error)
		runtime.LogError(im.ctx, err.Error())
		done <- 0
		return
	}

	err = writer.Flush()

	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to flush file: %s", fileInfo.Name), manager.Error)
		runtime.LogError(im.ctx, err.Error())
		done <- 0
		return
	}

	if im.config.App.PreserveCreationTime {

		if err := metadata.SetCTime(destFile, fileInfo.DateCreated); err != nil {
			im.manager.Notify(fmt.Sprintf("Failed to change file's date of creation: %s", dest), manager.Error)
			runtime.LogError(im.ctx, err.Error())
			done <- 0
			return
		}
	}

	stats, err := destFile.Stat()
	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to get file's stats: %s", dest), manager.Error)
		runtime.LogError(im.ctx, err.Error())
		done <- 0
		return
	}

	fileInfo.ConvertedPath = filepath.Clean(dest)
	fileInfo.Converted = true
	fileInfo.ConvertedSize = stats.Size()

	runtime.EventsEmit(im.ctx, "conversion:image:progress", map[string]interface{}{
		"id":      fileInfo.ID,
		"newSize": fileInfo.ConvertedSize,
		"ratio":   math.Round((1 - float64(fileInfo.ConvertedSize)/float64(fileInfo.Size)) * 100),
		"time":    time.Since(start).Milliseconds(),
	})

	// Send the conversion time to the channel
	done <- time.Since(start)
}

func decodeImage(filePath string, reader *bufio.Reader) (image.Image, error) {
	var img image.Image
	var err error
	switch filepath.Ext(filePath) {
	case ".jpg", ".jpeg":
		img, err = jpeg.DecodeJPEG(reader)
	case ".png":
		img, err = png.DecodePNG(reader)
	case ".webp":
		img, err = webp.DecodeWebp(reader)
	default:
		err = fmt.Errorf("unsupported file type: %s", filePath)
	}
	return img, err
}

func (im *ImageManager) encodeImage(writer *bufio.Writer, img image.Image) error {
	var err error
	switch im.config.App.Target {
	case "jpg":
		err = jpeg.EncodeJPEG(img, writer, im.config.App.ImageOpt.JpegOpt)
	case "png":
		err = png.EncodePNG(img, writer, im.config.App.ImageOpt.PngOpt)
	case "webp":
		err = webp.EncodeWebp(img, writer, im.config.App.ImageOpt.WebpOpt)
	}
	return err
}

/*func (im *ImageManager) worker(fileInfo ImageFileInfo, wg *sync.WaitGroup, done chan<- time.Duration) {
	defer wg.Done()
	start := time.Now()

	if _, err := os.Stat(fileInfo.Path); os.IsNotExist(err) {
		im.manager.Notify(fmt.Sprintf("Source file does not exist: %s", fileInfo.Name), manager.Error)
		done <- 0
		return
	}

	file, err := os.Open(fileInfo.Path)
	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to open file: %s", fileInfo.Name), manager.Error)
		done <- 0
		return
	}
	defer file.Close()

	// Create a buffered reader from the file
	reader := bufio.NewReader(file)
	img, err := decodeImage(fileInfo.Path, reader)

	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to decode file: %s", fileInfo.Name), manager.Error)
		done <- 0
		return
	}

	if img == nil {
		im.manager.Notify(fmt.Sprintf("Decoded image is nil: %s", fileInfo.Name), manager.Error)
		done <- 0
		return
	}

	filenameWithoutExt := strings.TrimSuffix(fileInfo.Name, filepath.Ext(fileInfo.Name))
	dest := path.Join(im.config.App.OutDir, im.config.App.Prefix+filenameWithoutExt+im.config.App.Suffix+"."+im.config.App.Target)

	if _, err := os.Stat(dest); err == nil {
		im.manager.Notify(fmt.Sprintf("File already exists: %s", dest), manager.Error)
		done <- 0
		return
	}

	destFile, err := os.Create(dest)

	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to create file: %s", dest), manager.Error)
		runtime.LogError(im.ctx, err.Error())
		done <- 0
		return
	}

	defer destFile.Close()

	writer := bufio.NewWriter(destFile)
	if err = im.encodeImage(writer, img); err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to encode file: %s", fileInfo.Name), manager.Error)
		done <- 0
		return
	}

	if err = writer.Flush(); err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to flush file: %s", fileInfo.Name), manager.Error)
		done <- 0
		return
	}

	if im.config.App.PreserveCreationTime {
		if err := metadata.SetCTime(destFile, fileInfo.DateCreated); err != nil {
			im.manager.Notify(fmt.Sprintf("Failed to change file's date of creation: %s", dest), manager.Error)
			runtime.LogError(im.ctx, err.Error())
			done <- 0
			return
		}
	}

	stats, err := destFile.Stat()
	if err != nil {
		im.manager.Notify(fmt.Sprintf("Failed to get file's stats: %s", dest), manager.Error)
		runtime.LogError(im.ctx, err.Error())
		done <- 0
		return
	}

	fileInfo.ConvertedPath = filepath.Clean(dest)
	fileInfo.Converted = true
	fileInfo.ConvertedSize = stats.Size()

	runtime.EventsEmit(im.ctx, "conversion:image:progress", map[string]interface{}{
		"id":      fileInfo.ID,
		"newSize": fileInfo.ConvertedSize,
		"ratio":   math.Round((1 - float64(fileInfo.ConvertedSize)/float64(fileInfo.Size)) * 100),
		"time":    time.Since(start).Milliseconds(),
	})

	// Send the conversion time to the channel
	done <- time.Since(start)
}*/

func (im *ImageManager) StopConversion() {
	im.stopRequested = true
	im.manager.Notify("Stopping conversion...", manager.Info)
}

func (im *ImageManager) StartConversion() {
	// Check if the destination directory exists
	if _, err := os.Stat(im.config.App.OutDir); os.IsNotExist(err) {
		im.manager.Notify("Destination directory does not exist, creating...", manager.Info)
		if err := os.MkdirAll(im.config.App.OutDir, 0755); err != nil {
			im.manager.Notify("Failed to create destination directory", manager.Error)
			runtime.LogError(im.ctx, err.Error())
			return
		}
	}

	// Limit the number of concurrent goroutines
	numCores := goruntime.NumCPU()
	maxGoroutines := int(float64(numCores) * 0.8)
	sem := make(chan struct{}, maxGoroutines)

	// Channel to collect conversion times
	done := make(chan time.Duration, len(im.Session.Files))

	var wg sync.WaitGroup
	converted := 0

	startTime := time.Now()

	for _, file := range im.Session.Files {
		if im.stopRequested {
			break
		}
		converted++
		sem <- struct{}{} // Acquire a token
		wg.Add(1)
		go func(fileInfo ImageFileInfo) {
			/*im.worker(fileInfo, &wg, done)*/
			im.worker(fileInfo, done)
			<-sem // Release the token
			wg.Done()
		}(file)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	// Collect and print the conversion times
	total := time.Duration(0)
	for t := range done {
		fmt.Println(t)
		total += t
	}

	timeElapsed := time.Since(startTime)

	im.stopRequested = false

	//im.manager.Notify(fmt.Sprintf("Converted %d files in %d", len(im.Session.Files), total.Milliseconds()), manager.Success)
	runtime.EventsEmit(im.ctx, "conversion:image:complete", map[string]interface{}{
		"length": converted,
		"time":   timeElapsed.Milliseconds(),
	})
	//add all times together, and send notification
}

func getFileType(t string) string {
	m, prs := mimes[t]
	if !prs {
		_ = fmt.Errorf("unsupported file type:" + t)
	}
	return m
}

var mimes = map[string]string{
	"image/.jpg": "jpg",
	"image/jpg":  "jpg",
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}
