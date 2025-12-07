export namespace models {
	
	export interface Conflict {
	    type: string;
	    severity: string;
	    description: string;
	    affected_mods: string[];
	    details?: Record<string, any>;
	}
	export interface ConflictRule {
	    id: number;
	    mod_id_a: string;
	    mod_id_b: string;
	    conflict_type: string;
	    description: string;
	    severity: string;
	}
	export interface VersionRange {
	    Raw: string;
	    Constraint?: semver.Constraints;
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
	    loader_type: string;
	    environment: string;
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
	    TotalMods: number;
	    NewMods: number;
	    UpdatedMods: number;
	    RemovedMods: number;
	    CacheHitRate: number;
	    MissingDependencies: number;
	    VersionConflicts: number;
	    CircularDependencies: number;
	    TotalConflicts: number;
	    CriticalConflicts: number;
	    WarningConflicts: number;
	    InfoConflicts: number;
	}
	export interface VersionConflict {
	    ModID: string;
	    ModName: string;
	    DependencyID: string;
	    DependencyName: string;
	    RequiredRange: string;
	    ActualVersion: string;
	}
	export interface MissingDependency {
	    ModID: string;
	    ModName: string;
	    DependencyID: string;
	    Required: boolean;
	    VersionRange: string;
	}
	export interface DependencyResult {
	    Graph?: models.Graph;
	    MissingDependencies: MissingDependency[];
	    VersionConflicts: VersionConflict[];
	    CircularDeps: string[][];
	}
	export interface ScanResult {
	    Modpack?: models.Modpack;
	    Mods: models.ModMetadata[];
	    NewMods: models.ModMetadata[];
	    UpdatedMods: models.ModMetadata[];
	    RemovedMods: models.ModMetadata[];
	    CacheHits: number;
	    CacheMisses: number;
	}
	export interface AnalysisReport {
	    ScanResult?: ScanResult;
	    DependencyResult?: DependencyResult;
	    Conflicts: models.Conflict[];
	    Summary?: AnalysisSummary;
	}
	
	
	
	

}

