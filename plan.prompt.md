# Plan: Backend Architecture Redesign for Modpack Analysis

Design and implement a modular, service-oriented architecture from scratch using dependency injection (Uber FX), with clear separation of concerns for mod metadata extraction, dependency resolution, conflict detection, and data persistence.

## Required Dependencies

Install the following Go packages:

```bash
# Dependency Injection
go get go.uber.org/fx

# Logging
go get github.com/sirupsen/logrus

# Database
go get github.com/mattn/go-sqlite3
go get github.com/jmoiron/sqlx  # SQL extensions

# Testing
go get github.com/stretchr/testify
go install github.com/vektra/mockery/v2@latest  # Code generation tool
```

## Steps

1. **Create domain models** - Design `internal/models` package with core entities:
   - `Mod` - Basic mod identity (ID, version)
   - `ModMetadata` - Full mod information with hash, name, loader type, environment, icon, dependencies
   - `Dependency` - Dependency relationship with version compatibility
   - `VersionRange` - Version constraint representation (replaces `Compat`)
   - `Modpack` - Modpack entity with path, mods list, scan timestamp
   - `Graph`, `Node`, `Edge` - Dependency graph visualization models
   - `Conflict` - Phase 2 conflict results (type, severity, description, affected mods)
   - `ModFeature` - Phase 1 feature data (future use, stub for now)
   - `LoaderType` - Enum for Fabric, ForgeModern, ForgeLegacy, NeoForge
   - `Environment` - Enum for Client, Server, Both

