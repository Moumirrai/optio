package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"optio/backend/image/jpeg"
	"optio/backend/image/png"
	"optio/backend/image/webp"
	"optio/backend/localstore"
	"optio/backend/models"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	goruntime "runtime"
	"strings"
)

const filename = "conf.json"

type App struct {
	OutDir               string               `json:"outDir"`
	Target               string               `json:"target"`
	Prefix               string               `json:"prefix"`
	Suffix               string               `json:"suffix"`
	Sizes                []*size              `json:"sizes"`
	ActiveSize           string               `json:"activeSize"`
	ImageOpt             *models.ImageOptions `json:"imageOpt"`
	VideoOpt             *models.VideoOptions `json:"videoOpt"`
	LastDir              string               `json:"lastDir"`
	PreserveCreationTime bool                 `json:"preserveCreationTime"`
	Advanced             *advanced            `json:"advanced"`
}

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

func (c *Config) GetAppConfig() map[string]interface{} {
	return map[string]interface{}{
		"outDir":               c.App.OutDir,
		"target":               c.App.Target,
		"prefix":               c.App.Prefix,
		"suffix":               c.App.Suffix,
		"sizes":                c.App.Sizes,
		"activeSize":           c.App.ActiveSize,
		"imageOpt":             c.App.ImageOpt,
		"videoOpt":             c.App.VideoOpt,
		"lastDir":              c.App.LastDir,
		"preserveCreationTime": c.App.PreserveCreationTime,
	}
}

func (c *Config) OpenOutputDir() {
	runtime.BrowserOpenURL(c.ctx, c.App.OutDir)
}

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

func (c *Config) SetOutDir() string {
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

func defaults() (*App, error) {
	a := &App{
		Target: "webp",
		ImageOpt: &models.ImageOptions{
			JpegOpt: &jpeg.Options{Quality: 80, PreserveMetadata: false},
			PngOpt:  &png.Options{Quality: 80},
			WebpOpt: &webp.Options{Lossless: false, Quality: 80},
		},
		VideoOpt: &models.VideoOptions{
			Width:          1920,
			Height:         1080,
			Bitrate:        50000,
			Codec:          "libx264",
			PercentageMode: true,
			Percentage:     50,
		},
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

/*func CheckFFmpegAvailable() bool {
	cmd := exec.Command("ffmpeg")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "ffmpeg")
}
*/
