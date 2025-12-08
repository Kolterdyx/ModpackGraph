package models

import "time"

// ModFeature represents features provided by a mod (Phase 1 analysis - stub for now)
type ModFeature struct {
	ModID       string      `json:"mod_id"`
	FeatureType FeatureType `json:"feature_type"`
	FeatureData string      `json:"feature_data"` // JSON encoded
	ExtractedAt time.Time   `json:"extracted_at"`
}

type FeatureType string

// FeatureType constants for future use
const (
	FeatureTypeWorldGeneration FeatureType = "world_generation"
	FeatureTypeEntities        FeatureType = "entities"
	FeatureTypeItems           FeatureType = "items"
	FeatureTypeBlocks          FeatureType = "blocks"
	FeatureTypeMechanics       FeatureType = "mechanics"
	FeatureTypeRecipes         FeatureType = "recipes"
)

var AllFeatureTypes = []FeatureType{
	FeatureTypeWorldGeneration,
	FeatureTypeEntities,
	FeatureTypeItems,
	FeatureTypeBlocks,
	FeatureTypeMechanics,
	FeatureTypeRecipes,
}

func (ft FeatureType) TSName() string {
	return string(ft)
}

func (ft FeatureType) String() string {
	return string(ft)
}

// NewModFeature creates a new ModFeature
func NewModFeature(modID string, featureType FeatureType, featureData string) *ModFeature {
	return &ModFeature{
		ModID:       modID,
		FeatureType: featureType,
		FeatureData: featureData,
		ExtractedAt: time.Now(),
	}
}
