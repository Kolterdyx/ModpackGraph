package services

import (
	"ModpackGraph/internal/models"
	"ModpackGraph/internal/services/loaders"
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// MetadataService handles JAR file processing and metadata extraction
type MetadataService struct {
	loaderRegistry *loaders.LoaderRegistry
}

// NewMetadataService creates a new MetadataService
func NewMetadataService() *MetadataService {
	return &MetadataService{
		loaderRegistry: loaders.NewLoaderRegistry(),
	}
}

// ExtractFromJAR extracts metadata from a JAR file
func (s *MetadataService) ExtractFromJAR(jarPath string) (*models.ModMetadata, error) {
	// Compute hash
	hash, err := s.ComputeHash(jarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to compute hash: %w", err)
	}

	// Open JAR as zip
	zipReader, err := zip.OpenReader(jarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JAR: %w", err)
	}
	defer zipReader.Close()

	// Extract metadata using loader registry
	metadata, err := s.loaderRegistry.ExtractMetadata(&zipReader.Reader, jarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata: %w", err)
	}

	// Set hash and file path
	metadata.Hash = hash
	metadata.FilePath = jarPath

	return metadata, nil
}

// ComputeHash computes SHA-256 hash of a file
func (s *MetadataService) ComputeHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("failed to compute hash: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// DetectLoader detects which loader a JAR uses
func (s *MetadataService) DetectLoader(jarPath string) (string, error) {
	zipReader, err := zip.OpenReader(jarPath)
	if err != nil {
		return "", fmt.Errorf("failed to open JAR: %w", err)
	}
	defer zipReader.Close()

	loader, err := s.loaderRegistry.DetectLoader(&zipReader.Reader)
	if err != nil {
		return "", err
	}

	// Determine loader type name
	switch loader.(type) {
	case *loaders.FabricLoader:
		return "fabric", nil
	case *loaders.ForgeModernLoader:
		return "forge_modern", nil
	case *loaders.ForgeLegacyLoader:
		return "forge_legacy", nil
	case *loaders.NeoForgeLoader:
		return "neoforge", nil
	default:
		return "unknown", nil
	}
}
