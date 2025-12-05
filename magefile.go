//go:build mage
// +build mage

package main

import (
	"ModpackGraph/internal/enums"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:embed wails.json
var wailsJSON []byte

type WailsConfig struct {
	Info struct {
		ProductVersion string `json:"productVersion"`
	} `json:"info"`
}

func Dev() error {
	fmt.Println("Starting development server...")
	xflags, err := metaFlags(enums.BuildTypeDevelopment)
	if err != nil {
		return err
	}

	xflagsStr := "-X '" + strings.Join(xflags, "' -X '") + "'"
	ldflags := xflagsStr
	cmd := exec.Command("wails", "dev", "-tags", "webkit2_41", "-ldflags", ldflags, "-nogorebuild")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func metaFlags(buildType string) ([]string, error) {

	var config WailsConfig
	err := json.Unmarshal(wailsJSON, &config)
	if err != nil {
		return nil, err
	}

	return []string{
		"ModpackGraph/internal.Version=" + config.Info.ProductVersion,
		"ModpackGraph/internal.BuildType=" + buildType,
	}, nil
}

func BuildLinux() error {
	fmt.Println("Building...")
	xflags, err := metaFlags(enums.BuildTypeProduction)
	if err != nil {
		return err
	}

	ldflags := "-X '" + strings.Join(xflags, "' -X '") + "'"

	cmd := exec.Command("wails", "build", "-ldflags", ldflags, "-platform", "linux/amd64", "-tags", "webkit2_41")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func BuildWindows() error {
	fmt.Println("Building Windows...")
	xflags, err := metaFlags(enums.BuildTypeProduction)
	if err != nil {
		return err
	}
	ldflags := "-X '" + strings.Join(xflags, "' -X '") + "'"

	cmd := exec.Command("wails", "build", "-ldflags", ldflags, "-platform", "windows/amd64", "-nsis", "-webview2", "embed")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func BuildMac() error {
	fmt.Println("Building Mac...")
	xflags, err := metaFlags(enums.BuildTypeProduction)
	if err != nil {
		return err
	}
	ldflags := "-X '" + strings.Join(xflags, "' -X '") + "'"

	cmd := exec.Command("wails", "build", "-ldflags", ldflags, "-platform", "darwin/universal")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
