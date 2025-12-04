package app

import (
	"ModpackGraph/internal/models"
	"ModpackGraph/internal/services"
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct holds the application state and services
type App struct {
	ctx               context.Context
	analysisService   *services.AnalysisService
	scanService       *services.ScanService
	dependencyService *services.DependencyService
	conflictService   *services.ConflictService
	cacheService      *services.CacheService
}

// NewApp creates a new App instance with injected dependencies
func NewApp(
	analysisService *services.AnalysisService,
	scanService *services.ScanService,
	dependencyService *services.DependencyService,
	conflictService *services.ConflictService,
	cacheService *services.CacheService,
) *App {
	return &App{
		analysisService:   analysisService,
		scanService:       scanService,
		dependencyService: dependencyService,
		conflictService:   conflictService,
		cacheService:      cacheService,
	}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// ScanModpack performs an initial scan with hash computation and cache lookup
func (a *App) ScanModpack(path string) (*services.ScanResult, error) {
	a.emitProgress("scan", "Starting modpack scan...", 0)

	result, err := a.scanService.ScanModpack(path)
	if err != nil {
		a.emitProgress("scan", fmt.Sprintf("Error: %v", err), 100)
		return nil, err
	}

	a.emitProgress("scan", "Scan complete", 100)
	return result, nil
}

// AnalyzeModpack runs the full two-phase analysis pipeline
func (a *App) AnalyzeModpack(path string) (*services.AnalysisReport, error) {
	a.emitProgress("analyze", "Starting analysis...", 0)

	report, err := a.analysisService.AnalyzeModpack(path)
	if err != nil {
		a.emitProgress("analyze", fmt.Sprintf("Error: %v", err), 100)
		return nil, err
	}

	a.emitProgress("analyze", "Analysis complete", 100)
	return report, nil
}

// GetDependencyGraph retrieves the dependency graph for visualization
func (a *App) GetDependencyGraph(modpackPath string) (*models.Graph, error) {
	return a.analysisService.GetDependencyGraph(modpackPath)
}

// GetModpackStatus checks cache status and change detection
func (a *App) GetModpackStatus(path string) (*services.ScanResult, error) {
	return a.analysisService.GetModpackStatus(path)
}

// GetModMetadata retrieves cached mod information
func (a *App) GetModMetadata(modID string) (*models.ModMetadata, error) {
	return a.cacheService.GetByID(modID)
}

// RefreshModpack forces a re-scan of all mods
func (a *App) RefreshModpack(path string) (*services.ScanResult, error) {
	a.emitProgress("refresh", "Refreshing modpack...", 0)

	result, err := a.scanService.ScanModpack(path)
	if err != nil {
		a.emitProgress("refresh", fmt.Sprintf("Error: %v", err), 100)
		return nil, err
	}

	a.emitProgress("refresh", "Refresh complete", 100)
	return result, nil
}

// GetConflictRules retrieves all known conflict rules
func (a *App) GetConflictRules() ([]*models.ConflictRule, error) {
	return a.conflictService.GetConflictRules()
}

// AddConflictRule adds a new conflict rule
func (a *App) AddConflictRule(modIDA, modIDB string, conflictType, description, severity string) error {
	return a.conflictService.AddConflictRule(
		modIDA,
		modIDB,
		models.ConflictType(conflictType),
		description,
		models.ConflictSeverity(severity),
	)
}

// DeleteConflictRule deletes a conflict rule
func (a *App) DeleteConflictRule(id int64) error {
	return a.conflictService.DeleteConflictRule(id)
}

// QuickScan performs a scan without full dependency analysis
func (a *App) QuickScan(path string) (*services.ScanResult, error) {
	a.emitProgress("quick_scan", "Quick scanning...", 0)

	result, err := a.analysisService.QuickScan(path)
	if err != nil {
		a.emitProgress("quick_scan", fmt.Sprintf("Error: %v", err), 100)
		return nil, err
	}

	a.emitProgress("quick_scan", "Quick scan complete", 100)
	return result, nil
}

// emitProgress emits progress events to the frontend
func (a *App) emitProgress(operation, message string, progress int) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "progress", map[string]interface{}{
			"operation": operation,
			"message":   message,
			"progress":  progress,
		})
	}
}
