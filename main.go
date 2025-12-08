package main

import (
	"ModpackGraph/internal"
	"ModpackGraph/internal/app"
	"ModpackGraph/internal/di"
	"ModpackGraph/internal/enums"
	"ModpackGraph/internal/logger"
	"ModpackGraph/internal/models"
	"ModpackGraph/internal/util"
	"bytes"
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"go.uber.org/fx"
	"io/fs"
	"net/http"
	"path"
	"text/template"
)

const DefaultUserLang = "en"

var supportedLangs = []string{"en", "es"}

//go:embed all:frontend/dist/frontend/browser
var assets embed.FS

//go:embed language-index.gohtml
var languageIndexTMPL string

type offsetFS struct {
	root string
}

func NewAssetFS(root string) fs.FS {
	return &offsetFS{
		root: root,
	}
}

func (c offsetFS) Open(name string) (fs.File, error) {
	return assets.Open(path.Join(c.root, name))
}

func main() {
	// Initialize logger
	logger.Init()
	logger.GetLogger().Infof("Starting ModpackGraph application %v", internal.BuildType)

	var application *app.App

	// Create FX app to wire dependencies
	fxApp := fx.New(
		fx.NopLogger,
		di.Module,
		fx.Provide(app.NewApp),
		fx.Populate(&application),
	)

	err := wails.Run(&options.App{
		Title:            "ModpackGraph",
		Width:            1024,
		Height:           768,
		AssetServer:      getAssetServerOptions(),
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		Bind:             []any{application},
		EnumBind: []any{
			models.AllLoaderTypes,
			models.AllEnvironments,
			models.AllConflictTypes,
			models.AllConflictSeverities,
			models.AllFeatureTypes,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: true,
		},
		Mac: &mac.Options{
			WindowIsTranslucent: true,
		},
		OnStartup: func(ctx context.Context) {
			go func() {
				if err := fxApp.Start(ctx); err != nil {
					logger.GetLogger().WithError(err).Fatal("Failed to start FX application")
				}
			}()
			application.Startup(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			return application.OnBeforeClose(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			if err := fxApp.Stop(ctx); err != nil {
				logger.GetLogger().WithError(err).Fatal("Failed to stop FX application")
			}
		},
	})

	if err != nil {
		logger.GetLogger().WithError(err).Fatal("Failed to start application")
	}
}

func getAssetServerOptions() *assetserver.Options {
	assetFS := NewAssetFS("frontend/dist/frontend/browser")
	assetServer := http.FileServer(http.FS(assetFS))

	return &assetserver.Options{
		Assets: util.If(internal.BuildType == enums.BuildTypeProduction, nil, assetFS),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				language := r.Header.Get("Accept-Language")
				selectedLang := DefaultUserLang
				for _, lang := range supportedLangs {
					if len(language) >= 2 && language[0:2] == lang {
						selectedLang = lang
						break
					}
				}
				tmpl, err := template.New("language-index").Parse(languageIndexTMPL)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				data := struct {
					Lang string
				}{
					Lang: selectedLang,
				}
				var buf bytes.Buffer
				err = tmpl.Execute(&buf, data)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write(buf.Bytes())
			} else {
				assetServer.ServeHTTP(w, r)
			}
		}),
	}
}
