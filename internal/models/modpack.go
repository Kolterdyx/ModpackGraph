package models

import "time"

// Modpack represents a modpack configuration
type Modpack struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	LastScanned time.Time `json:"last_scanned"`
	ModCount    int       `json:"mod_count"`
	Mods        []*Mod    `json:"mods,omitempty"`
}

// NewModpack creates a new Modpack
func NewModpack(name, path string) *Modpack {
	return &Modpack{
		Name:        name,
		Path:        path,
		LastScanned: time.Now(),
		Mods:        make([]*Mod, 0),
	}
}

// ModpackMod represents a mod's presence in a modpack
type ModpackMod struct {
	ModpackID int64  `json:"modpack_id"`
	ModID     string `json:"mod_id"`
	Hash      string `json:"hash"`
	FilePath  string `json:"file_path"`
}
