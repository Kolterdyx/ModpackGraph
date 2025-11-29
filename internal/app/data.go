package app

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-graphviz"
)

var ignoredMods = map[string]struct{}{
	"minecraft": {}, "forge": {}, "neoforge": {}, "fabricloader": {},
	"fabric-loader": {}, "fabric": {}, "fabric-api": {}, "fabric_api": {},
	"fabric-resource-loader-v0": {}, "fabric-screen-api-v1": {},
	"fabric-networking-api-v1": {}, "fabric-lifecycle-events-v1": {},
	"fabric-renderer-api-v1": {}, "fabric-registry-sync-v0": {},
	"fabric-api-base": {}, "fabric-events-interaction-v0": {},
	"fabric-permissions-api-v0": {}, "fabric-command-api-v2": {},
	"fabric-kotlin": {}, "java": {},
}

func shouldIgnore(modid string) bool {
	if modid == "" {
		return true
	}
	_, ok := ignoredMods[strings.ToLower(modid)]
	return ok
}

type Mod struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type ModMetadata struct {
	Mod
	Name     string `json:"name"`
	Depends  []Dep  `json:"depends"`
	Embedded []Mod  `json:"embedded"`
	Path     string `json:"path"`
}

type Dep struct {
	ID            string `json:"id"`
	Required      bool   `json:"required"`
	Compatibility string `json:"compatibility,omitempty"`
}

// Extract metadata from bytes of a jar
func extractMetadataFromBytes(rawBytes []byte) (*ModMetadata, error) {
	r, err := zip.NewReader(bytes.NewReader(rawBytes), int64(len(rawBytes)))
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		log.Debugf("Checking file in jar: %s", f.Name)
		var meta *ModMetadata
		switch f.Name {
		// Fabric
		case "fabric.mod.json":
			meta, err = getFabricMetadata(r, f)
		// Forge modern
		case "META-INF/mods.toml":
			meta, err = getForgeMetadata(r, f)
		// Forge old mcmod.info
		case "mcmod.info":
			meta, err = getOldForgeMetadata(r, f)
		default:
			continue
		}

		if err != nil {
			log.WithError(err).Error("Error extracting metadata from %s", f.Name)
			continue
		}
		return meta, nil
	}

	return nil, nil
}

func getFabricMetadata(r *zip.Reader, f *zip.File) (*ModMetadata, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	var data map[string]any
	err = json.NewDecoder(rc).Decode(&data)
	if err != nil {
		return nil, err
	}
	err = rc.Close()
	if err != nil {
		return nil, err
	}

	modID, _ := data["id"].(string)
	name, _ := data["name"].(string)
	version, _ := data["version"].(string)
	if name == "" {
		name = modID
	}
	var depends []Dep
	for _, key := range []string{"depends", "recommends", "suggests"} {
		if val, ok := data[key].(map[string]any); ok {
			required := key == "depends"
			for k := range val {
				depends = append(depends, Dep{
					ID:            k,
					Compatibility: fmt.Sprintf("%v", val[k]),
					Required:      required,
				})
			}
		}
	}

	embedded, err := getEmbeddedMods(r)
	if err != nil {
		log.WithError(err).Error("Error getting embedded mods")
	}
	return &ModMetadata{
		Mod: Mod{
			ID:      modID,
			Version: version,
		},
		Name:     name,
		Depends:  depends,
		Embedded: embedded,
	}, nil
}

