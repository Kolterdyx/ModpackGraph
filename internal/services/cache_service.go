package services

import (
	"ModpackGraph/internal/logger"
	"ModpackGraph/internal/models"
	"ModpackGraph/internal/repository"
	"fmt"
	"github.com/sirupsen/logrus"
)

// CacheService handles hash-based caching of mod metadata
type CacheService struct {
	modRepo         *repository.ModRepository
	metadataService *MetadataService
}

// NewCacheService creates a new CacheService
func NewCacheService(modRepo *repository.ModRepository, metadataService *MetadataService) *CacheService {
	return &CacheService{
		modRepo:         modRepo,
		metadataService: metadataService,
	}
}

// CacheResult represents the result of a cache lookup
type CacheResult struct {
	Metadata *models.ModMetadata
	Hash     string
	CacheHit bool
}

// GetOrExtract gets metadata from cache or extracts it from JAR
func (s *CacheService) GetOrExtract(jarPath string) (*CacheResult, error) {
	log := logger.GetLogger()

	// Compute hash
	hash, err := s.metadataService.ComputeHash(jarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to compute hash: %w", err)
	}

	// Check cache
	cached, err := s.modRepo.GetByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to query cache: %w", err)
	}

	if cached != nil {
		// Cache hit - update file path
		log.WithFields(logrus.Fields{
			logger.FieldModID: cached.ID,
			logger.FieldHash:  hash,
			logger.FieldPath:  jarPath,
		}).Debug("Cache hit")

		cached.FilePath = jarPath
		return &CacheResult{
			Metadata: cached,
			Hash:     hash,
			CacheHit: true,
		}, nil
	}

	// Cache miss - extract metadata
	log.WithFields(logrus.Fields{
		logger.FieldHash: hash,
		logger.FieldPath: jarPath,
	}).Debug("Cache miss, extracting metadata")

	metadata, err := s.metadataService.ExtractFromJAR(jarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata: %w", err)
	}

	// Save to cache
	if err := s.modRepo.Save(metadata); err != nil {
		return nil, fmt.Errorf("failed to save to cache: %w", err)
	}

	log.WithFields(logrus.Fields{
		logger.FieldModID:      metadata.ID,
		logger.FieldModName:    metadata.Name,
		logger.FieldModVersion: metadata.Version,
		logger.FieldLoaderType: metadata.LoaderType,
	}).Info("Cached new mod metadata")

	return &CacheResult{
		Metadata: metadata,
		Hash:     hash,
		CacheHit: false,
	}, nil
}

// SaveBatch saves multiple metadata entries in a batch
func (s *CacheService) SaveBatch(metadataList []*models.ModMetadata) error {
	return s.modRepo.SaveBatch(metadataList)
}

// IsCached checks if a hash exists in cache
func (s *CacheService) IsCached(hash string) (bool, error) {
	cached, err := s.modRepo.GetByHash(hash)
	if err != nil {
		return false, err
	}
	return cached != nil, nil
}

// GetByHash retrieves cached metadata by hash
func (s *CacheService) GetByHash(hash string) (*models.ModMetadata, error) {
	return s.modRepo.GetByHash(hash)
}

// GetByID retrieves cached metadata by mod ID
func (s *CacheService) GetByID(modID string) (*models.ModMetadata, error) {
	return s.modRepo.GetByID(modID)
}
