package repository

import (
	"ModpackGraph/internal/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// ModRepository handles mod metadata persistence
type ModRepository struct {
	db *DB
}

// NewModRepository creates a new ModRepository
func NewModRepository(db *DB) *ModRepository {
	return &ModRepository{db: db}
}

// GetByID retrieves a mod by its ID
func (r *ModRepository) GetByID(id string) (*models.ModMetadata, error) {
	query := `
		SELECT id, hash, version, name, description, loader_type, environment, 
		       icon_data, metadata_json, created_at, updated_at
		FROM mods
		WHERE id = ?
	`

	var metadata models.ModMetadata
	var description, iconData, metadataJSON sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&metadata.ID,
		&metadata.Hash,
		&metadata.Version,
		&metadata.Name,
		&description,
		&metadata.LoaderType,
		&metadata.Environment,
		&iconData,
		&metadataJSON,
		&metadata.CreatedAt,
		&metadata.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get mod by id: %w", err)
	}

	metadata.Description = description.String
	metadata.IconData = iconData.String
	metadata.MetadataJSON = metadataJSON.String

	// Load authors
	authors, err := r.getAuthors(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get authors: %w", err)
	}
	metadata.Authors = authors

	// Load dependencies
	deps, err := r.getDependencies(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}
	metadata.Dependencies = deps

	return &metadata, nil
}

// GetByHash retrieves a mod by its hash
func (r *ModRepository) GetByHash(hash string) (*models.ModMetadata, error) {
	query := `
		SELECT id, hash, version, name, description, loader_type, environment, 
		       icon_data, metadata_json, created_at, updated_at
		FROM mods
		WHERE hash = ?
	`

	var metadata models.ModMetadata
	var description, iconData, metadataJSON sql.NullString
	err := r.db.QueryRow(query, hash).Scan(
		&metadata.ID,
		&metadata.Hash,
		&metadata.Version,
		&metadata.Name,
		&description,
		&metadata.LoaderType,
		&metadata.Environment,
		&iconData,
		&metadataJSON,
		&metadata.CreatedAt,
		&metadata.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get mod by hash: %w", err)
	}

	metadata.Description = description.String
	metadata.IconData = iconData.String
	metadata.MetadataJSON = metadataJSON.String

	// Load authors
	authors, err := r.getAuthors(metadata.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get authors: %w", err)
	}
	metadata.Authors = authors

	// Load dependencies
	deps, err := r.getDependencies(metadata.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}
	metadata.Dependencies = deps

	return &metadata, nil
}

