package models

// ConflictType represents the type of conflict
type ConflictType string

const (
	ConflictTypeMissingDependency   ConflictType = "missing_dependency"
	ConflictTypeVersionConflict     ConflictType = "version_conflict"
	ConflictTypeKnownIncompatible   ConflictType = "known_incompatible"
	ConflictTypeFeatureOverlap      ConflictType = "feature_overlap"
	ConflictTypeEnvironmentMismatch ConflictType = "environment_mismatch"
	ConflictTypeCircularDependency  ConflictType = "circular_dependency"
)

// ConflictSeverity represents how severe a conflict is
type ConflictSeverity string

const (
	ConflictSeverityCritical ConflictSeverity = "critical"
	ConflictSeverityWarning  ConflictSeverity = "warning"
	ConflictSeverityInfo     ConflictSeverity = "info"
)

// Conflict represents a detected conflict between mods
type Conflict struct {
	Type         ConflictType     `json:"type"`
	Severity     ConflictSeverity `json:"severity"`
	Description  string           `json:"description"`
	AffectedMods []string         `json:"affected_mods"`
	Details      map[string]any   `json:"details,omitempty"`
}

// NewConflict creates a new conflict
func NewConflict(conflictType ConflictType, severity ConflictSeverity, description string, affectedMods []string) *Conflict {
	return &Conflict{
		Type:         conflictType,
		Severity:     severity,
		Description:  description,
		AffectedMods: affectedMods,
		Details:      make(map[string]any),
	}
}
