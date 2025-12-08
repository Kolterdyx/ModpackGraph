package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) OpenDirectoryDialog() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a directory",
	})
}
