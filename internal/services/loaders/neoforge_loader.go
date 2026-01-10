package loaders

import (
	"ModpackGraph/internal/models"
	"archive/zip"
)

// NeoForgeLoader handles NeoForge mod metadata extraction
// Currently uses the same format as modern Forge
type NeoForgeLoader struct {
	loader ModLoader
}

// NewNeoForgeLoader creates a new NeoForgeLoader
func NewNeoForgeLoader(extractor IconExtractor) ModLoader {
	return &NeoForgeLoader{
		loader: NewForgeModernLoader(extractor),
	}
}

func (nfl *NeoForgeLoader) CanHandle(zipReader *zip.Reader) bool {
	return nfl.loader.CanHandle(zipReader)
}

// ExtractMetadata extracts metadata from a NeoForge mod
func (nfl *NeoForgeLoader) ExtractMetadata(zipReader *zip.Reader, jarPath string) (*models.ModMetadata, error) {
	metadata, err := nfl.loader.ExtractMetadata(zipReader, jarPath)
	if err != nil {
		return nil, err
	}

	// Override loader type to NeoForge
	metadata.LoaderType = models.LoaderTypeNeoForge

	return metadata, nil
}
