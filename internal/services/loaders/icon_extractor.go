package loaders

import (
	"archive/zip"
	"encoding/base64"
	"io"
	"strings"
)

// IconExtractor is a shared utility for extracting icons from JARs
type IconExtractor interface {
	Extract(zipReader *zip.Reader, modID string) (string, error)
	ExtractFile(f *zip.File) (string, error)
	ExtractWithFallback(zipReader *zip.Reader, modID string, defaultIcon string) string
}

type iconExtractor struct{}

// NewIconExtractor creates a new iconExtractor
func NewIconExtractor() IconExtractor {
	return &iconExtractor{}
}

// Extract tries to find and extract an icon from the JAR
func (ie *iconExtractor) Extract(zipReader *zip.Reader, modID string) (string, error) {
	// Common icon paths to search
	iconPaths := []string{
		modID + ".png",
		"logo.png",
		"icon.png",
		"pack.png",
		"assets/" + modID + "/icon.png",
		"assets/" + modID + "/logo.png",
		"assets/" + modID + "/pack.png",
		"assets/" + modID + "/textures/logo.png",
	}

	for _, path := range iconPaths {
		for _, f := range zipReader.File {
			if strings.EqualFold(f.Name, path) {
				return ie.ExtractFile(f)
			}
		}
	}

	// No icon found, return empty string (will use default later)
	return "", nil
}

// ExtractFile extracts and base64 encodes a file
func (ie *iconExtractor) ExtractFile(f *zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	// Base64 encode the image data
	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:image/png;base64," + encoded, nil
}

// ExtractWithFallback extracts icon with a default fallback
func (ie *iconExtractor) ExtractWithFallback(zipReader *zip.Reader, modID string, defaultIcon string) string {
	icon, err := ie.Extract(zipReader, modID)
	if err != nil || icon == "" {
		return defaultIcon
	}
	return icon
}
