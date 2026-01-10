package loaders

import (
	"ModpackGraph/internal/logger"
	"ModpackGraph/internal/models"
	"archive/zip"
	"fmt"
	"go.uber.org/fx"
)

type loaderRegistryParams struct {
	fx.In
	Loaders []ModLoader `group:"mod_loader"`
}

// LoaderRegistry manages all mod loaders and detects the correct one for a JAR
type LoaderRegistry struct {
	loaders []ModLoader
}

// NewLoaderRegistry creates a new LoaderRegistry with all available loaders
func NewLoaderRegistry(p loaderRegistryParams) *LoaderRegistry {
	return &LoaderRegistry{
		loaders: p.Loaders,
	}
}

// DetectLoader finds the appropriate loader for a JAR
func (lr *LoaderRegistry) DetectLoader(zipReader *zip.Reader) (ModLoader, error) {
	logger.GetLogger().Debugf("Detecting mod loader, checking %d loaders...", len(lr.loaders))
	for _, loader := range lr.loaders {
		if loader.CanHandle(zipReader) {
			return loader, nil
		}
	}
	return nil, fmt.Errorf("no compatible loader found for JAR")
}

// ExtractMetadata detects the loader and extracts metadata
func (lr *LoaderRegistry) ExtractMetadata(zipReader *zip.Reader, jarPath string) (*models.ModMetadata, error) {
	loader, err := lr.DetectLoader(zipReader)
	if err != nil {
		return nil, err
	}

	return loader.ExtractMetadata(zipReader, jarPath)
}

// RegisterLoader adds a custom loader to the registry
func (lr *LoaderRegistry) RegisterLoader(loader ModLoader) {
	lr.loaders = append(lr.loaders, loader)
}
