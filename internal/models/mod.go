package models

import "time"

// Mod represents basic mod identity
type Mod struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

// NewMod creates a new Mod
func NewMod(id, version string) *Mod {
	return &Mod{
		ID:      id,
		Version: version,
	}
}

// ModMetadata represents full mod information with metadata
type ModMetadata struct {
	ID           string        `json:"id"`
	Hash         string        `json:"hash"`
	Version      string        `json:"version"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Authors      []string      `json:"authors"`
	LoaderType   LoaderType    `json:"loader_type"`
	Environment  Environment   `json:"environment"`
	IconData     string        `json:"icon_data,omitempty"` // Base64 encoded
	Dependencies []*Dependency `json:"dependencies"`
	MetadataJSON string        `json:"metadata_json,omitempty"` // Raw metadata
	FilePath     string        `json:"file_path,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// NewModMetadata creates a new ModMetadata
func NewModMetadata() *ModMetadata {
	now := time.Now()
	return &ModMetadata{
		Dependencies: make([]*Dependency, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// ToMod converts ModMetadata to a simple Mod
func (m *ModMetadata) ToMod() *Mod {
	return NewMod(m.ID, m.Version)
}
