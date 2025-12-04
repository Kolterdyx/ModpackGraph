package di

import (
	"ModpackGraph/internal/repository"
	"ModpackGraph/internal/services"
	"ModpackGraph/internal/services/loaders"
	"os"
	"path/filepath"

	"go.uber.org/fx"
)

// Module provides all application dependencies
var Module = fx.Options(
	// Database
	fx.Provide(NewDatabase),

	// Repositories
	fx.Provide(repository.NewModRepository),
	fx.Provide(repository.NewModpackRepository),
	fx.Provide(repository.NewConflictRuleRepository),

	// Loader services
	fx.Provide(
		loaders.NewLoaderRegistry,
		loaders.NewIconExtractor,
		fx.Annotate(loaders.NewFabricLoader, fx.ResultTags(`group:"mod_loader"`)),
		fx.Annotate(loaders.NewForgeModernLoader, fx.ResultTags(`group:"mod_loader"`)),
		fx.Annotate(loaders.NewForgeLegacyLoader, fx.ResultTags(`group:"mod_loader"`)),
		fx.Annotate(loaders.NewNeoForgeLoader, fx.ResultTags(`group:"mod_loader"`)),
	),

	// Core services
	fx.Provide(services.NewMetadataService),
	fx.Provide(services.NewCacheService),
	fx.Provide(services.NewScanService),
	fx.Provide(services.NewDependencyService),
	fx.Provide(services.NewConflictService),
	fx.Provide(services.NewAnalysisService),
)

// NewDatabase creates and initializes the database
func NewDatabase() (*repository.DB, error) {
	// Get user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}

	dbDir := filepath.Join(configDir, "ModpackGraph")
	dbPath := filepath.Join(dbDir, "modpack_graph.db")

	return repository.NewDB(dbPath)
}
