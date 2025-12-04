package models

// Dependency represents a mod dependency relationship
type Dependency struct {
	ModID        string        `json:"mod_id"`
	DependencyID string        `json:"dependency_id"`
	Required     bool          `json:"required"`
	VersionRange *VersionRange `json:"version_range,omitempty"`
}

// NewDependency creates a new dependency
func NewDependency(modID, dependencyID string, required bool, versionRange *VersionRange) *Dependency {
	return &Dependency{
		ModID:        modID,
		DependencyID: dependencyID,
		Required:     required,
		VersionRange: versionRange,
	}
}
