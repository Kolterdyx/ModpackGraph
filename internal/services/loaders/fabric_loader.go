package loaders

import (
	"ModpackGraph/internal/assets"
	"ModpackGraph/internal/logger"
	"ModpackGraph/internal/models"
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
)

// FabricLoader handles Fabric mod metadata extraction
type FabricLoader struct {
	iconExtractor IconExtractor
}

// NewFabricLoader creates a new FabricLoader
func NewFabricLoader(
	iconExtractor IconExtractor,
) ModLoader {
	return &FabricLoader{
		iconExtractor: iconExtractor,
	}
}

// CanHandle checks if this is a Fabric mod
func (fl *FabricLoader) CanHandle(zipReader *zip.Reader) bool {
	logger.GetLogger().Debugf("Checking for Fabric mod...")
	for _, f := range zipReader.File {
		if f.Name == "fabric.mod.json" {
			return true
		}
	}
	return false
}

// ExtractMetadata extracts metadata from a Fabric mod
func (fl *FabricLoader) ExtractMetadata(zipReader *zip.Reader, jarPath string) (*models.ModMetadata, error) {
	// Find fabric.mod.json
	var fabricModJSON *zip.File
	for _, f := range zipReader.File {
		if f.Name == "fabric.mod.json" {
			fabricModJSON = f
			break
		}
	}

	if fabricModJSON == nil {
		return nil, fmt.Errorf("fabric.mod.json not found")
	}

	// Read and parse fabric.mod.json
	rc, err := fabricModJSON.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open fabric.mod.json: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read fabric.mod.json: %w", err)
	}

	var fabricMod FabricModJSONV1
	if err := json.Unmarshal(data, &fabricMod); err != nil {
		return nil, fmt.Errorf("failed to parse fabric.mod.json: %w", err)
	}

	// Create ModMetadata
	metadata := models.NewModMetadata()
	metadata.ID = fabricMod.ID
	metadata.Version = fabricMod.Version
	metadata.Name = fabricMod.Name
	metadata.Description = fabricMod.Description
	metadata.LoaderType = models.LoaderTypeFabric
	metadata.FilePath = jarPath
	metadata.MetadataJSON = string(data)

	// Extract authors
	if len(fabricMod.Authors) > 0 {
		metadata.Authors = make([]string, 0, len(fabricMod.Authors))
		for _, author := range fabricMod.Authors {
			if author.Name != "" {
				metadata.Authors = append(metadata.Authors, author.Name)
			}
		}
	}

	// Map environment
	metadata.Environment = mapFabricEnvironment(fabricMod.Environment)

	// Extract dependencies
	metadata.Dependencies = fl.extractDependencies(fabricMod, metadata.ID)

	// Extract icon
	iconPath := fabricMod.Icon
	if iconPath != "" {
		metadata.IconData = fl.extractIconFromPath(zipReader, iconPath)
	}
	if metadata.IconData == "" {
		metadata.IconData = fl.iconExtractor.ExtractWithFallback(zipReader, metadata.ID, assets.DefaultModIconData)
	}

	return metadata, nil
}

// extractDependencies extracts dependencies from Fabric mod metadata
func (fl *FabricLoader) extractDependencies(fabricMod FabricModJSONV1, modID string) []*models.Dependency {
	deps := make([]*models.Dependency, 0)

	// Required dependencies (depends)
	for depID, constraint := range fabricMod.Depends {
		if depID == "fabricloader" || depID == "fabric" || depID == "minecraft" || depID == "java" {
			// Skip loader and game dependencies
			continue
		}
		versionRange, _ := parseVersionConstraint(constraint)
		deps = append(deps, models.NewDependency(modID, depID, true, versionRange))
	}

	// Recommended dependencies
	for depID, constraint := range fabricMod.Recommends {
		if depID == "fabricloader" || depID == "fabric" || depID == "minecraft" || depID == "java" {
			continue
		}
		versionRange, _ := parseVersionConstraint(constraint)
		deps = append(deps, models.NewDependency(modID, depID, false, versionRange))
	}

	// Suggested dependencies (also optional)
	for depID, constraint := range fabricMod.Suggests {
		if depID == "fabricloader" || depID == "fabric" || depID == "minecraft" || depID == "java" {
			continue
		}
		versionRange, _ := parseVersionConstraint(constraint)
		deps = append(deps, models.NewDependency(modID, depID, false, versionRange))
	}

	return deps
}

// extractIconFromPath extracts icon from a specific path
func (fl *FabricLoader) extractIconFromPath(zipReader *zip.Reader, iconPath string) string {
	for _, f := range zipReader.File {
		if f.Name == iconPath {
			icon, err := fl.iconExtractor.ExtractFile(f)
			if err == nil {
				return icon
			}
		}
	}
	return ""
}

// mapFabricEnvironment maps Fabric environment to our model
func mapFabricEnvironment(env FabricV1Environment) models.Environment {
	switch env {
	case FabricV1EnvironmentClient:
		return models.EnvironmentClient
	case FabricV1EnvironmentServer:
		return models.EnvironmentServer
	case FabricV1EnvironmentUniversal:
		return models.EnvironmentBoth
	default:
		return models.EnvironmentBoth
	}
}

// parseVersionConstraint parses a version constraint string
func parseVersionConstraint(constraint interface{}) (*models.VersionRange, error) {
	var constraintStr string

	switch v := constraint.(type) {
	case string:
		constraintStr = v
	case []interface{}:
		// Handle array format (e.g., [">=1.0.0", "<2.0.0"])
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				constraintStr = s
			}
		}
	case map[string]interface{}:
		// Handle object format with version field
		if ver, ok := v["version"].(string); ok {
			constraintStr = ver
		}
	}

	if constraintStr == "" {
		constraintStr = "*"
	}

	return models.NewVersionRange(constraintStr)
}

type FabricBaseJSON struct {
	SchemaVersion int `json:"schemaVersion"`
}

// FabricModJSONV1 represents the structure of fabric.mod.json
type FabricModJSONV1 struct {
	FabricBaseJSON
	ID            string                            `json:"id"`
	Version       string                            `json:"version"`
	Provides      []string                          `json:"provides"`
	Environment   FabricV1Environment               `json:"environment"`
	Entrypoints   map[string][]FabricV1Entrypoint   `json:"entrypoints"`
	Jars          []FabricV1JarEntry                `json:"jars"`
	Mixins        []FabricV1MixinEntry              `json:"mixins"`
	AccessWidener string                            `json:"accessWidener"`
	Depends       map[string]FabricV1VersionMatcher `json:"depends"`
	Recommends    map[string]FabricV1VersionMatcher `json:"recommends"`
	Suggests      map[string]FabricV1VersionMatcher `json:"suggests"`
	Conflicts     map[string]FabricV1VersionMatcher `json:"conflicts"`
	Breaks        map[string]FabricV1VersionMatcher `json:"breaks"`
	Requires      []string                          `json:"requires"`
	Name          string                            `json:"name"`
	Description   string                            `json:"description"`
	Authors       []FabricV1Person                  `json:"authors"`
	Contributors  []FabricV1Person                  `json:"contributors"`
	Contact       map[string]string                 `json:"contact"`
	Icon          string                            `json:"icon"`
}