func getForgeMetadata(r *zip.Reader, f *zip.File) (*ModMetadata, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	data, _ := io.ReadAll(rc)
	err = rc.Close()
	if err != nil {
		return nil, err
	}

	var tomlData map[string]any
	if err := toml.Unmarshal(data, &tomlData); err != nil {
		return nil, err
	}

	modsArr, ok := tomlData["mods"].([]any)
	if !ok || len(modsArr) == 0 {
		return nil, fmt.Errorf("no mod entries found in mods.toml")
	}
	modEntry := modsArr[0].(map[string]any)
	modID, _ := modEntry["modId"].(string)
	version, _ := modEntry["version"].(string)
	name, _ := modEntry["displayName"].(string)
	if name == "" {
		name = modID
	}
	var depends []Dep
	if deps, ok := tomlData["dependencies"].(map[string]any); ok {
		if modDeps, ok := deps[modID].([]any); ok {
			for _, d := range modDeps {
				dm := d.(map[string]any)
				depID, _ := dm["modId"].(string)
				mandatory, _ := dm["mandatory"].(bool)
				compat, _ := dm["versionRange"].(string)
				compat = rangeToCompat(compat)
				depends = append(depends, Dep{
					ID:            depID,
					Compatibility: compat,
					Required:      mandatory,
				})
			}
		}
	}
	embedded, err := getEmbeddedMods(r)
	if err != nil {
		log.WithError(err).Error("Error getting embedded mods")
	}
	return &ModMetadata{
		Mod: Mod{
			ID:      modID,
			Version: version,
		},
		Name:     name,
		Depends:  depends,
		Embedded: embedded,
	}, nil
}

func getOldForgeMetadata(r *zip.Reader, f *zip.File) (*ModMetadata, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	var data []map[string]any

	if err = json.NewDecoder(rc).Decode(&data); err != nil {
		return nil, err
	}

	if err = rc.Close(); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("no mod entries found in mcmod.info")
	}
	entry := data[0]
	modID, _ := entry["modid"].(string)
	name, _ := entry["name"].(string)
	version, _ := entry["version"].(string)
	if name == "" {
		name = modID
	}
	var depends []Dep
	for _, dep := range entry["dependencies"].([]any) {
		if s, ok := dep.(string); ok {
			depends = append(depends, Dep{
				ID:       s,
				Required: true,
			})
		}
	}
	for _, dep := range entry["requiredMods"].([]any) {
		if s, ok := dep.(string); ok {
			depends = append(depends, Dep{
				ID:       s,
				Required: true,
			})
		}
	}
	embedded, err := getEmbeddedMods(r)
	if err != nil {
		log.WithError(err).Error("Error getting embedded mods")
	}
	return &ModMetadata{
		Mod: Mod{
			ID:      modID,
			Version: version,
		},
		Name:     name,
		Depends:  depends,
		Embedded: embedded,
	}, nil
}

func getEmbeddedMods(r *zip.Reader) ([]Mod, error) {

	var mods []Mod
	var paths []string
	for _, f := range r.File {
		paths = append(paths, f.Name)
	}

	// --- 1. META-INF/jars/ and jarjar/ embedded jars ---
	embeddedPrefixes := []string{
		"META-INF/jars/",
		"META-INF/jarjar/",
	}
	for _, prefix := range embeddedPrefixes {
		for _, name := range paths {
			if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ".jar") {
				f, err := r.Open(name)
				if err != nil {
					continue
				}
				data, err := io.ReadAll(f)
				_ = f.Close()
				if err != nil {
					continue
				}
				meta, err := extractMetadataFromBytes(data)
				if err != nil {
					continue
				}
				if meta != nil {
					mods = append(mods, meta.Mod)
				}
			}
		}
	}

	for _, f := range r.File {
		if f.Name == "fabric.mod.json" {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			var data map[string]any
			err = json.NewDecoder(rc).Decode(&data)
			if err != nil {
				_ = rc.Close()
				continue
			}
			_ = rc.Close()
			if jars, ok := data["jars"].([]any); ok {
				for _, entry := range jars {
					if emap, ok := entry.(map[string]any); ok {
						if id, ok := emap["id"].(string); ok {
							if version, ok := emap["version"].(string); ok {
								mods = append(mods, Mod{
									ID:      id,
									Version: version,
								})
							}
						}
					}
				}
			}
		}
	}
	return mods, nil
}

