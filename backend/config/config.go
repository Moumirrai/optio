package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"optio/backend/jpeg"
	"optio/backend/localstore"
	"optio/backend/png"
	video "optio/backend/videoEnc"
	"optio/backend/webp"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	goruntime "runtime"
	"strings"
)

const filename = "conf.json"

// App represents application persistent configuration values.
type App struct {
	OutDir               string         `json:"outDir"`
	Target               string         `json:"target"`
	Prefix               string         `json:"prefix"`
	Suffix               string         `json:"suffix"`
	Sizes                []*size        `json:"sizes"`
	ActiveSize           string         `json:"activeSize"`
	JpegOpt              *jpeg.Options  `json:"jpegOpt"`
	PngOpt               *png.Options   `json:"pngOpt"`
	WebpOpt              *webp.Options  `json:"webpOpt"`
	VideoOpt             *video.Options `json:"videoOpt"`
	LastDir              string         `json:"lastDir"`
	PreserveCreationTime bool           `json:"preserveCreationTime"`
	Advanced             *advanced      `json:"advanced"`
}

// Config represents the application settings.
type Config struct {
	App        *App
	ctx        context.Context
	localStore *localstore.LocalStore
}

type advanced struct {
	CoresUsed    int  `json:"coresUsed"`
	UseNvidiaGPU bool `json:"useNvidiaGPU"`
}

type size struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Strategy int    `json:"strategy"`
	Name     string `json:"name"`
}

// WailsInit performs setup when Wails is ready.
// Startup performs setup when Wails is ready.
func (c *Config) Startup(ctx context.Context) {
	runtime.LogInfo(ctx, "Config initialized...")
	c.ctx = ctx
}

// NewConfig returns a new instance of Config.
func NewConfig() *Config {
	c := &Config{}
	c.localStore = localstore.NewLocalStore()

	a, err := c.localStore.Load(filename)
	if err != nil {
		c.App, _ = defaults()
	}
	if err = json.Unmarshal(a, &c.App); err != nil {
		fmt.Printf("error")
	}
	return c
}

// GetAppConfig returns the application configuration.
func (c *Config) GetAppConfig() map[string]interface{} {
	return map[string]interface{}{
		"outDir":               c.App.OutDir,
		"target":               c.App.Target,
		"prefix":               c.App.Prefix,
		"suffix":               c.App.Suffix,
		"sizes":                c.App.Sizes,
		"activeSize":           c.App.ActiveSize,
		"jpegOpt":              c.App.JpegOpt,
		"pngOpt":               c.App.PngOpt,
		"webpOpt":              c.App.WebpOpt,
		"videoOpt":             c.App.VideoOpt,
		"lastDir":              c.App.LastDir,
		"preserveCreationTime": c.App.PreserveCreationTime,
	}
}

// OpenOutputDir opens the output directory using the native system browser.
func (c *Config) OpenOutputDir() {
	runtime.BrowserOpenURL(c.ctx, c.App.OutDir)
}

// RestoreDefaults sets the app configuration to defaults.
func (c *Config) RestoreDefaults() (err error) {
	var a *App
	a, err = defaults()
	if err != nil {
		return err
	}
	c.App = a
	if err = c.store(); err != nil {
		return err
	}
	return nil
}

// SetConfig sets and stores the given configuration.
func (c *Config) SetConfig(cfg string) error {
	a := &App{}
	if err := json.Unmarshal([]byte(cfg), &a); err != nil {
		runtime.LogErrorf(c.ctx, "failed to unmarshal config: %v", err)
		return err
	}
	c.App = a
	if err := c.store(); err != nil {
		runtime.LogErrorf(c.ctx, "failed to store config: %v", err)
		return err
	}
	return nil
}

// SetOutDir opens a directory select dialog and sets the output directory to
// the chosen directory.
func (c *Config) SetOutDir() string {
	//dir := c.Runtime.Dialog.SelectDirectory()
	dir, err := runtime.OpenDirectoryDialog(c.ctx, *&runtime.OpenDialogOptions{Title: "Select Output Directory"})
	if err != nil {
		runtime.LogErrorf(c.ctx, "failed to open directory dialog: %v", err)
		return ""
	}
	if dir != "" {
		c.App.OutDir = dir
		runtime.LogInfof(c.ctx, "set output directory: %s", dir)
		if err := c.store(); err != nil {
			runtime.LogErrorf(c.ctx, "failed to store config: %v", err)
		}
	}
	return c.App.OutDir
}

func (c *Config) SetLastDir(dir string) {
	c.App.LastDir = dir
	if err := c.store(); err != nil {
		runtime.LogErrorf(c.ctx, "failed to store config: %v", err)
	}
}

func (c *Config) GetCpuCores() int {
	return goruntime.NumCPU()
}

// defaults returns the application configuration defaults.
func defaults() (*App, error) {
	a := &App{
		Target:               "webp",
		JpegOpt:              &jpeg.Options{Quality: 80, PreserveMetadata: false},
		PngOpt:               &png.Options{Quality: 80},
		WebpOpt:              &webp.Options{Lossless: false, Quality: 80},
		PreserveCreationTime: false,
		Advanced: &advanced{
			CoresUsed:    int(float64(goruntime.NumCPU()) * 0.8),
			UseNvidiaGPU: CheckNVEncAvailable(),
		},
	}
	ud, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get user directory: %v", err)
		return nil, err
	}

	od := path.Join(ud, "optio")
	cp := filepath.Clean(od)

	if _, err = os.Stat(od); os.IsNotExist(err) {
		if err = os.Mkdir(od, 0777); err != nil {
			od = "./"
			fmt.Printf("failed to create default output directory: %v", err)
			return nil, err
		}
	}
	a.OutDir = cp
	a.LastDir = cp
	return a, nil
}

// store stores the configuration state to the file system.
func (c *Config) store() error {
	js, err := json.Marshal(c.GetAppConfig())
	if err != nil {
		runtime.LogErrorf(c.ctx, "failed to marshal config: %v", err)
		return err
	}
	if err = c.localStore.Store(js, filename); err != nil {
		runtime.LogErrorf(c.ctx, "failed to store config: %v", err)
		return err
	}
	return nil
}

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
