package main

import (
	"context"
	"embed"
	manager "optio/backend"
	"optio/backend/config"
	imageManager "optio/backend/image"
	"optio/backend/stat"
	videoManager "optio/backend/video"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	cfg := config.NewConfig()
	m := manager.NewManager(cfg, stat.NewStat())
	st := stat.NewStat()
	im := imageManager.NewManager(cfg, st, m)
	vm := videoManager.NewManager(cfg, st, m)

	onStartup := func(ctx context.Context) {
		cfg.Startup(ctx)
		st.Startup(ctx)
		m.Startup(ctx)
		im.Startup(ctx)
		vm.Startup(ctx)
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:                    "Optio",
		Width:                    1200,
		Height:                   800,
		MinHeight:                600,
		MinWidth:                 1200,
		EnableDefaultContextMenu: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 33, G: 33, B: 33, A: 1},
		OnStartup:        onStartup,
		Bind: []interface{}{
			cfg,
			st,
			m,
			im,
			vm,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
