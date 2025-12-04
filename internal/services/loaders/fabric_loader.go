package loaders

import (
	"ModpackGraph/internal/assets"
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
) *FabricLoader {
	return &FabricLoader{
		iconExtractor: iconExtractor,
	}
}

// CanHandle checks if this is a Fabric mod
func (fl *FabricLoader) CanHandle(zipReader *zip.Reader) bool {
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

	var fabricMod FabricModJSON
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
func (fl *FabricLoader) extractDependencies(fabricMod FabricModJSON, modID string) []*models.Dependency {
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
func mapFabricEnvironment(env string) models.Environment {
	switch env {
	case "client":
		return models.EnvironmentClient
	case "server":
		return models.EnvironmentServer
	case "*", "":
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

// FabricModJSON represents the structure of fabric.mod.json
type FabricModJSON struct {
	SchemaVersion int                    `json:"schemaVersion"`
	ID            string                 `json:"id"`
	Version       string                 `json:"version"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Authors       []FabricAuthor         `json:"authors"`
	Contact       map[string]string      `json:"contact"`
	License       interface{}            `json:"license"` // Can be string or array
	Icon          string                 `json:"icon"`
	Environment   string                 `json:"environment"`
	Depends       map[string]interface{} `json:"depends"`
	Recommends    map[string]interface{} `json:"recommends"`
	Suggests      map[string]interface{} `json:"suggests"`
	Breaks        map[string]interface{} `json:"breaks"`
	Conflicts     map[string]interface{} `json:"conflicts"`
}

// FabricAuthor represents an author entry
type FabricAuthor struct {
	Name    string            `json:"name"`
	Contact map[string]string `json:"contact"`
}
