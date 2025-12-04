package loaders

import (
	"ModpackGraph/internal/models"
	"archive/zip"
)

// ModLoader is the interface for mod metadata extraction
type ModLoader interface {
	// CanHandle returns true if this loader can handle the given JAR
	CanHandle(zipReader *zip.Reader) bool

	// ExtractMetadata extracts metadata from the JAR
	ExtractMetadata(zipReader *zip.Reader, jarPath string) (*models.ModMetadata, error)
}
