package services

import (
	"ModpackGraph/internal/logger"
	"ModpackGraph/internal/models"
	"ModpackGraph/internal/repository"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ScanService handles modpack directory scanning
type ScanService struct {
	cacheService    *CacheService
	modpackRepo     *repository.ModpackRepository
	metadataService *MetadataService
}

// NewScanService creates a new ScanService
func NewScanService(cacheService *CacheService, modpackRepo *repository.ModpackRepository, metadataService *MetadataService) *ScanService {
	return &ScanService{
		cacheService:    cacheService,
		modpackRepo:     modpackRepo,
		metadataService: metadataService,
	}
}

// ScanResult represents the result of scanning a modpack
type ScanResult struct {
	Modpack     *models.Modpack
	Mods        []*models.ModMetadata
	NewMods     []*models.ModMetadata
	UpdatedMods []*models.ModMetadata
	RemovedMods []*models.ModMetadata
	CacheHits   int
	CacheMisses int
}

// ScanModpack scans a modpack directory for JAR files
func (s *ScanService) ScanModpack(modpackPath string) (*ScanResult, error) {
	log := logger.GetLogger()
	log.WithField(logger.FieldPath, modpackPath).Info("Starting modpack scan")

	// Normalize path
	absPath, err := filepath.Abs(modpackPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", absPath)
	}

	// Get or create modpack record
	modpack, err := s.modpackRepo.GetByPath(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack: %w", err)
	}

	if modpack == nil {
		// Create new modpack
		modpack = models.NewModpack(filepath.Base(absPath), absPath)
		if err := s.modpackRepo.Save(modpack); err != nil {
			return nil, fmt.Errorf("failed to save modpack: %w", err)
		}
	}

	// Get existing mods in modpack
	existingMods, err := s.modpackRepo.GetMods(modpack.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing mods: %w", err)
	}

	existingModMap := make(map[string]*models.ModpackMod)
	for _, mod := range existingMods {
		existingModMap[mod.FilePath] = mod
	}

	// Scan for JAR files
	jarFiles, err := s.findJARFiles(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find JAR files: %w", err)
	}

	result := &ScanResult{
		Modpack:     modpack,
		Mods:        make([]*models.ModMetadata, 0),
		NewMods:     make([]*models.ModMetadata, 0),
		UpdatedMods: make([]*models.ModMetadata, 0),
		RemovedMods: make([]*models.ModMetadata, 0),
		CacheHits:   0,
		CacheMisses: 0,
	}

	// Process each JAR file
	currentMods := make([]*models.ModpackMod, 0)
	processedPaths := make(map[string]bool)

	for _, jarPath := range jarFiles {
		relPath, err := filepath.Rel(absPath, jarPath)
		if err != nil {
			relPath = jarPath
		}

		processedPaths[relPath] = true

		// Get or extract metadata
		cacheResult, err := s.cacheService.GetOrExtract(jarPath)
		if err != nil {
			// Log error but continue
			log.WithFields(logrus.Fields{
				logger.FieldPath:  jarPath,
				logger.FieldError: err,
			}).Warn("Failed to process JAR file")
			continue
		}

		result.Mods = append(result.Mods, cacheResult.Metadata)

		if cacheResult.CacheHit {
			result.CacheHits++
		} else {
			result.CacheMisses++
			result.NewMods = append(result.NewMods, cacheResult.Metadata)
		}

		// Check if this is an update
		if existing, ok := existingModMap[relPath]; ok {
			if existing.Hash != cacheResult.Hash {
				result.UpdatedMods = append(result.UpdatedMods, cacheResult.Metadata)
			}
		}

		// Add to current mods list
		currentMods = append(currentMods, &models.ModpackMod{
			ModpackID: modpack.ID,
			ModID:     cacheResult.Metadata.ID,
			Hash:      cacheResult.Hash,
			FilePath:  relPath,
		})
	}

	// Find removed mods
	for relPath, existing := range existingModMap {
		if !processedPaths[relPath] {
			// Mod was removed
			if metadata, err := s.cacheService.GetByID(existing.ModID); err == nil && metadata != nil {
				result.RemovedMods = append(result.RemovedMods, metadata)
			}
		}
	}

	// Update modpack_mods table
	if err := s.modpackRepo.SetMods(modpack.ID, currentMods); err != nil {
		return nil, fmt.Errorf("failed to update modpack mods: %w", err)
	}

	// Update modpack timestamp
	modpack.LastScanned = time.Now()
	modpack.ModCount = len(currentMods)
	if err := s.modpackRepo.Save(modpack); err != nil {
		return nil, fmt.Errorf("failed to update modpack: %w", err)
	}

	result.Modpack = modpack

	log.WithFields(logrus.Fields{
		logger.FieldPath:   absPath,
		logger.TotalMods:   len(currentMods),
		logger.NewMods:     len(result.NewMods),
		logger.UpdatedMods: len(result.UpdatedMods),
		logger.RemovedMods: len(result.RemovedMods),
		logger.CacheHits:   result.CacheHits,
		logger.CacheMisses: result.CacheMisses,
	}).Info("Modpack scan completed")

	return result, nil
}

// findJARFiles recursively finds all JAR files in a directory
func (s *ScanService) findJARFiles(root string) ([]string, error) {
	jarFiles := make([]string, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".jar" {
			jarFiles = append(jarFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return jarFiles, nil
}

// GetModpackByPath retrieves a modpack by its path
func (s *ScanService) GetModpackByPath(path string) (*models.Modpack, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	return s.modpackRepo.GetByPath(absPath)
}

// GetModpackMods retrieves all mods in a modpack
func (s *ScanService) GetModpackMods(modpackID int64) ([]*models.ModMetadata, error) {
	modpackMods, err := s.modpackRepo.GetMods(modpackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack mods: %w", err)
	}

	if len(modpackMods) == 0 {
		return []*models.ModMetadata{}, nil
	}

	// Get mod IDs
	modIDs := make([]string, len(modpackMods))
	for i, mm := range modpackMods {
		modIDs[i] = mm.ModID
	}

	// Fetch metadata for all mods
	mods, err := s.cacheService.modRepo.GetByIDs(modIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get mods: %w", err)
	}

	return mods, nil
}