// Save saves or updates a mod
func (r *ModRepository) Save(metadata *models.ModMetadata) error {
	tx, err := r.db.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Upsert mod
	query := `
		INSERT INTO mods (id, hash, version, name, description, loader_type, environment, icon_data, metadata_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			hash = excluded.hash,
			version = excluded.version,
			name = excluded.name,
			description = excluded.description,
			loader_type = excluded.loader_type,
			environment = excluded.environment,
			icon_data = excluded.icon_data,
			metadata_json = excluded.metadata_json,
			updated_at = excluded.updated_at
	`

	now := time.Now()
	if metadata.CreatedAt.IsZero() {
		metadata.CreatedAt = now
	}
	metadata.UpdatedAt = now

	_, err = tx.Exec(query,
		metadata.ID,
		metadata.Hash,
		metadata.Version,
		metadata.Name,
		nullString(metadata.Description),
		metadata.LoaderType,
		metadata.Environment,
		nullString(metadata.IconData),
		nullString(metadata.MetadataJSON),
		metadata.CreatedAt,
		metadata.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save mod: %w", err)
	}

	// Delete old authors and dependencies
	if _, err := tx.Exec("DELETE FROM mod_authors WHERE mod_id = ?", metadata.ID); err != nil {
		return fmt.Errorf("failed to delete old authors: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM mod_dependencies WHERE mod_id = ?", metadata.ID); err != nil {
		return fmt.Errorf("failed to delete old dependencies: %w", err)
	}

	// Insert authors
	if err := r.saveAuthors(tx, metadata.ID, metadata.Authors); err != nil {
		return fmt.Errorf("failed to save authors: %w", err)
	}

	// Insert dependencies
	if err := r.saveDependencies(tx, metadata.Dependencies); err != nil {
		return fmt.Errorf("failed to save dependencies: %w", err)
	}

	return tx.Commit()
}

// SaveBatch saves multiple mods in a single transaction
func (r *ModRepository) SaveBatch(metadataList []*models.ModMetadata) error {
	tx, err := r.db.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, metadata := range metadataList {
		if err := r.saveTx(tx, metadata); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// saveTx saves a mod within a transaction
func (r *ModRepository) saveTx(tx *sql.Tx, metadata *models.ModMetadata) error {
	query := `
		INSERT INTO mods (id, hash, version, name, description, loader_type, environment, icon_data, metadata_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			hash = excluded.hash,
			version = excluded.version,
			name = excluded.name,
			description = excluded.description,
			loader_type = excluded.loader_type,
			environment = excluded.environment,
			icon_data = excluded.icon_data,
			metadata_json = excluded.metadata_json,
			updated_at = excluded.updated_at
	`

	now := time.Now()
	if metadata.CreatedAt.IsZero() {
		metadata.CreatedAt = now
	}
	metadata.UpdatedAt = now

	_, err := tx.Exec(query,
		metadata.ID,
		metadata.Hash,
		metadata.Version,
		metadata.Name,
		nullString(metadata.Description),
		metadata.LoaderType,
		metadata.Environment,
		nullString(metadata.IconData),
		nullString(metadata.MetadataJSON),
		metadata.CreatedAt,
		metadata.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save mod: %w", err)
	}

	// Delete old authors and dependencies
	if _, err := tx.Exec("DELETE FROM mod_authors WHERE mod_id = ?", metadata.ID); err != nil {
		return fmt.Errorf("failed to delete old authors: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM mod_dependencies WHERE mod_id = ?", metadata.ID); err != nil {
		return fmt.Errorf("failed to delete old dependencies: %w", err)
	}

	// Insert authors
	if err := r.saveAuthors(tx, metadata.ID, metadata.Authors); err != nil {
		return fmt.Errorf("failed to save authors: %w", err)
	}

	// Insert dependencies
	if err := r.saveDependencies(tx, metadata.Dependencies); err != nil {
		return fmt.Errorf("failed to save dependencies: %w", err)
	}

	return nil
}

// getAuthors retrieves authors for a mod
func (r *ModRepository) getAuthors(modID string) ([]string, error) {
	query := "SELECT author FROM mod_authors WHERE mod_id = ? ORDER BY author"
	rows, err := r.db.Query(query, modID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	authors := make([]string, 0)
	for rows.Next() {
		var author string
		if err := rows.Scan(&author); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, rows.Err()
}

// saveAuthors saves authors for a mod
func (r *ModRepository) saveAuthors(tx *sql.Tx, modID string, authors []string) error {
	if len(authors) == 0 {
		return nil
	}

	query := "INSERT INTO mod_authors (mod_id, author) VALUES (?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, author := range authors {
		if _, err := stmt.Exec(modID, author); err != nil {
			return err
		}
	}

	return nil
}

// getDependencies retrieves dependencies for a mod
func (r *ModRepository) getDependencies(modID string) ([]*models.Dependency, error) {
	query := "SELECT mod_id, dependency_id, required, version_range FROM mod_dependencies WHERE mod_id = ?"
	rows, err := r.db.Query(query, modID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deps := make([]*models.Dependency, 0)
	for rows.Next() {
		var dep models.Dependency
		var versionRange sql.NullString
		if err := rows.Scan(&dep.ModID, &dep.DependencyID, &dep.Required, &versionRange); err != nil {
			return nil, err
		}

		if versionRange.Valid && versionRange.String != "" {
			vr, _ := models.NewVersionRange(versionRange.String)
			dep.VersionRange = vr
		}

		deps = append(deps, &dep)
	}

	return deps, rows.Err()
}

// saveDependencies saves dependencies for a mod
func (r *ModRepository) saveDependencies(tx *sql.Tx, deps []*models.Dependency) error {
	if len(deps) == 0 {
		return nil
	}

	query := "INSERT INTO mod_dependencies (mod_id, dependency_id, required, version_range) VALUES (?, ?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, dep := range deps {
		var versionRange string
		if dep.VersionRange != nil {
			versionRange = dep.VersionRange.String()
		}
		if _, err := stmt.Exec(dep.ModID, dep.DependencyID, dep.Required, nullString(versionRange)); err != nil {
			return err
		}
	}

	return nil
}

// GetAll retrieves all mods
func (r *ModRepository) GetAll() ([]*models.ModMetadata, error) {
	query := `
		SELECT id, hash, version, name, description, loader_type, environment, 
		       icon_data, metadata_json, created_at, updated_at
		FROM mods
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all mods: %w", err)
	}
	defer rows.Close()

	mods := make([]*models.ModMetadata, 0)
	for rows.Next() {
		var metadata models.ModMetadata
		var description, iconData, metadataJSON sql.NullString
		if err := rows.Scan(
			&metadata.ID,
			&metadata.Hash,
			&metadata.Version,
			&metadata.Name,
			&description,
			&metadata.LoaderType,
			&metadata.Environment,
			&iconData,
			&metadataJSON,
			&metadata.CreatedAt,
			&metadata.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan mod: %w", err)
		}

		metadata.Description = description.String
		metadata.IconData = iconData.String
		metadata.MetadataJSON = metadataJSON.String

		// Load authors and dependencies
		metadata.Authors, _ = r.getAuthors(metadata.ID)
		metadata.Dependencies, _ = r.getDependencies(metadata.ID)

		mods = append(mods, &metadata)
	}

	return mods, rows.Err()
}

// Delete deletes a mod by ID
func (r *ModRepository) Delete(id string) error {
	query := "DELETE FROM mods WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete mod: %w", err)
	}
	return nil
}

// GetByIDs retrieves multiple mods by their IDs
func (r *ModRepository) GetByIDs(ids []string) ([]*models.ModMetadata, error) {
	if len(ids) == 0 {
		return []*models.ModMetadata{}, nil
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
		SELECT id, hash, version, name, description, loader_type, environment, 
		       icon_data, metadata_json, created_at, updated_at
		FROM mods
		WHERE id IN (%s)
	`, placeholders)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get mods by ids: %w", err)
	}
	defer rows.Close()

	mods := make([]*models.ModMetadata, 0)
	for rows.Next() {
		var metadata models.ModMetadata
		var description, iconData, metadataJSON sql.NullString
		if err := rows.Scan(
			&metadata.ID,
			&metadata.Hash,
			&metadata.Version,
			&metadata.Name,
			&description,
			&metadata.LoaderType,
			&metadata.Environment,
			&iconData,
			&metadataJSON,
			&metadata.CreatedAt,
			&metadata.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan mod: %w", err)
		}

		metadata.Description = description.String
		metadata.IconData = iconData.String
		metadata.MetadataJSON = metadataJSON.String

		// Load authors and dependencies
		metadata.Authors, _ = r.getAuthors(metadata.ID)
		metadata.Dependencies, _ = r.getDependencies(metadata.ID)

		mods = append(mods, &metadata)
	}

	return mods, rows.Err()
}

// nullString returns sql.NullString
func nullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
