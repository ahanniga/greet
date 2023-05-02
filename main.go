package main

import (
	"embed"
	"github.com/fstanis/screenresolution"
	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var app *App

func main() {
	// Create an instance of the app structure
	app = NewApp()
	res := screenresolution.GetPrimary()
	width := int(float64(res.Width) * 0.7)
	height := int(float64(res.Height) * 0.9)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Greet",
		Width:  width,
		Height: height,
		//Menu:              AppMenu,
		MinWidth:          1024, // Buttons start overlapping under this
		MinHeight:         480,
		DisableResize:     false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		LogLevel:         logger.DEBUG,
		OnStartup:        app.startup,
		OnDomReady:       app.OnDomReady,
		OnBeforeClose:    app.OnBeforeClose,
		OnShutdown:       app.OnShutdown,
		Bind: []interface{}{
			app,
		},
		//Debug: options.Debug{
		//	OpenInspectorOnStartup: true,
		//},
	})

	if err != nil {
		log.Err(err)
	}
}
