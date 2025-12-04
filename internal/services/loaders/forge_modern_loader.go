package loaders

import (
	"ModpackGraph/internal/assets"
	"ModpackGraph/internal/models"
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// ForgeModernLoader handles Forge 1.13+ mod metadata extraction
type ForgeModernLoader struct {
	iconExtractor IconExtractor
}

// NewForgeModernLoader creates a new ForgeModernLoader
func NewForgeModernLoader(
	iconExtractor IconExtractor,
) *ForgeModernLoader {
	return &ForgeModernLoader{
		iconExtractor: iconExtractor,
	}
}

// CanHandle checks if this is a modern Forge mod
func (fml *ForgeModernLoader) CanHandle(zipReader *zip.Reader) bool {
	for _, f := range zipReader.File {
		if f.Name == "META-INF/mods.toml" {
			return true
		}
	}
	return false
}

// ExtractMetadata extracts metadata from a modern Forge mod
func (fml *ForgeModernLoader) ExtractMetadata(zipReader *zip.Reader, jarPath string) (*models.ModMetadata, error) {
	// Find mods.toml
	var modsToml *zip.File
	for _, f := range zipReader.File {
		if f.Name == "META-INF/mods.toml" {
			modsToml = f
			break
		}
	}

	if modsToml == nil {
		return nil, fmt.Errorf("META-INF/mods.toml not found")
	}

	// Read and parse mods.toml
	rc, err := modsToml.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open mods.toml: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read mods.toml: %w", err)
	}

	var forgeMod ForgeModsToml
	if err := toml.Unmarshal(data, &forgeMod); err != nil {
		return nil, fmt.Errorf("failed to parse mods.toml: %w", err)
	}

	// Get the first mod entry (most mods.toml files have one mod)
	if len(forgeMod.Mods) == 0 {
		return nil, fmt.Errorf("no mods defined in mods.toml")
	}

	mod := forgeMod.Mods[0]

	// Handle version placeholders
	version := mod.Version
	if strings.Contains(version, "${") {
		version = fml.resolveVersionPlaceholder(zipReader, version)
	}

	// Create ModMetadata
	metadata := models.NewModMetadata()
	metadata.ID = mod.ModID
	metadata.Version = version
	metadata.Name = mod.DisplayName
	metadata.Description = mod.Description
	metadata.LoaderType = models.LoaderTypeForgeModern
	metadata.Environment = models.EnvironmentBoth // Forge mods default to both
	metadata.FilePath = jarPath
	metadata.MetadataJSON = string(data)

	// Extract authors
	if mod.Authors != "" {
		// Authors can be comma-separated
		authors := strings.Split(mod.Authors, ",")
		metadata.Authors = make([]string, 0, len(authors))
		for _, author := range authors {
			author = strings.TrimSpace(author)
			if author != "" {
				metadata.Authors = append(metadata.Authors, author)
			}
		}
	}

	// Extract dependencies
	metadata.Dependencies = fml.extractDependencies(forgeMod.Dependencies, metadata.ID)

	// Extract icon
	if mod.LogoFile != "" {
		metadata.IconData = fml.extractIconFromPath(zipReader, mod.LogoFile)
	}
	if metadata.IconData == "" {
		metadata.IconData = fml.iconExtractor.ExtractWithFallback(zipReader, metadata.ID, assets.DefaultModIconData)
	}

	return metadata, nil
}

// resolveVersionPlaceholder tries to resolve version placeholders like ${file.jarVersion}
func (fml *ForgeModernLoader) resolveVersionPlaceholder(zipReader *zip.Reader, version string) string {
	// Try to read from MANIFEST.MF
	for _, f := range zipReader.File {
		if f.Name == "META-INF/MANIFEST.MF" {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				continue
			}

			manifest := string(data)
			lines := strings.Split(manifest, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "Implementation-Version:") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	}

	// If we can't resolve, return a default or the original
	return strings.ReplaceAll(version, "${file.jarVersion}", "unknown")
}

// extractDependencies extracts dependencies from Forge mod metadata
func (fml *ForgeModernLoader) extractDependencies(deps []ForgeDependency, modID string) []*models.Dependency {
	result := make([]*models.Dependency, 0)

	for _, dep := range deps {
		// Skip forge, minecraft, and other core dependencies
		if dep.ModID == "forge" || dep.ModID == "minecraft" || dep.ModID == "java" {
			continue
		}

		required := dep.Mandatory
		versionRange, _ := models.NewVersionRange(dep.VersionRange)

		result = append(result, models.NewDependency(modID, dep.ModID, required, versionRange))
	}

	return result
}

// extractIconFromPath extracts icon from a specific path
func (fml *ForgeModernLoader) extractIconFromPath(zipReader *zip.Reader, iconPath string) string {
	for _, f := range zipReader.File {
		if f.Name == iconPath {
			icon, err := fml.iconExtractor.ExtractFile(f)
			if err == nil {
				return icon
			}
		}
	}
	return ""
}

// ForgeModsToml represents the structure of mods.toml
type ForgeModsToml struct {
	ModLoader          string            `toml:"modLoader"`
	LoaderVersion      string            `toml:"loaderVersion"`
	License            string            `toml:"license"`
	ShowAsResourcePack bool              `toml:"showAsResourcePack"`
	Mods               []ForgeMod        `toml:"mods"`
	Dependencies       []ForgeDependency `toml:"dependencies"`
}

// ForgeMod represents a mod entry in mods.toml
type ForgeMod struct {
	ModID       string `toml:"modId"`
	Version     string `toml:"version"`
	DisplayName string `toml:"displayName"`
	Description string `toml:"description"`
	LogoFile    string `toml:"logoFile"`
	Authors     string `toml:"authors"`
	Credits     string `toml:"credits"`
	DisplayURL  string `toml:"displayURL"`
}

// ForgeDependency represents a dependency entry
type ForgeDependency struct {
	ModID        string `toml:"modId"`
	Mandatory    bool   `toml:"mandatory"`
	VersionRange string `toml:"versionRange"`
	Ordering     string `toml:"ordering"`
	Side         string `toml:"side"`
}
