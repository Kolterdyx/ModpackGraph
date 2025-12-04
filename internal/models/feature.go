package models

import "time"

// ModFeature represents features provided by a mod (Phase 1 analysis - stub for now)
type ModFeature struct {
	ModID       string    `json:"mod_id"`
	FeatureType string    `json:"feature_type"`
	FeatureData string    `json:"feature_data"` // JSON encoded
	ExtractedAt time.Time `json:"extracted_at"`
}

// FeatureType constants for future use
const (
	FeatureTypeWorldGeneration = "world_generation"
	FeatureTypeEntities        = "entities"
	FeatureTypeItems           = "items"
	FeatureTypeBlocks          = "blocks"
	FeatureTypeMechanics       = "mechanics"
	FeatureTypeRecipes         = "recipes"
)

// NewModFeature creates a new ModFeature
func NewModFeature(modID, featureType, featureData string) *ModFeature {
	return &ModFeature{
		ModID:       modID,
		FeatureType: featureType,
		FeatureData: featureData,
		ExtractedAt: time.Now(),
	}
}
