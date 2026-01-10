package loaders

import (
	"ModpackGraph/internal/assets"
	"ModpackGraph/internal/logger"
	"ModpackGraph/internal/models"
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ForgeLegacyLoader handles Forge 1.12.2 and earlier mod metadata extraction
type ForgeLegacyLoader struct {
	iconExtractor IconExtractor
}

// NewForgeLegacyLoader creates a new ForgeLegacyLoader
func NewForgeLegacyLoader(
	iconExtractor IconExtractor,
) ModLoader {
	return &ForgeLegacyLoader{
		iconExtractor: iconExtractor,
	}
}

// CanHandle checks if this is a legacy Forge mod
func (fll *ForgeLegacyLoader) CanHandle(zipReader *zip.Reader) bool {
	logger.GetLogger().Debugf("Checking for Forge legacy mod...")
	for _, f := range zipReader.File {
		if f.Name == "mcmod.info" {
			return true
		}
	}
	return false
}

// ExtractMetadata extracts metadata from a legacy Forge mod
func (fll *ForgeLegacyLoader) ExtractMetadata(zipReader *zip.Reader, jarPath string) (*models.ModMetadata, error) {
	// Find mcmod.info
	var mcmodInfo *zip.File
	for _, f := range zipReader.File {
		if f.Name == "mcmod.info" {
			mcmodInfo = f
			break
		}
	}

	if mcmodInfo == nil {
		return nil, fmt.Errorf("mcmod.info not found")
	}

	// Read and parse mcmod.info
	rc, err := mcmodInfo.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open mcmod.info: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read mcmod.info: %w", err)
	}

	// mcmod.info can be an array or an object with "modList" array
	var mods []LegacyModInfo

	// Try parsing as array first
	if err := json.Unmarshal(data, &mods); err != nil {
		// Try parsing as object with modList
		var wrapper struct {
			ModList []LegacyModInfo `json:"modList"`
		}
		if err2 := json.Unmarshal(data, &wrapper); err2 != nil {
			return nil, fmt.Errorf("failed to parse mcmod.info: %w", err)
		}
		mods = wrapper.ModList
	}

	if len(mods) == 0 {
		return nil, fmt.Errorf("no mods defined in mcmod.info")
	}

	mod := mods[0]

	// Handle version placeholders
	version := mod.Version
	if strings.Contains(version, "${") {
		version = fll.resolveVersionPlaceholder(zipReader, version)
	}

	// Create ModMetadata
	metadata := models.NewModMetadata()
	metadata.ID = mod.ModID
	metadata.Version = version
	metadata.Name = mod.Name
	metadata.Description = mod.Description
	metadata.LoaderType = models.LoaderTypeForgeLegacy
	metadata.Environment = models.EnvironmentBoth // Legacy Forge mods default to both
	metadata.FilePath = jarPath
	metadata.MetadataJSON = string(data)

	// Extract authors
	if len(mod.AuthorList) > 0 {
		metadata.Authors = mod.AuthorList
	} else if mod.Authors != "" {
		// Sometimes authors is a single string
		metadata.Authors = []string{mod.Authors}
	}

	// Extract dependencies
	metadata.Dependencies = fll.extractDependencies(mod, metadata.ID)

	// Extract icon
	if mod.LogoFile != "" {
		metadata.IconData = fll.extractIconFromPath(zipReader, mod.LogoFile)
	}
	if metadata.IconData == "" {
		metadata.IconData = fll.iconExtractor.ExtractWithFallback(zipReader, metadata.ID, assets.DefaultModIconData)
	}

	return metadata, nil
}

// resolveVersionPlaceholder tries to resolve version placeholders like ${version}
func (fll *ForgeLegacyLoader) resolveVersionPlaceholder(zipReader *zip.Reader, version string) string {
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
	return strings.ReplaceAll(version, "${version}", "unknown")
}

// extractDependencies extracts dependencies from legacy Forge mod metadata
func (fll *ForgeLegacyLoader) extractDependencies(mod LegacyModInfo, modID string) []*models.Dependency {
	result := make([]*models.Dependency, 0)

	// Parse requiredMods (format: "modid@[version]")
	for _, req := range mod.RequiredMods {
		dep := fll.parseDependencyString(req, modID, true)
		if dep != nil {
			result = append(result, dep)
		}
	}

	// Parse dependencies (format: "modid@[version]")
	for _, dep := range mod.Dependencies {
		parsed := fll.parseDependencyString(dep, modID, true)
		if parsed != nil {
			result = append(result, parsed)
		}
	}

	return result
}

// parseDependencyString parses dependency strings like "modid@[1.0,2.0)"
func (fll *ForgeLegacyLoader) parseDependencyString(depStr string, modID string, required bool) *models.Dependency {
	// Skip forge, minecraft, and other core dependencies
	if depStr == "" || strings.HasPrefix(depStr, "Forge") || strings.HasPrefix(depStr, "forge") {
		return nil
	}

	parts := strings.SplitN(depStr, "@", 2)
	depID := parts[0]

	if depID == "minecraft" || depID == "Minecraft" {
		return nil
	}

	var versionRange *models.VersionRange
	if len(parts) == 2 {
		versionRange, _ = models.NewVersionRange(parts[1])
	} else {
		versionRange, _ = models.NewVersionRange("*")
	}

	return models.NewDependency(modID, depID, required, versionRange)
}

// extractIconFromPath extracts icon from a specific path
func (fll *ForgeLegacyLoader) extractIconFromPath(zipReader *zip.Reader, iconPath string) string {
	for _, f := range zipReader.File {
		if f.Name == iconPath {
			icon, err := fll.iconExtractor.ExtractFile(f)
			if err == nil {
				return icon
			}
		}
	}
	return ""
}

// LegacyModInfo represents the structure of mcmod.info
type LegacyModInfo struct {
	ModID                    string   `json:"modid"`
	Name                     string   `json:"name"`
	Description              string   `json:"description"`
	Version                  string   `json:"version"`
	MCVersion                string   `json:"mcversion"`
	URL                      string   `json:"url"`
	UpdateURL                string   `json:"updateUrl"`
	AuthorList               []string `json:"authorList"`
	Authors                  string   `json:"authors"` // Fallback for single author
	Credits                  string   `json:"credits"`
	LogoFile                 string   `json:"logoFile"`
	Screenshots              []string `json:"screenshots"`
	Parent                   string   `json:"parent"`
	RequiredMods             []string `json:"requiredMods"`
	Dependencies             []string `json:"dependencies"`
	Dependants               []string `json:"dependants"`
	UseDependencyInformation bool     `json:"useDependencyInformation"`
}
