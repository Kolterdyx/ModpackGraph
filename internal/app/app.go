package app

import (
	"context"
	"encoding/json"
	"github.com/goccy/go-graphviz"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
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

type Layout graphviz.Layout

var AllLayouts = []Layout{
	Layout(graphviz.CIRCO),
	Layout(graphviz.DOT),
	Layout(graphviz.FDP),
	Layout(graphviz.NEATO),
	Layout(graphviz.NOP),
	Layout(graphviz.NOP1),
	Layout(graphviz.NOP2),
	Layout(graphviz.OSAGE),
	Layout(graphviz.PATCHWORK),
	Layout(graphviz.SFDP),
	Layout(graphviz.TWOPI),
}

func (l Layout) TSName() string {
	return string(l)
}

func (l Layout) Graphviz() graphviz.Layout {
	return graphviz.Layout(l)
}

type GraphOptions struct {
	Path   string `json:"path,omitempty"`
	Layout Layout `json:"layout,omitempty"`
}

func (a *App) GenerateDependencyGraphJSON(options GraphOptions) (string, error) {
	modData, err := scanModFolder(options.Path)
	if err != nil {
		return "", err
	}
	graphData, err := json.Marshal(modData)
	return string(graphData), err
}

func (a *App) GenerateDependencyGraphSVG(options GraphOptions) (string, error) {
	modData, err := scanModFolder(options.Path)
	if err != nil {
		return "", err
	}
	log.Debug("Scanned mod folder")
	svgData, err := generateDependencyGraphSVG(a.ctx, modData, options)
	if err != nil {
		return "", err
	}
	log.Debug("Generated SVG data")
	content := string(svgData)
	// remove everything before <svg
	svgIndex := -1
	for i := 0; i < len(content)-4; i++ {
		if content[i:i+4] == "<svg" {
			svgIndex = i
			break
		}
	}
	if svgIndex != -1 {
		content = content[svgIndex:]
	}
	return content, nil
}
