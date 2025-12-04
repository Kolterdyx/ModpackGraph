package repository

import (
	"ModpackGraph/internal/models"
	"database/sql"
	"fmt"
	"time"
)

// ModpackRepository handles modpack persistence
type ModpackRepository struct {
	db *DB
}

// NewModpackRepository creates a new ModpackRepository
func NewModpackRepository(db *DB) *ModpackRepository {
	return &ModpackRepository{db: db}
}

// GetByID retrieves a modpack by its ID
func (r *ModpackRepository) GetByID(id int64) (*models.Modpack, error) {
	query := `
		SELECT id, name, path, last_scanned, mod_count
		FROM modpacks
		WHERE id = ?
	`

	var modpack models.Modpack
	err := r.db.QueryRow(query, id).Scan(
		&modpack.ID,
		&modpack.Name,
		&modpack.Path,
		&modpack.LastScanned,
		&modpack.ModCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack by id: %w", err)
	}

	return &modpack, nil
}

// GetByPath retrieves a modpack by its path
func (r *ModpackRepository) GetByPath(path string) (*models.Modpack, error) {
	query := `
		SELECT id, name, path, last_scanned, mod_count
		FROM modpacks
		WHERE path = ?
	`

	var modpack models.Modpack
	err := r.db.QueryRow(query, path).Scan(
		&modpack.ID,
		&modpack.Name,
		&modpack.Path,
		&modpack.LastScanned,
		&modpack.ModCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack by path: %w", err)
	}

	return &modpack, nil
}

// Save saves or updates a modpack
func (r *ModpackRepository) Save(modpack *models.Modpack) error {
	if modpack.ID == 0 {
		// Insert new modpack
		query := `
			INSERT INTO modpacks (name, path, last_scanned, mod_count)
			VALUES (?, ?, ?, ?)
		`
		result, err := r.db.Exec(query, modpack.Name, modpack.Path, modpack.LastScanned, modpack.ModCount)
		if err != nil {
			return fmt.Errorf("failed to insert modpack: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		modpack.ID = id
	} else {
		// Update existing modpack
		query := `
			UPDATE modpacks
			SET name = ?, path = ?, last_scanned = ?, mod_count = ?
			WHERE id = ?
		`
		_, err := r.db.Exec(query, modpack.Name, modpack.Path, modpack.LastScanned, modpack.ModCount, modpack.ID)
		if err != nil {
			return fmt.Errorf("failed to update modpack: %w", err)
		}
	}

	return nil
}

// Delete deletes a modpack by ID
func (r *ModpackRepository) Delete(id int64) error {
	query := "DELETE FROM modpacks WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete modpack: %w", err)
	}
	return nil
}

// GetAll retrieves all modpacks
func (r *ModpackRepository) GetAll() ([]*models.Modpack, error) {
	query := `
		SELECT id, name, path, last_scanned, mod_count
		FROM modpacks
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all modpacks: %w", err)
	}
	defer rows.Close()

	modpacks := make([]*models.Modpack, 0)
	for rows.Next() {
		var modpack models.Modpack
		if err := rows.Scan(
			&modpack.ID,
			&modpack.Name,
			&modpack.Path,
			&modpack.LastScanned,
			&modpack.ModCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan modpack: %w", err)
		}
		modpacks = append(modpacks, &modpack)
	}

	return modpacks, rows.Err()
}

// GetMods retrieves all mods in a modpack
func (r *ModpackRepository) GetMods(modpackID int64) ([]*models.ModpackMod, error) {
	query := `
		SELECT modpack_id, mod_id, hash, file_path
		FROM modpack_mods
		WHERE modpack_id = ?
		ORDER BY file_path
	`

	rows, err := r.db.Query(query, modpackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack mods: %w", err)
	}
	defer rows.Close()

	mods := make([]*models.ModpackMod, 0)
	for rows.Next() {
		var mod models.ModpackMod
		if err := rows.Scan(
			&mod.ModpackID,
			&mod.ModID,
			&mod.Hash,
			&mod.FilePath,
		); err != nil {
			return nil, fmt.Errorf("failed to scan modpack mod: %w", err)
		}
		mods = append(mods, &mod)
	}

	return mods, rows.Err()
}

// AddMod adds a mod to a modpack
func (r *ModpackRepository) AddMod(modpackID int64, modID, hash, filePath string) error {
	query := `
		INSERT INTO modpack_mods (modpack_id, mod_id, hash, file_path)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(modpack_id, mod_id) DO UPDATE SET
			hash = excluded.hash,
			file_path = excluded.file_path
	`

	_, err := r.db.Exec(query, modpackID, modID, hash, filePath)
	if err != nil {
		return fmt.Errorf("failed to add mod to modpack: %w", err)
	}

	return nil
}

// RemoveMod removes a mod from a modpack
func (r *ModpackRepository) RemoveMod(modpackID int64, modID string) error {
	query := "DELETE FROM modpack_mods WHERE modpack_id = ? AND mod_id = ?"
	_, err := r.db.Exec(query, modpackID, modID)
	if err != nil {
		return fmt.Errorf("failed to remove mod from modpack: %w", err)
	}
	return nil
}

// ClearMods removes all mods from a modpack
func (r *ModpackRepository) ClearMods(modpackID int64) error {
	query := "DELETE FROM modpack_mods WHERE modpack_id = ?"
	_, err := r.db.Exec(query, modpackID)
	if err != nil {
		return fmt.Errorf("failed to clear modpack mods: %w", err)
	}
	return nil
}

// UpdateLastScanned updates the last_scanned timestamp
func (r *ModpackRepository) UpdateLastScanned(modpackID int64) error {
	query := "UPDATE modpacks SET last_scanned = ? WHERE id = ?"
	_, err := r.db.Exec(query, time.Now(), modpackID)
	if err != nil {
		return fmt.Errorf("failed to update last_scanned: %w", err)
	}
	return nil
}

// UpdateModCount updates the mod_count field
func (r *ModpackRepository) UpdateModCount(modpackID int64, count int) error {
	query := "UPDATE modpacks SET mod_count = ? WHERE id = ?"
	_, err := r.db.Exec(query, count, modpackID)
	if err != nil {
		return fmt.Errorf("failed to update mod_count: %w", err)
	}
	return nil
}

// SetMods replaces all mods in a modpack with the given list
func (r *ModpackRepository) SetMods(modpackID int64, mods []*models.ModpackMod) error {
	tx, err := r.db.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Clear existing mods
	if _, err := tx.Exec("DELETE FROM modpack_mods WHERE modpack_id = ?", modpackID); err != nil {
		return fmt.Errorf("failed to clear mods: %w", err)
	}

	// Insert new mods
	if len(mods) > 0 {
		query := "INSERT INTO modpack_mods (modpack_id, mod_id, hash, file_path) VALUES (?, ?, ?, ?)"
		stmt, err := tx.Prepare(query)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, mod := range mods {
			if _, err := stmt.Exec(modpackID, mod.ModID, mod.Hash, mod.FilePath); err != nil {
				return fmt.Errorf("failed to insert mod: %w", err)
			}
		}
	}

	// Update mod count
	if _, err := tx.Exec("UPDATE modpacks SET mod_count = ? WHERE id = ?", len(mods), modpackID); err != nil {
		return fmt.Errorf("failed to update mod count: %w", err)
	}

	return tx.Commit()
}
