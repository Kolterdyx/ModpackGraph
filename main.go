package main

import (
	app2 "ModpackGraph/internal/app"
	"bytes"
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

const DefaultUserLang = "en"

var supportedLangs = []string{"en", "es"}

//go:embed all:frontend/dist/frontend/browser
var assets embed.FS

//go:embed language-index.gohtml
var languageIndexTMPL string

//go:embed wails.json
var configJSON []byte

type assetFS struct {
	root string
}

func NewAssetFS(root string) fs.FS {
	return &assetFS{
		root: root,
	}
}

func init() {
	println("hello world")
}

func (c assetFS) Open(name string) (fs.File, error) {
	return assets.Open(path.Join(c.root, name))
}

func main() {

	// setup logging to file
	logFile, _ := os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()
	//log.SetOutput(logFile)
	// Create an instance of the app structure
	defer func() {
		if r := recover(); r != nil {
			//log.Fatal(r)
		}
	}()

	var config app2.Config
	err := json.Unmarshal(configJSON, &config)
	if err != nil {
		//log.Fatalf("Failed to parse wails.json: %v", err)
	}

	app := app2.NewApp(config)

	assetServer := http.FileServer(http.FS(NewAssetFS("frontend/dist/frontend/browser")))

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "ModpackGraph",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: nil,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/" {
					language := r.Header.Get("Accept-Language")
					//log.Debugf("Detected language: %s", language)
					selectedLang := DefaultUserLang
					for _, lang := range supportedLangs {
						if len(language) >= 2 && language[0:2] == lang {
							selectedLang = lang
							break
						}
					}
					//log.Debugf("Selected language: %s", selectedLang)
					tmpl, err := template.New("language-index").Parse(languageIndexTMPL)
					if err != nil {
						//log.Error("Failed to parse language index template: ", err)
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
						//log.Error("Failed to execute language index template: ", err)
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
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		OnStartup:        app.Startup,
		Menu:             app.Menu(),
		Bind: []any{
			app,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: true,
		},
		Mac: &mac.Options{
			WindowIsTranslucent: true,
		},
	})

	if err != nil {
		//log.Fatalf("Error: %v", err.Error())
	}
}
