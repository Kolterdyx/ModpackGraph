package repository

const Schema = `
-- Mods table: cached mod metadata indexed by hash
CREATE TABLE IF NOT EXISTS mods (
    id TEXT PRIMARY KEY,
    hash TEXT UNIQUE NOT NULL,
    version TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    loader_type TEXT NOT NULL,
    environment TEXT NOT NULL,
    icon_data TEXT,
    metadata_json TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_mods_hash ON mods(hash);
CREATE INDEX IF NOT EXISTS idx_mods_loader_type ON mods(loader_type);

-- Mod authors (many-to-many)
CREATE TABLE IF NOT EXISTS mod_authors (
    mod_id TEXT NOT NULL,
    author TEXT NOT NULL,
    PRIMARY KEY (mod_id, author),
    FOREIGN KEY (mod_id) REFERENCES mods(id) ON DELETE CASCADE
);

-- Mod dependencies table
CREATE TABLE IF NOT EXISTS mod_dependencies (
    mod_id TEXT NOT NULL,
    dependency_id TEXT NOT NULL,
    required BOOLEAN NOT NULL DEFAULT 1,
    version_range TEXT,
    PRIMARY KEY (mod_id, dependency_id),
    FOREIGN KEY (mod_id) REFERENCES mods(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_mod_dependencies_mod_id ON mod_dependencies(mod_id);
CREATE INDEX IF NOT EXISTS idx_mod_dependencies_dependency_id ON mod_dependencies(dependency_id);

-- Modpacks table
CREATE TABLE IF NOT EXISTS modpacks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    path TEXT UNIQUE NOT NULL,
    last_scanned TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    mod_count INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_modpacks_path ON modpacks(path);

-- Modpack mods table (which mods are in which modpack)
CREATE TABLE IF NOT EXISTS modpack_mods (
    modpack_id INTEGER NOT NULL,
    mod_id TEXT NOT NULL,
    hash TEXT NOT NULL,
    file_path TEXT NOT NULL,
    PRIMARY KEY (modpack_id, mod_id),
    FOREIGN KEY (modpack_id) REFERENCES modpacks(id) ON DELETE CASCADE,
    FOREIGN KEY (mod_id) REFERENCES mods(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_modpack_mods_modpack_id ON modpack_mods(modpack_id);
CREATE INDEX IF NOT EXISTS idx_modpack_mods_hash ON modpack_mods(hash);

-- Mod features table (Phase 1 analysis)
CREATE TABLE IF NOT EXISTS mod_features (
    mod_id TEXT NOT NULL,
    feature_type TEXT NOT NULL,
    feature_data TEXT,
    extracted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (mod_id, feature_type),
    FOREIGN KEY (mod_id) REFERENCES mods(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_mod_features_feature_type ON mod_features(feature_type);

-- Conflict rules table (known incompatibilities)
CREATE TABLE IF NOT EXISTS conflict_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    mod_id_a TEXT NOT NULL,
    mod_id_b TEXT NOT NULL,
    conflict_type TEXT NOT NULL,
    description TEXT,
    severity TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_conflict_rules_mod_id_a ON conflict_rules(mod_id_a);
CREATE INDEX IF NOT EXISTS idx_conflict_rules_mod_id_b ON conflict_rules(mod_id_b);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_mods_updated_at
AFTER UPDATE ON mods
FOR EACH ROW
BEGIN
    UPDATE mods SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
`
