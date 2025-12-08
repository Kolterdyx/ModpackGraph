export namespace models {
	
	export enum LoaderType {
	    fabric = "fabric",
	    forge_modern = "forge_modern",
	    forge_legacy = "forge_legacy",
	    neoforge = "neoforge",
	    quilt = "quilt",
	}
	export enum Environment {
	    client = "client",
	    server = "server",
	    both = "both",
	}
	export enum ConflictType {
	    missing_dependency = "missing_dependency",
	    version_conflict = "version_conflict",
	    known_incompatible = "known_incompatible",
	    feature_overlap = "feature_overlap",
	    environment_mismatch = "environment_mismatch",
	    circular_dependency = "circular_dependency",
	}
	export enum ConflictSeverity {
	    critical = "critical",
	    warning = "warning",
	    info = "info",
	}
	export enum FeatureType {
	    world_generation = "world_generation",
	    entities = "entities",
	    items = "items",
	    blocks = "blocks",
	    mechanics = "mechanics",
	    recipes = "recipes",
	}
	export interface Conflict {
	    type: ConflictType;
	    severity: ConflictSeverity;
	    description: string;
	    affected_mods: string[];
	    details?: Record<string, any>;
	}
	export interface ConflictRule {
	    id: number;
	    mod_id_a: string;
	    mod_id_b: string;
	    conflict_type: ConflictType;
	    description: string;
	    severity: ConflictSeverity;
	}
	export interface VersionRange {
	    raw: string;
	}
	export interface Dependency {
	    mod_id: string;
	    dependency_id: string;
	    required: boolean;
	    version_range?: VersionRange;
	}
	export interface Edge {
	    from: string;
	    to: string;
	    required: boolean;
	    label?: string;
	}
	export interface Node {
	    id: string;
	    label: string;
	    version: string;
	    group?: string;
	}
	export interface Graph {
	    nodes: Node[];
	    edges: Edge[];
	}
	export interface Mod {
	    id: string;
	    version: string;
	}
	export interface ModMetadata {
	    id: string;
	    hash: string;
	    version: string;
	    name: string;
	    description: string;
	    authors: string[];
	    loader_type: LoaderType;
	    environment: Environment;
	    icon_data?: string;
	    dependencies: Dependency[];
	    metadata_json?: string;
	    file_path?: string;
	    created_at: Date;
	    updated_at: Date;
	}
	export interface Modpack {
	    id: number;
	    name: string;
	    path: string;
	    // Go type: time
	    last_scanned: any;
	    mod_count: number;
	    mods?: Mod[];
	}
	

}

export namespace semver {
	
	export interface Constraints {
	    IncludePrerelease: boolean;
	}

}

export namespace services {
	
	export interface AnalysisSummary {
	    total_mods: number;
	    new_mods: number;
	    updated_mods: number;
	    removed_mods: number;
	    cache_hit_rate: number;
	    missing_dependencies: number;
	    version_conflicts: number;
	    circular_dependencies: number;
	    total_conflicts: number;
	    critical_conflicts: number;
	    warning_conflicts: number;
	    info_conflicts: number;
	}
	export interface VersionConflict {
	    mod_id: string;
	    mod_name: string;
	    dependency_id: string;
	    dependency_name: string;
	    required_range: string;
	    actual_version: string;
	}
	export interface MissingDependency {
	    mod_id: string;
	    mod_name: string;
	    dependency_id: string;
	    required: boolean;
	    version_range: string;
	}
	export interface DependencyResult {
	    graph?: models.Graph;
	    missing_dependencies: MissingDependency[];
	    version_conflicts: VersionConflict[];
	    circular_dependencies: string[][];
	}
	export interface ScanResult {
	    modpack?: models.Modpack;
	    mods: models.ModMetadata[];
	    new_mods: models.ModMetadata[];
	    updated_mods: models.ModMetadata[];
	    removed_mods: models.ModMetadata[];
	    cache_hits: number;
	    cache_misses: number;
	}
	export interface AnalysisReport {
	    scan_result?: ScanResult;
	    dependency_result?: DependencyResult;
	    conflicts: models.Conflict[];
	    summary?: AnalysisSummary;
	}
	
	
	
	

}

