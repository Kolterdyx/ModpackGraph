package services

import (
	"ModpackGraph/internal/models"
)

// DependencyService handles dependency resolution and graph generation
type DependencyService struct{}

// NewDependencyService creates a new DependencyService
func NewDependencyService() *DependencyService {
	return &DependencyService{}
}

// DependencyResult represents the result of dependency analysis
type DependencyResult struct {
	Graph               *models.Graph
	MissingDependencies []*MissingDependency
	VersionConflicts    []*VersionConflict
	CircularDeps        [][]string
}

// MissingDependency represents a missing dependency
type MissingDependency struct {
	ModID        string
	ModName      string
	DependencyID string
	Required     bool
	VersionRange string
}

// VersionConflict represents a version conflict
type VersionConflict struct {
	ModID          string
	ModName        string
	DependencyID   string
	DependencyName string
	RequiredRange  string
	ActualVersion  string
}

// Analyze analyzes dependencies for a set of mods
func (s *DependencyService) Analyze(mods []*models.ModMetadata) (*DependencyResult, error) {
	result := &DependencyResult{
		Graph:               models.NewGraph(),
		MissingDependencies: make([]*MissingDependency, 0),
		VersionConflicts:    make([]*VersionConflict, 0),
		CircularDeps:        make([][]string, 0),
	}

	// Create mod map for quick lookup
	modMap := make(map[string]*models.ModMetadata)
	for _, mod := range mods {
		modMap[mod.ID] = mod
	}

	// Build graph nodes
	for _, mod := range mods {
		result.Graph.AddNode(mod.ID, mod.Name, mod.Version)
	}

	// Build graph edges and check for issues
	for _, mod := range mods {
		for _, dep := range mod.Dependencies {
			// Add edge to graph
			versionLabel := ""
			if dep.VersionRange != nil {
				versionLabel = dep.VersionRange.String()
			}
			result.Graph.AddEdge(mod.ID, dep.DependencyID, dep.Required, versionLabel)

			// Check if dependency exists
			depMod, exists := modMap[dep.DependencyID]
			if !exists {
				// Missing dependency
				result.MissingDependencies = append(result.MissingDependencies, &MissingDependency{
					ModID:        mod.ID,
					ModName:      mod.Name,
					DependencyID: dep.DependencyID,
					Required:     dep.Required,
					VersionRange: versionLabel,
				})
				continue
			}

			// Check version compatibility
			if dep.VersionRange != nil && !dep.VersionRange.Check(depMod.Version) {
				result.VersionConflicts = append(result.VersionConflicts, &VersionConflict{
					ModID:          mod.ID,
					ModName:        mod.Name,
					DependencyID:   dep.DependencyID,
					DependencyName: depMod.Name,
					RequiredRange:  versionLabel,
					ActualVersion:  depMod.Version,
				})
			}
		}
	}

	// Detect circular dependencies
	result.CircularDeps = s.detectCircularDependencies(mods, modMap)

	return result, nil
}

// detectCircularDependencies detects circular dependency chains
func (s *DependencyService) detectCircularDependencies(mods []*models.ModMetadata, modMap map[string]*models.ModMetadata) [][]string {
	cycles := make([][]string, 0)
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := make([]string, 0)

	var detectCycle func(modID string) bool
	detectCycle = func(modID string) bool {
		visited[modID] = true
		recStack[modID] = true
		path = append(path, modID)

		mod, exists := modMap[modID]
		if !exists {
			// Remove from path and recStack before returning
			path = path[:len(path)-1]
			recStack[modID] = false
			return false
		}

		for _, dep := range mod.Dependencies {
			if !visited[dep.DependencyID] {
				if detectCycle(dep.DependencyID) {
					return true
				}
			} else if recStack[dep.DependencyID] {
				// Found a cycle - extract the cycle from path
				cycleStart := -1
				for i, id := range path {
					if id == dep.DependencyID {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					cycle := make([]string, len(path)-cycleStart)
					copy(cycle, path[cycleStart:])
					cycles = append(cycles, cycle)
				}
				return true
			}
		}

		path = path[:len(path)-1]
		recStack[modID] = false
		return false
	}

	for _, mod := range mods {
		if !visited[mod.ID] {
			detectCycle(mod.ID)
		}
	}

	return cycles
}

// BuildGraph builds a dependency graph for visualization
func (s *DependencyService) BuildGraph(mods []*models.ModMetadata) *models.Graph {
	graph := models.NewGraph()

	// Add nodes
	for _, mod := range mods {
		graph.AddNode(mod.ID, mod.Name, mod.Version)
	}

	// Add edges
	for _, mod := range mods {
		for _, dep := range mod.Dependencies {
			versionLabel := ""
			if dep.VersionRange != nil {
				versionLabel = dep.VersionRange.String()
			}
			graph.AddEdge(mod.ID, dep.DependencyID, dep.Required, versionLabel)
		}
	}

	return graph
}

// GetDependents returns all mods that depend on a given mod
func (s *DependencyService) GetDependents(modID string, mods []*models.ModMetadata) []*models.ModMetadata {
	dependents := make([]*models.ModMetadata, 0)

	for _, mod := range mods {
		for _, dep := range mod.Dependencies {
			if dep.DependencyID == modID {
				dependents = append(dependents, mod)
				break
			}
		}
	}

	return dependents
}

// GetDependencies returns all dependencies of a given mod
func (s *DependencyService) GetDependencies(modID string, mods []*models.ModMetadata) []*models.ModMetadata {
	// Find the mod
	var targetMod *models.ModMetadata
	modMap := make(map[string]*models.ModMetadata)
	for _, mod := range mods {
		modMap[mod.ID] = mod
		if mod.ID == modID {
			targetMod = mod
		}
	}

	if targetMod == nil {
		return []*models.ModMetadata{}
	}

	dependencies := make([]*models.ModMetadata, 0)
	for _, dep := range targetMod.Dependencies {
		if depMod, exists := modMap[dep.DependencyID]; exists {
			dependencies = append(dependencies, depMod)
		}
	}

	return dependencies
}
