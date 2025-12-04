package main

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
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

	assetServer := http.FileServer(http.FS(NewAssetFS("frontend/dist/frontend/browser")))

	_ = wails.Run(&options.App{
		Title:  "ModpackGraph",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: nil,
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
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		Bind:             []any{},
		Linux: &linux.Options{
			WindowIsTranslucent: true,
		},
		Mac: &mac.Options{
			WindowIsTranslucent: true,
		},
	})
}
