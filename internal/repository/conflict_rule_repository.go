package repository

import (
	"ModpackGraph/internal/models"
	"database/sql"
	"fmt"
)

// ConflictRuleRepository handles conflict rules persistence
type ConflictRuleRepository struct {
	db *DB
}

// NewConflictRuleRepository creates a new ConflictRuleRepository
func NewConflictRuleRepository(db *DB) *ConflictRuleRepository {
	return &ConflictRuleRepository{db: db}
}

// GetByID retrieves a conflict rule by its ID
func (r *ConflictRuleRepository) GetByID(id int64) (*models.ConflictRule, error) {
	query := `
		SELECT id, mod_id_a, mod_id_b, conflict_type, description, severity
		FROM conflict_rules
		WHERE id = ?
	`

	var rule models.ConflictRule
	var description sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&rule.ID,
		&rule.ModIDA,
		&rule.ModIDB,
		&rule.ConflictType,
		&description,
		&rule.Severity,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get conflict rule by id: %w", err)
	}

	rule.Description = description.String
	return &rule, nil
}

// GetByModPair retrieves conflict rules for a pair of mods
func (r *ConflictRuleRepository) GetByModPair(modIDA, modIDB string) ([]*models.ConflictRule, error) {
	query := `
		SELECT id, mod_id_a, mod_id_b, conflict_type, description, severity
		FROM conflict_rules
		WHERE (mod_id_a = ? AND mod_id_b = ?) OR (mod_id_a = ? AND mod_id_b = ?)
	`

	rows, err := r.db.Query(query, modIDA, modIDB, modIDB, modIDA)
	if err != nil {
		return nil, fmt.Errorf("failed to get conflict rules: %w", err)
	}
	defer rows.Close()

	rules := make([]*models.ConflictRule, 0)
	for rows.Next() {
		var rule models.ConflictRule
		var description sql.NullString
		if err := rows.Scan(
			&rule.ID,
			&rule.ModIDA,
			&rule.ModIDB,
			&rule.ConflictType,
			&description,
			&rule.Severity,
		); err != nil {
			return nil, fmt.Errorf("failed to scan conflict rule: %w", err)
		}
		rule.Description = description.String
		rules = append(rules, &rule)
	}

	return rules, rows.Err()
}

// GetByMod retrieves all conflict rules involving a mod
func (r *ConflictRuleRepository) GetByMod(modID string) ([]*models.ConflictRule, error) {
	query := `
		SELECT id, mod_id_a, mod_id_b, conflict_type, description, severity
		FROM conflict_rules
		WHERE mod_id_a = ? OR mod_id_b = ?
	`

	rows, err := r.db.Query(query, modID, modID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conflict rules: %w", err)
	}
	defer rows.Close()

	rules := make([]*models.ConflictRule, 0)
	for rows.Next() {
		var rule models.ConflictRule
		var description sql.NullString
		if err := rows.Scan(
			&rule.ID,
			&rule.ModIDA,
			&rule.ModIDB,
			&rule.ConflictType,
			&description,
			&rule.Severity,
		); err != nil {
			return nil, fmt.Errorf("failed to scan conflict rule: %w", err)
		}
		rule.Description = description.String
		rules = append(rules, &rule)
	}

	return rules, rows.Err()
}

// Save saves or updates a conflict rule
func (r *ConflictRuleRepository) Save(rule *models.ConflictRule) error {
	if rule.ID == 0 {
		// Insert new rule
		query := `
			INSERT INTO conflict_rules (mod_id_a, mod_id_b, conflict_type, description, severity)
			VALUES (?, ?, ?, ?, ?)
		`
		result, err := r.db.Exec(query, rule.ModIDA, rule.ModIDB, rule.ConflictType, nullString(rule.Description), rule.Severity)
		if err != nil {
			return fmt.Errorf("failed to insert conflict rule: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		rule.ID = id
	} else {
		// Update existing rule
		query := `
			UPDATE conflict_rules
			SET mod_id_a = ?, mod_id_b = ?, conflict_type = ?, description = ?, severity = ?
			WHERE id = ?
		`
		_, err := r.db.Exec(query, rule.ModIDA, rule.ModIDB, rule.ConflictType, nullString(rule.Description), rule.Severity, rule.ID)
		if err != nil {
			return fmt.Errorf("failed to update conflict rule: %w", err)
		}
	}

	return nil
}

// Delete deletes a conflict rule by ID
func (r *ConflictRuleRepository) Delete(id int64) error {
	query := "DELETE FROM conflict_rules WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete conflict rule: %w", err)
	}
	return nil
}

// GetAll retrieves all conflict rules
func (r *ConflictRuleRepository) GetAll() ([]*models.ConflictRule, error) {
	query := `
		SELECT id, mod_id_a, mod_id_b, conflict_type, description, severity
		FROM conflict_rules
		ORDER BY mod_id_a, mod_id_b
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all conflict rules: %w", err)
	}
	defer rows.Close()

	rules := make([]*models.ConflictRule, 0)
	for rows.Next() {
		var rule models.ConflictRule
		var description sql.NullString
		if err := rows.Scan(
			&rule.ID,
			&rule.ModIDA,
			&rule.ModIDB,
			&rule.ConflictType,
			&description,
			&rule.Severity,
		); err != nil {
			return nil, fmt.Errorf("failed to scan conflict rule: %w", err)
		}
		rule.Description = description.String
		rules = append(rules, &rule)
	}

	return rules, rows.Err()
}

// CheckConflicts checks if any of the given mods have conflicts with each other
func (r *ConflictRuleRepository) CheckConflicts(modIDs []string) ([]*models.ConflictRule, error) {
	if len(modIDs) < 2 {
		return []*models.ConflictRule{}, nil
	}

	// Build query with placeholders
	placeholders := ""
	args := make([]interface{}, 0)
	for _, modID := range modIDs {
		if placeholders != "" {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, modID)
	}

	query := fmt.Sprintf(`
		SELECT id, mod_id_a, mod_id_b, conflict_type, description, severity
		FROM conflict_rules
		WHERE mod_id_a IN (%s) AND mod_id_b IN (%s)
	`, placeholders, placeholders)

	// Add args twice (for both IN clauses)
	args = append(args, args...)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to check conflicts: %w", err)
	}
	defer rows.Close()

	rules := make([]*models.ConflictRule, 0)
	for rows.Next() {
		var rule models.ConflictRule
		var description sql.NullString
		if err := rows.Scan(
			&rule.ID,
			&rule.ModIDA,
			&rule.ModIDB,
			&rule.ConflictType,
			&description,
			&rule.Severity,
		); err != nil {
			return nil, fmt.Errorf("failed to scan conflict rule: %w", err)
		}
		rule.Description = description.String
		rules = append(rules, &rule)
	}

	return rules, rows.Err()
}
