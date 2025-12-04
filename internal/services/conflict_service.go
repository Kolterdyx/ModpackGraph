package services

import (
	"ModpackGraph/internal/models"
	"ModpackGraph/internal/repository"
	"fmt"
)

// ConflictService handles Phase 2 conflict detection
type ConflictService struct {
	conflictRuleRepo  *repository.ConflictRuleRepository
	dependencyService *DependencyService
}

// NewConflictService creates a new ConflictService
func NewConflictService(conflictRuleRepo *repository.ConflictRuleRepository, dependencyService *DependencyService) *ConflictService {
	return &ConflictService{
		conflictRuleRepo:  conflictRuleRepo,
		dependencyService: dependencyService,
	}
}

// DetectConflicts detects conflicts in a set of mods
func (s *ConflictService) DetectConflicts(mods []*models.ModMetadata, depResult *DependencyResult) ([]*models.Conflict, error) {
	conflicts := make([]*models.Conflict, 0)

	// 1. Check for missing dependencies
	for _, missing := range depResult.MissingDependencies {
		severity := models.ConflictSeverityCritical
		if !missing.Required {
			severity = models.ConflictSeverityWarning
		}

		description := fmt.Sprintf("Mod '%s' requires '%s' which is not installed", missing.ModName, missing.DependencyID)
		if missing.VersionRange != "" {
			description += fmt.Sprintf(" (version: %s)", missing.VersionRange)
		}

		conflict := models.NewConflict(
			models.ConflictTypeMissingDependency,
			severity,
			description,
			[]string{missing.ModID, missing.DependencyID},
		)
		conflict.Details["required"] = missing.Required
		conflict.Details["version_range"] = missing.VersionRange

		conflicts = append(conflicts, conflict)
	}

	// 2. Check for version conflicts
	for _, vc := range depResult.VersionConflicts {
		description := fmt.Sprintf("Mod '%s' requires '%s' version %s, but version %s is installed",
			vc.ModName, vc.DependencyName, vc.RequiredRange, vc.ActualVersion)

		conflict := models.NewConflict(
			models.ConflictTypeVersionConflict,
			models.ConflictSeverityWarning,
			description,
			[]string{vc.ModID, vc.DependencyID},
		)
		conflict.Details["required_range"] = vc.RequiredRange
		conflict.Details["actual_version"] = vc.ActualVersion

		conflicts = append(conflicts, conflict)
	}

	// 3. Check for circular dependencies
	for _, cycle := range depResult.CircularDeps {
		if len(cycle) > 0 {
			description := fmt.Sprintf("Circular dependency detected: %v", cycle)

			conflict := models.NewConflict(
				models.ConflictTypeCircularDependency,
				models.ConflictSeverityCritical,
				description,
				cycle,
			)

			conflicts = append(conflicts, conflict)
		}
	}

	// 4. Check known incompatibilities from database
	modIDs := make([]string, len(mods))
	for i, mod := range mods {
		modIDs[i] = mod.ID
	}

	knownConflicts, err := s.conflictRuleRepo.CheckConflicts(modIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to check known conflicts: %w", err)
	}

	for _, rule := range knownConflicts {
		conflict := models.NewConflict(
			rule.ConflictType,
			rule.Severity,
			rule.Description,
			[]string{rule.ModIDA, rule.ModIDB},
		)
		conflict.Details["rule_id"] = rule.ID

		conflicts = append(conflicts, conflict)
	}

	// 5. Check for environment mismatches (if server/client info is available)
	envConflicts := s.checkEnvironmentMismatches(mods)
	conflicts = append(conflicts, envConflicts...)

	return conflicts, nil
}

// checkEnvironmentMismatches checks for environment compatibility issues
func (s *ConflictService) checkEnvironmentMismatches(mods []*models.ModMetadata) []*models.Conflict {
	conflicts := make([]*models.Conflict, 0)

	// Check if we have any server-only or client-only mods
	hasClientOnly := false
	hasServerOnly := false
	clientOnlyMods := make([]string, 0)
	serverOnlyMods := make([]string, 0)

	for _, mod := range mods {
		switch mod.Environment {
		case models.EnvironmentClient:
			hasClientOnly = true
			clientOnlyMods = append(clientOnlyMods, mod.ID)
		case models.EnvironmentServer:
			hasServerOnly = true
			serverOnlyMods = append(serverOnlyMods, mod.ID)
		}
	}

	// If we have both client-only and server-only mods, that's unusual
	if hasClientOnly && hasServerOnly {
		description := fmt.Sprintf("Modpack contains both client-only and server-only mods. This may indicate a configuration issue.")

		conflict := models.NewConflict(
			models.ConflictTypeEnvironmentMismatch,
			models.ConflictSeverityInfo,
			description,
			append(clientOnlyMods, serverOnlyMods...),
		)
		conflict.Details["client_only_mods"] = clientOnlyMods
		conflict.Details["server_only_mods"] = serverOnlyMods

		conflicts = append(conflicts, conflict)
	}

	return conflicts
}

// AddConflictRule adds a new conflict rule to the database
func (s *ConflictService) AddConflictRule(modIDA, modIDB string, conflictType models.ConflictType, description string, severity models.ConflictSeverity) error {
	rule := models.NewConflictRule(modIDA, modIDB, conflictType, description, severity)
	return s.conflictRuleRepo.Save(rule)
}

// GetConflictRules retrieves all conflict rules
func (s *ConflictService) GetConflictRules() ([]*models.ConflictRule, error) {
	return s.conflictRuleRepo.GetAll()
}

// DeleteConflictRule deletes a conflict rule
func (s *ConflictService) DeleteConflictRule(id int64) error {
	return s.conflictRuleRepo.Delete(id)
}
