package models

// ConflictRule represents a known incompatibility in the database
type ConflictRule struct {
	ID           int64            `json:"id"`
	ModIDA       string           `json:"mod_id_a"`
	ModIDB       string           `json:"mod_id_b"`
	ConflictType ConflictType     `json:"conflict_type"`
	Description  string           `json:"description"`
	Severity     ConflictSeverity `json:"severity"`
}

// NewConflictRule creates a new conflict rule
func NewConflictRule(modIDA, modIDB string, conflictType ConflictType, description string, severity ConflictSeverity) *ConflictRule {
	return &ConflictRule{
		ModIDA:       modIDA,
		ModIDB:       modIDB,
		ConflictType: conflictType,
		Description:  description,
		Severity:     severity,
	}
}
