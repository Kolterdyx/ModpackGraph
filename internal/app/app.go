package app

import (
	"context"
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

func (a *App) OpenSelectFolderDialog(path, title string) (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: path,
	})
}

func (a *App) GenerateDependencyGraph(path string) (string, error) {
	modData, err := scanModFolder(path)
	if err != nil {
		return "", err
	}
	svgData, err := generateDependencyGraphSVG(a.ctx, modData)
	if err != nil {
		return "", err
	}
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
