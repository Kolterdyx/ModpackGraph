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

// getWebkitGTKVersion returns "4.0", "4.1", or an error
func getWebkitGTKVersion() (string, error) {
	// Try webkit2gtk-4.1 first
	if exists("webkit2gtk-4.1") {
		return "4.1", nil
	}
	// Then fallback to 4.0
	if exists("webkit2gtk-4.0") {
		return "4.0", nil
	}
	return "", fmt.Errorf("no supported webkit2gtk version found")
}

func exists(pkg string) bool {
	cmd := exec.Command("pkg-config", "--exists", pkg)
	return cmd.Run() == nil
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

func Dev() error {
	fmt.Println("Starting development server...")
	xflags, err := metaFlags(enums.BuildTypeDevelopment)
	if err != nil {
		return err
	}

	xflagsStr := "-X '" + strings.Join(xflags, "' -X '") + "'"
	ldflags := xflagsStr

	args := []string{
		"dev",
		"-ldflags", ldflags,
		"-nogorebuild",
	}

	if webkitVersion, err := getWebkitGTKVersion(); err == nil && webkitVersion == "4.1" {
		args = append(args, "-tags", "webkit2_41")
	} else if err != nil {
		return err
	}

	cmd := exec.Command("wails", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func BuildLinux() error {
	fmt.Println("Building...")
	xflags, err := metaFlags(enums.BuildTypeProduction)
	if err != nil {
		return err
	}

	ldflags := "-X '" + strings.Join(xflags, "' -X '") + "'"

	args := []string{
		"build",
		"-ldflags", ldflags,
		"-platform", "linux/amd64",
		"-nogorebuild",
	}

	if webkitVersion, err := getWebkitGTKVersion(); err == nil && webkitVersion == "4.1" {
		args = append(args, "-tags", "webkit2_41")
	} else if err != nil {
		return err
	}

	cmd := exec.Command("wails", args...)
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