2. **Create loader services** - Implement `internal/services/loaders` with mod metadata extraction for each loader type:
   
   **Core Interface**:
   ```go
   type ModLoader interface {
       // Returns true if this loader can handle the given JAR
       CanHandle(zipReader *zip.Reader) bool
       // Extract metadata from the JAR
       ExtractMetadata(zipReader *zip.Reader, jarPath string) (*models.ModMetadata, error)
   }
   ```
   
   **Implementations**:
   - `FabricLoader` - Parse `fabric.mod.json`:
     - JSON unmarshaling for mod metadata
     - Dependency extraction from `depends`, `recommends`, `suggests` fields
     - Environment field mapping (client/server/*)
     - Icon path detection from `icon` field
   
   - `ForgeModernLoader` - Parse `META-INF/mods.toml`:
     - TOML unmarshaling (github.com/pelletier/go-toml/v2)
     - Handle `${file.jarVersion}` placeholders via MANIFEST.MF
     - Dependency extraction from `dependencies` section
     - Icon path detection from `logoFile` field
     - Default to "both" environment
   
   - `ForgeLegacyLoader` - Parse `mcmod.info`:
     - JSON unmarshaling of array format
     - Handle `${version}` placeholders via MANIFEST.MF
     - Dependencies from `requiredMods` and `dependencies` arrays
     - Default to "both" environment
   
   - `NeoForgeLoader` - Initially alias to ForgeModernLoader (format similarity)
   
   - `LoaderRegistry` - Manages all loaders:
     - Register all loader implementations
     - Detect correct loader for JAR via `CanHandle` checks
     - Return appropriate loader or error
   
   - `IconExtractor` - Shared utility:
     - Search common paths (logo.png, icon.png, pack.png, assets/{modid}/icon.png)
     - Extract and base64 encode
     - Provide default icon fallback

3. **Setup database layer** - Create `internal/repository` package with:
   - SQLite initialization and schema migrations (using `golang-migrate` or embedded SQL)
   - Database models matching schema (mods, modpacks, mod_features, etc.)
   - Repository pattern implementations for data access
   - Transaction management utilities

4. **Implement core services** - Establish `internal/services` with:
   
   - `MetadataService` - JAR processing:
     - Open JAR files and create zip.Reader
     - Compute SHA-256 hash of JAR file
     - Delegate to LoaderRegistry for metadata extraction
     - Handle embedded JARs (recursive scanning)
     - Return ModMetadata with hash and loader type
   
   - `CacheService` - Hash-based caching:
     - Check if hash exists in repository
     - Return cached metadata on hit
     - Coordinate with MetadataService on miss
     - Track which mods are new/modified/unchanged
     - Batch update operations for performance
   
   - `ScanService` - Modpack scanning:
     - Walk directory tree for .jar files
     - Coordinate with CacheService for each JAR
     - Build list of all mods in modpack
     - Detect removed mods (in DB but not in scan)
     - Update modpack_mods table
   
   - `DependencyService` - Dependency resolution:
     - Build dependency graph from mod metadata
     - Topological sort to detect cycles
     - Check version compatibility using VersionRange
     - Identify missing dependencies
     - Flag incompatible versions
     - Generate dependency graph for visualization
   
   - `ConflictService` - Phase 2 conflict detection:
     - Check for missing dependencies (from DependencyService)
     - Query mod_features for overlapping features (future)
     - Match against conflict_rules table
     - Check environment mismatches (client-only on server)
     - Classify conflicts by severity (critical/warning/info)
     - Generate Conflict models with descriptions
   
   - `FeatureService` - Phase 1 feature extraction (STUB):
     - Interface defined but not implemented
     - Placeholder that returns empty feature sets
     - TODO comment referencing future bytecode analysis plan
   
   - `AnalysisService` - Orchestration:
     - ScanModpack: Coordinate ScanService + CacheService
     - AnalyzeDependencies: Call DependencyService
     - DetectConflicts: Call ConflictService with dependency results
     - Full pipeline: Scan → Dependencies → Conflicts
     - Emit progress events via callback/channel
   
   - `RepositoryService` - Data persistence:
     - CRUD operations for all tables
     - Transaction management
     - Query builders for complex lookups
     - Connection pooling
     - Migration execution

5. **Setup dependency injection** - Configure Uber FX for service wiring:
   - Create `internal/di` package with FX module definitions
   - Provide constructors for all services (loaders, metadata, cache, dependency, conflict, analysis, repository)
   - Inject logger into all services
   - Create `App` constructor that receives all required services via FX
   - Update `main.go` to use FX application lifecycle instead of direct instantiation
   - Maintain Wails `Startup` hook integration

6. **Add logging and observability** - Integrate Logrus throughout services:
   - Configure logger in `main.go` with appropriate level and format
   - Add structured logging with fields (modID, path, hash, operation)
   - Log key operations: cache hits/misses, metadata extraction, dependency resolution
   - Create logging middleware for service layer
   - Emit progress events to Wails runtime for frontend display

7. **Design Wails API layer** - Create new `App` struct in `internal/app` with injected services:
   - `ScanModpack(path string)` - Initial scan with hash computation and cache lookup
   - `AnalyzeModpack(path string)` - Run two-phase analysis with progress events
   - `GetDependencyGraph(modpackPath string)` - Retrieve dependency graph for visualization
   - `GetModpackStatus(path string)` - Check cache status and change detection results
   - `GetConflicts(modpackPath string)` - Query Phase 2 conflict detection results
   - `GetModMetadata(modID string)` - Retrieve cached mod information
   - `RefreshModpack(path string)` - Force re-scan of all mods
   - Emit Wails runtime events for progress tracking during long operations

## Database Design

### Schema

**mods table** - Cached mod metadata indexed by hash
- `id` (TEXT PRIMARY KEY) - Mod ID from metadata
- `hash` (TEXT UNIQUE) - SHA-256 hash of JAR file
- `version` (TEXT) - Mod version
- `name` (TEXT) - Display name
- `loader_type` (TEXT) - fabric, forge_modern, forge_legacy, neoforge
- `icon_data` (TEXT) - Base64 encoded icon
- `environment` (TEXT) - client, server, or both
- `metadata_json` (TEXT) - Full metadata as JSON
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

**mod_dependencies table** - Mod dependency relationships
- `mod_id` (TEXT) - References mods(id)
- `dependency_id` (TEXT) - Dependency mod ID
- `required` (BOOLEAN) - Is dependency mandatory
- `version_range` (TEXT) - Compatible version range
- PRIMARY KEY (`mod_id`, `dependency_id`)

**modpacks table** - Analyzed modpack configurations
- `id` (INTEGER PRIMARY KEY)
- `name` (TEXT)
- `path` (TEXT UNIQUE) - Filesystem path to mods folder
- `last_scanned` (TIMESTAMP)
- `mod_count` (INTEGER)

**modpack_mods table** - Mods present in each modpack
- `modpack_id` (INTEGER) - References modpacks(id)
- `mod_id` (TEXT) - References mods(id)
- `hash` (TEXT) - Current hash for change detection
- `file_path` (TEXT) - Relative path in modpack
- PRIMARY KEY (`modpack_id`, `mod_id`)

**mod_features table** - Features provided by mods (Phase 1 analysis)
- `mod_id` (TEXT) - References mods(id)
- `feature_type` (TEXT) - world_generation, entities, items, blocks, mechanics, etc.
- `feature_data` (TEXT) - JSON with feature details
- `extracted_at` (TIMESTAMP)
- PRIMARY KEY (`mod_id`, `feature_type`)

**conflict_rules table** - Known incompatibilities database
- `id` (INTEGER PRIMARY KEY)
- `mod_id_a` (TEXT)
- `mod_id_b` (TEXT)
- `conflict_type` (TEXT) - known_incompatible, version_conflict, feature_overlap
- `description` (TEXT)
- `severity` (TEXT) - critical, warning, info

### Caching Strategy

1. **Initial scan** - Compute SHA-256 hash for each JAR file in modpack
2. **Cache lookup** - Query `mods` table by hash to find existing metadata
3. **Selective loading** - Only extract metadata from JARs with new/changed hashes
4. **Incremental updates** - Update `modpack_mods` table with current state
5. **Orphan cleanup** - Optionally remove cached mods no longer referenced by any modpack

## Two-Phase Analysis

### Phase 1: Mod Feature Analysis (Per-Mod, Cached)

**Status**: Implementation deferred to future planning phase

Extract and cache features that each mod provides:
- World generation modifications
- Entity additions/modifications
- Item registrations
- Block registrations  
- Game mechanics changes
- Configuration options

This data is stored in `mod_features` table and reused across modpacks. Only runs when:
- New mod detected (hash not in cache)
- Mod updated (hash changed)
- Feature extraction logic updated (version flag)

### Phase 2: Conflict Detection (Per-Modpack, Always Run)

Use cached feature data to identify conflicts:
1. **Dependency conflicts** - Missing dependencies, incompatible versions
2. **Feature overlap** - Multiple mods modifying same game elements
3. **Known incompatibilities** - Match against `conflict_rules` table
4. **Environment mismatches** - Client-only mods on server, etc.

Results are computed on-demand and not cached, ensuring accurate detection after modpack changes.

## Testing Strategy

### Unit Tests (Testify + Mockery)

**Loader Services**
- Mock `zip.Reader` to test each loader (Fabric, Forge Modern, Forge Legacy, NeoForge)
- Verify correct metadata extraction from sample JSON/TOML structures
- Test error handling for malformed metadata files
- Validate version range parsing in `Compat` type

**MetadataService**
- Mock loader implementations
- Test loader selection logic based on JAR contents
- Verify icon extraction and fallback behavior
- Test embedded JAR detection and recursion

**DependencyService**
- Mock `RepositoryService` for mod lookups
- Test dependency resolution with various version constraints
- Verify circular dependency detection
- Test missing dependency identification

**ConflictService** (Phase 2)
- Mock `RepositoryService` for feature data
- Test conflict detection algorithms
- Verify severity classification
- Test known incompatibility matching

**RepositoryService**
- Use in-memory SQLite (`:memory:`) for tests
- Test CRUD operations for all tables
- Verify hash-based caching logic
- Test concurrent access scenarios

### Integration Tests

**Sample Modpacks**
- Create fixture directories with real/synthetic JAR files
- Include scenarios: clean modpack, missing dependencies, version conflicts, known incompatibilities
- Test full pipeline: scan → cache → analyze → report

**Regression Tests**
- Store known-good analysis results for real modpacks (subset of mods from workspace)
- Verify consistent behavior across refactorings
- Use golden file pattern for graph outputs

**Performance Tests**
- Benchmark large modpack scanning (500+ mods)
- Measure cache hit rates
- Profile database query performance

## Implementation Notes

### Hash Computation

Use `crypto/sha256` to compute JAR file hashes:
```go
func ComputeJARHash(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()
    
    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }
    
    return hex.EncodeToString(h.Sum(nil)), nil
}
```

### Change Detection Flow

1. Scan modpack directory for `.jar` files
2. Compute hash for each JAR
3. Query `modpack_mods` table for previous hashes
4. Identify changes:
   - **New**: Hash not in `modpack_mods` for this modpack
   - **Modified**: Hash differs from stored hash
   - **Removed**: Stored entry not in current scan
   - **Unchanged**: Hash matches stored hash
5. Process only new/modified JARs through metadata extraction
6. Update `modpack_mods` table with current state

### Environment Compatibility

Extract from loader-specific metadata:

**Fabric** (`fabric.mod.json`):
```json
{
  "environment": "client" | "server" | "*"
}
```

**Forge** (annotation inspection or manifest):
- Check for `Dist.CLIENT` or `Dist.DEDICATED_SERVER` in `@Mod` annotations
- Default to "both" if not specified

### Dependency Resolution Algorithm

1. Build dependency graph from cached metadata
2. Perform topological sort to detect circular dependencies
3. Check version compatibility using `VersionRange.Intersect()`
4. Identify missing dependencies (not present in modpack)
5. Flag incompatible versions based on version ranges

### Future: Phase 1 Feature Extraction

**Deferred to future planning** - Will require:
- JAR introspection (class analysis, resource scanning)
- Registry pattern detection (item/block registrations)
- Mixin analysis for Fabric mods
- Event handler detection for Forge mods
- Configuration parsing for feature flags

Potential approaches:
- Bytecode analysis (ASM library)
- Static analysis of annotation usage
- JSON/NBT data file scanning
- Integration with mod loader dev documentation



