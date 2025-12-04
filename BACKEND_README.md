# ModpackGraph Backend Implementation

## Overview

This implementation provides a complete backend architecture for modpack analysis with the following features:

- **Multi-loader support**: Fabric, Forge (Modern & Legacy), NeoForge
- **Hash-based caching**: SHA-256 hashing for efficient metadata caching
- **Dependency analysis**: Graph-based dependency resolution with conflict detection
- **SQLite persistence**: Local database for mod metadata and modpack state
- **Dependency injection**: Uber FX for clean service architecture

## Architecture

### Domain Models (`internal/models`)

Core entities representing the domain:
- `Mod` & `ModMetadata` - Mod information with hash, dependencies, loader type
- `Dependency` & `VersionRange` - Dependency relationships with version constraints
- `Modpack` & `ModpackMod` - Modpack configuration and mod associations
- `Graph`, `Node`, `Edge` - Dependency graph for visualization
- `Conflict` - Conflict detection results (Phase 2)
- `LoaderType` & `Environment` - Type-safe enums

### Loader Services (`internal/services/loaders`)

Modular metadata extraction for different mod loaders:
- `FabricLoader` - Parses `fabric.mod.json`
- `ForgeModernLoader` - Parses `META-INF/mods.toml` (Forge 1.13+)
- `ForgeLegacyLoader` - Parses `mcmod.info` (Forge 1.12.2 and earlier)
- `NeoForgeLoader` - Handles NeoForge mods (uses modern Forge format)
- `LoaderRegistry` - Auto-detects correct loader for each JAR
- `IconExtractor` - Shared utility for icon extraction and base64 encoding

### Repository Layer (`internal/repository`)

SQLite-based data persistence:
- `DB` - Database connection and schema initialization
- `ModRepository` - CRUD operations for mod metadata
- `ModpackRepository` - Modpack and modpack_mods management
- `ConflictRuleRepository` - Known incompatibilities database

### Core Services (`internal/services`)

Business logic layer:
- `MetadataService` - JAR processing and hash computation
- `CacheService` - Hash-based caching with cache hit/miss tracking
- `ScanService` - Directory scanning and change detection
- `DependencyService` - Dependency graph analysis and conflict detection
- `ConflictService` - Phase 2 conflict detection (missing deps, version conflicts, known incompatibilities)
- `AnalysisService` - Orchestrates full analysis pipeline

### Dependency Injection (`internal/di`)

Uber FX module for wiring all services together with automatic dependency resolution.

### Application Layer (`internal/app`)

Wails-exposed API:
- `ScanModpack(path)` - Scan modpack with hash computation
- `AnalyzeModpack(path)` - Run full two-phase analysis
- `GetDependencyGraph(path)` - Retrieve dependency graph for visualization
- `GetModpackStatus(path)` - Check cache status and changes
- `GetModMetadata(modID)` - Retrieve cached mod info
- `RefreshModpack(path)` - Force re-scan
- Progress events via Wails runtime

## Database Schema

### mods table
Cached mod metadata indexed by SHA-256 hash:
- id, hash, version, name, description
- loader_type, environment, icon_data
- metadata_json (full metadata as JSON)
- created_at, updated_at

### mod_dependencies table
Dependency relationships:
- mod_id, dependency_id, required, version_range

### modpacks table
Analyzed modpacks:
- id, name, path, last_scanned, mod_count

### modpack_mods table
Mods in each modpack (with hash for change detection):
- modpack_id, mod_id, hash, file_path

### conflict_rules table
Known incompatibilities:
- id, mod_id_a, mod_id_b, conflict_type, description, severity

## Caching Strategy

1. **Initial Scan**: Compute SHA-256 hash for each JAR
2. **Cache Lookup**: Query mods table by hash
3. **Selective Loading**: Only extract metadata from new/changed JARs
4. **Incremental Updates**: Update modpack_mods with current state
5. **Change Detection**: Track new, modified, removed mods

## Two-Phase Analysis

### Phase 1: Mod Feature Analysis (DEFERRED)
Feature extraction is stubbed for future implementation. Will require:
- Bytecode analysis (ASM library)
- Registry pattern detection
- Mixin/event handler analysis

### Phase 2: Conflict Detection (IMPLEMENTED)
Uses dependency metadata to identify:
- **Missing dependencies** - Required mods not present
- **Version conflicts** - Incompatible version ranges
- **Circular dependencies** - Dependency cycles
- **Known incompatibilities** - From conflict_rules table
- **Environment mismatches** - Client-only on server, etc.

## Usage Example

```go
// Scan a modpack
result, err := app.ScanModpack("/path/to/modpack/mods")

// Run full analysis
report, err := app.AnalyzeModpack("/path/to/modpack/mods")

// Get dependency graph for visualization
graph, err := app.GetDependencyGraph("/path/to/modpack/mods")
```

## Dependencies

- **github.com/Masterminds/semver/v3** - Version constraint handling
- **github.com/pelletier/go-toml/v2** - TOML parsing for Forge mods
- **github.com/mattn/go-sqlite3** - SQLite driver
- **go.uber.org/fx** - Dependency injection
- **github.com/sirupsen/logrus** - Structured logging
- **github.com/wailsapp/wails/v2** - Desktop app framework

## Testing Strategy

### Unit Tests (To Be Implemented)
- Mock loaders for metadata extraction tests
- In-memory SQLite for repository tests
- Dependency service tests with various scenarios

### Integration Tests (To Be Implemented)
- Sample modpacks with real/synthetic JARs
- Known-good analysis results (golden files)
- Performance benchmarks for large modpacks (500+ mods)

## Future Enhancements

1. **Phase 1 Feature Extraction**
   - JAR introspection and bytecode analysis
   - Registry detection (items, blocks, entities)
   - Configuration parsing

2. **Advanced Conflict Detection**
   - Feature overlap detection
   - Performance impact analysis
   - Automatic conflict resolution suggestions

3. **Performance Optimizations**
   - Parallel JAR processing
   - Progressive caching
   - Index optimizations

4. **API Extensions**
   - Batch operations
   - Webhook notifications
   - Export/import configurations