func rangeToCompat(compat string) string {
	if compat == "" {
		return "not specified"
	}
	var left, a, b, right string
	if strings.HasPrefix(compat, "[") || strings.HasPrefix(compat, "(") {
		left = string(compat[0])
		compat = compat[1:]
	}
	if strings.HasSuffix(compat, "]") || strings.HasSuffix(compat, ")") {
		right = string(compat[len(compat)-1])
		compat = compat[:len(compat)-1]
	}
	// now split by comma
	parts := strings.SplitN(compat, ",", 2)
	if len(parts) == 2 {
		a = strings.TrimSpace(parts[0])
		b = strings.TrimSpace(parts[1])
	} else {
		a = strings.TrimSpace(parts[0])
	}
	// reconstruct using >=, <=, >, <
	var result strings.Builder
	if a != "" {
		if left == "[" {
			result.WriteString(">=")
		} else if left == "(" {
			result.WriteString(">")
		}
		result.WriteString(a)
	}
	if b != "" {
		if result.Len() > 0 {
			result.WriteString(" and ")
		}
		if right == "]" {
			result.WriteString("<=")
		} else if right == ")" {
			result.WriteString("<")
		}
		result.WriteString(b)
	}
	return result.String()
}

// Extract metadata from a jar path
func extractModMetadata(jarPath string) *ModMetadata {
	data, err := os.ReadFile(jarPath)
	if err != nil {
		return nil
	}
	meta, err := extractMetadataFromBytes(data)
	if err != nil {
		log.Warnf("Failed to extract metadata from %s", jarPath)
	}
	if meta != nil {
		meta.Path = jarPath
	}
	return meta
}

// Scan folder
func scanModFolder(folder string) (map[string]*ModMetadata, error) {
	mods := make(map[string]*ModMetadata)
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".jar") {
			return nil
		}
		info := extractModMetadata(path)
		if info != nil && !shouldIgnore(info.ID) {
			var filtered []Dep
			for _, dep := range info.Depends {
				if !shouldIgnore(dep.ID) {
					filtered = append(filtered, dep)
				}
			}
			info.Depends = filtered
			mods[info.ID] = info
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return mods, nil
}

// generateDependencyGraphSVG generates the dependency graph SVG
func generateDependencyGraphSVG(ctx context.Context, mods map[string]*ModMetadata) ([]byte, error) {
	g, err := graphviz.New(ctx)
	if err != nil {
		return nil, err
	}
	graph, err := g.Graph()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = graph.Close()
		_ = g.Close()
	}()

	nodes := make(map[string]*graphviz.Node)
	for _, mod := range mods {
		node, err := graph.CreateNodeByName(fmt.Sprintf("%s\n%s", mod.Name, mod.Version))
		if err != nil {
			return nil, err
		}
		node.SetID(mod.ID)
		node.SetShape(graphviz.BoxShape)
		nodes[mod.ID] = node
	}
	for _, mod := range mods {
		for _, dep := range mod.Depends {
			depNode, exists := nodes[dep.ID]
			if !exists {

				if slices.ContainsFunc(mod.Embedded, func(m Mod) bool {
					return m.ID == dep.ID
				}) {
					// Dependency is embedded within the mod jar, nothing to plot
					continue
				} else {
					depNode, err = graph.CreateNodeByName(fmt.Sprintf("%s (missing)\nrequires: %s", dep.ID, dep.Compatibility))
					if err != nil {
						return nil, err
					}
					depNode.SetID(mod.ID)
					depNode.SetShape(graphviz.BoxShape)
					depNode.SetStyle(graphviz.FilledNodeStyle)
					if !dep.Required {
						depNode.SetFillColor("yellow")
					} else {
						depNode.SetFillColor("red")
					}
				}
			}
			_, err = graph.CreateEdgeByName(fmt.Sprintf("%s -> %s", mod.ID, dep.ID), nodes[mod.ID], depNode)

		}
	}
	g.SetLayout(graphviz.FDP)
	var buf bytes.Buffer
	if err = g.Render(ctx, graph, graphviz.SVG, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
