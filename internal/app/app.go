package app

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	config Config
}

func NewApp(config Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	loadDefaultIconData()
}

type FileFilter struct {
	DisplayName string `json:"displayName"`
	Pattern     string `json:"pattern"`
}

type OpenDialogOptions struct {
	Title                      string       `json:"title,omitempty"`
	DefaultDirectory           string       `json:"defaultDirectory,omitempty"`
	DefaultFilename            string       `json:"defaultFilename,omitempty"`
	Filters                    []FileFilter `json:"filters,omitempty"`
	ShowHiddenFiles            bool         `json:"showHiddenFiles,omitempty"`
	CanCreateDirectories       bool         `json:"canCreateDirectories,omitempty"`
	ResolvesAliases            bool         `json:"resolvesAliases,omitempty"`
	TreatPackagesAsDirectories bool         `json:"treatPackagesAsDirectories,omitempty"`
}

func (a *App) OpenDirectoryDialog(options OpenDialogOptions) (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            options.Title,
		DefaultDirectory: options.DefaultDirectory,
	})
}

func (a *App) GenerateDependencyGraph(options GraphGenerationOptions) (*Graph, error) {
	return scanModFolder(options.Path)
}

func (a *App) Menu() *menu.Menu {
	m := menu.NewMenu()

	fileMenu := m.AddSubmenu("File")
	fileMenu.AddText("Quit", nil, func(_ *menu.CallbackData) {
		runtime.Quit(a.ctx)
	})

	aboutMenu := m.AddSubmenu("About")
	aboutMenu.AddText("ModpackGraph", nil, func(_ *menu.CallbackData) {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "About ModpackGraph",
			Message: fmt.Sprintf(`ModpackGraph
Version: %s
%s
GitHub: https://github.com/Kolterdyx/ModpackGraph
%s
License: MIT`,
				a.config.Info.Version,
				a.config.Info.Comments,
				a.config.Info.Copyright,
			),
		})
	})
	aboutMenu.AddText("Author", nil, func(_ *menu.CallbackData) {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Author",
			Message: fmt.Sprintf(
				`Developed by %s
GitHub: https://github.com/Kolterdyx`,
				a.config.Author.Name,
			),
		})
	})
	return m
}

//func (a *App) GenerateDependencyGraphSVG(modGraph *Graph) (string, error) {
//
//	content, err := modGraph.Graphviz(a.ctx)
//	if err != nil {
//		//log.Error("Failed to generate graphviz content: ", err)
//		return "", err
//	}
//	// remove everything before <svg
//	svgIndex := -1
//	for i := 0; i < len(content)-4; i++ {
//		if content[i:i+4] == "<svg" {
//			svgIndex = i
//			break
//		}
//	}
//	if svgIndex != -1 {
//		content = content[svgIndex:]
//	}
//	return content, nil
//}
