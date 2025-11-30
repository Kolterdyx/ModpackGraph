package main

import (
	app2 "ModpackGraph/internal/app"
	"bytes"
	"embed"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
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

type assetFS struct {
	root string
}

func NewAssetFS(root string) fs.FS {
	return &assetFS{
		root: root,
	}
}

func (c assetFS) Open(name string) (fs.File, error) {
	return assets.Open(path.Join(c.root, name))
}

func main() {
	log.SetLevel(log.DebugLevel)
	// Create an instance of the app structure
	app := app2.NewApp()

	assetServer := http.FileServer(http.FS(NewAssetFS("frontend/dist/frontend/browser")))

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ModpackGraph",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: nil,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/" {
					language := r.Header.Get("Accept-Language")
					log.Debugf("Detected language: %s", language)
					selectedLang := DefaultUserLang
					for _, lang := range supportedLangs {
						if len(language) >= 2 && language[0:2] == lang {
							selectedLang = lang
							break
						}
					}
					log.Debugf("Selected language: %s", selectedLang)
					tmpl, err := template.New("language-index").Parse(languageIndexTMPL)
					if err != nil {
						log.Error("Failed to parse language index template: ", err)
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
						log.Error("Failed to execute language index template: ", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					_, _ = w.Write(buf.Bytes())
				} else {
					assetServer.ServeHTTP(w, r)
				}
			}),
		},
		EnableDefaultContextMenu: true,
		BackgroundColour:         &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:                app.Startup,
		Bind: []any{
			app,
		},
		EnumBind: []any{
			app2.AllLayouts,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
