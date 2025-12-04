package loaders

import (
	"ModpackGraph/internal/models"
	"archive/zip"
)

// NeoForgeLoader handles NeoForge mod metadata extraction
// Currently uses the same format as modern Forge
type NeoForgeLoader struct {
	*ForgeModernLoader
}

// NewNeoForgeLoader creates a new NeoForgeLoader
func NewNeoForgeLoader(f *ForgeModernLoader) *NeoForgeLoader {
	return &NeoForgeLoader{
		ForgeModernLoader: f,
	}
}

// ExtractMetadata extracts metadata from a NeoForge mod
func (nfl *NeoForgeLoader) ExtractMetadata(zipReader *zip.Reader, jarPath string) (*models.ModMetadata, error) {
	metadata, err := nfl.ForgeModernLoader.ExtractMetadata(zipReader, jarPath)
	if err != nil {
		return nil, err
	}

	// Override loader type to NeoForge
	metadata.LoaderType = models.LoaderTypeNeoForge

	return metadata, nil
}
