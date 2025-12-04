package services

import (
	"ModpackGraph/internal/models"
	"fmt"
)

// AnalysisService orchestrates the full analysis pipeline
type AnalysisService struct {
	scanService       *ScanService
	dependencyService *DependencyService
	conflictService   *ConflictService
}

// NewAnalysisService creates a new AnalysisService
func NewAnalysisService(scanService *ScanService, dependencyService *DependencyService, conflictService *ConflictService) *AnalysisService {
	return &AnalysisService{
		scanService:       scanService,
		dependencyService: dependencyService,
		conflictService:   conflictService,
	}
}

// AnalysisReport represents the complete analysis results
type AnalysisReport struct {
	ScanResult       *ScanResult
	DependencyResult *DependencyResult
	Conflicts        []*models.Conflict
	Summary          *AnalysisSummary
}

// AnalysisSummary provides high-level statistics
type AnalysisSummary struct {
	TotalMods            int
	NewMods              int
	UpdatedMods          int
	RemovedMods          int
	CacheHitRate         float64
	MissingDependencies  int
	VersionConflicts     int
	CircularDependencies int
	TotalConflicts       int
	CriticalConflicts    int
	WarningConflicts     int
	InfoConflicts        int
}

// AnalyzeModpack runs the complete analysis pipeline
func (s *AnalysisService) AnalyzeModpack(modpackPath string) (*AnalysisReport, error) {
	// Step 1: Scan modpack
	scanResult, err := s.scanService.ScanModpack(modpackPath)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	// Step 2: Analyze dependencies
	depResult, err := s.dependencyService.Analyze(scanResult.Mods)
	if err != nil {
		return nil, fmt.Errorf("dependency analysis failed: %w", err)
	}

	// Step 3: Detect conflicts
	conflicts, err := s.conflictService.DetectConflicts(scanResult.Mods, depResult)
	if err != nil {
		return nil, fmt.Errorf("conflict detection failed: %w", err)
	}

	// Generate summary
	summary := s.generateSummary(scanResult, depResult, conflicts)

	return &AnalysisReport{
		ScanResult:       scanResult,
		DependencyResult: depResult,
		Conflicts:        conflicts,
		Summary:          summary,
	}, nil
}

// generateSummary generates analysis summary statistics
func (s *AnalysisService) generateSummary(scanResult *ScanResult, depResult *DependencyResult, conflicts []*models.Conflict) *AnalysisSummary {
	summary := &AnalysisSummary{
		TotalMods:            len(scanResult.Mods),
		NewMods:              len(scanResult.NewMods),
		UpdatedMods:          len(scanResult.UpdatedMods),
		RemovedMods:          len(scanResult.RemovedMods),
		MissingDependencies:  len(depResult.MissingDependencies),
		VersionConflicts:     len(depResult.VersionConflicts),
		CircularDependencies: len(depResult.CircularDeps),
		TotalConflicts:       len(conflicts),
	}

	// Calculate cache hit rate
	total := scanResult.CacheHits + scanResult.CacheMisses
	if total > 0 {
		summary.CacheHitRate = float64(scanResult.CacheHits) / float64(total) * 100.0
	}

	// Count conflicts by severity
	for _, conflict := range conflicts {
		switch conflict.Severity {
		case models.ConflictSeverityCritical:
			summary.CriticalConflicts++
		case models.ConflictSeverityWarning:
			summary.WarningConflicts++
		case models.ConflictSeverityInfo:
			summary.InfoConflicts++
		}
	}

	return summary
}

// GetDependencyGraph returns the dependency graph for a modpack
func (s *AnalysisService) GetDependencyGraph(modpackPath string) (*models.Graph, error) {
	// Get modpack
	modpack, err := s.scanService.GetModpackByPath(modpackPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack: %w", err)
	}

	if modpack == nil {
		return nil, fmt.Errorf("modpack not found: %s", modpackPath)
	}

	// Get mods
	mods, err := s.scanService.GetModpackMods(modpack.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mods: %w", err)
	}

	// Build graph
	graph := s.dependencyService.BuildGraph(mods)
	return graph, nil
}

// QuickScan performs a scan without full dependency analysis
func (s *AnalysisService) QuickScan(modpackPath string) (*ScanResult, error) {
	return s.scanService.ScanModpack(modpackPath)
}

// GetModpackStatus checks the cache status of a modpack
func (s *AnalysisService) GetModpackStatus(modpackPath string) (*ScanResult, error) {
	// This could be optimized to not re-extract metadata, just check hashes
	return s.scanService.ScanModpack(modpackPath)
}
