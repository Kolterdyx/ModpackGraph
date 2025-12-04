package assets

import (
	"embed"
	"encoding/base64"
	"strings"
)

//go:embed pack.png
var defaultIconFS embed.FS

var DefaultModIconData string

func init() {
	data, err := defaultIconFS.ReadFile("pack.png")
	if err != nil {
		panic("failed to load default mod icon: " + err.Error())
	}
	var buf strings.Builder
	_, err = base64.NewEncoder(base64.StdEncoding, &buf).Write(data)
	if err != nil {
		panic("failed to base64 encode default mod icon: " + err.Error())
	}
	DefaultModIconData = "data:image/png;base64," + buf.String()
}
