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
	"strings"

	"github.com/goccy/go-graphviz"
)

var ignoredMods = map[string]struct{}{
	"minecraft":                    {},
	"forge":                        {},
	"neoforge":                     {},
	"fabricloader":                 {},
	"fabric-loader":                {},
	"fabric":                       {},
	"fabric-api":                   {},
	"fabric_api":                   {},
	"fabric-resource-loader-v0":    {},
	"fabric-screen-api-v1":         {},
	"fabric-networking-api-v1":     {},
	"fabric-lifecycle-events-v1":   {},
	"fabric-renderer-api-v1":       {},
	"fabric-registry-sync-v0":      {},
	"fabric-api-base":              {},
	"fabric-events-interaction-v0": {},
	"fabric-permissions-api-v0":    {},
	"fabric-command-api-v2":        {},
	"fabric-kotlin":                {},
	"java":                         {},
}

func shouldIgnore(modid string, ignored map[string]struct{}) bool {
	if modid == "" {
		return true
	}
	_, ok := ignoredMods[strings.ToLower(modid)]
	if ok {
		return true
	}
	_, ok = ignored[modid]
	return ok
}

type Mod struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type ModMetadata struct {
	Mod
	Name    string `json:"name"`
	Depends []Dep  `json:"depends"`
	Path    string `json:"path"`
}

type Dep struct {
	ID            string `json:"id"`
	Required      bool   `json:"required"`
	Compatibility string `json:"compatibility,omitempty"`
}

func getFabricMetadata(f *zip.File) (ModMetadata, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered in getFabricMetadata: %v", r)
		}
	}()
	rc, err := f.Open()
	if err != nil {
		return ModMetadata{}, err
	}
	var data map[string]any
	err = json.NewDecoder(rc).Decode(&data)
	if err != nil {
		return ModMetadata{}, err
	}
	err = rc.Close()
	if err != nil {
		return ModMetadata{}, err
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

	return ModMetadata{
		Mod: Mod{
			ID:      modID,
			Version: version,
		},
		Name:    name,
		Depends: depends,
	}, nil
}

func getForgeMetadata(r *zip.Reader, f *zip.File) (ModMetadata, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered in getForgeMetadata: %v", r)
		}
	}()
	rc, err := f.Open()
	if err != nil {
		return ModMetadata{}, err
	}
	data, _ := io.ReadAll(rc)
	err = rc.Close()
	if err != nil {
		return ModMetadata{}, err
	}

	var tomlData map[string]any
	if err := toml.Unmarshal(data, &tomlData); err != nil {
		return ModMetadata{}, err
	}
	modsAny, ok := tomlData["mods"]
	if !ok {
		return ModMetadata{}, fmt.Errorf("no mods section found in mods.toml")
	}
	modsArr, ok := modsAny.([]any)
	if !ok || len(modsArr) == 0 {
		return ModMetadata{}, fmt.Errorf("no mod entries found in mods.toml")
	}
	modEntry, ok := modsArr[0].(map[string]any)
	if !ok {
		return ModMetadata{}, fmt.Errorf("invalid mod entry in mods.toml")
	}
	modID, ok := modEntry["modId"].(string)
	if !ok {
		return ModMetadata{}, fmt.Errorf("modId not found in mods.toml")
	}
	version, ok := modEntry["version"].(string)
	if !ok || version == "" {
		version = "<not specified>"
	}
	if version == "${file.jarVersion}" {
		// extract from MANIFEST.MF
		for _, mfFile := range r.File {
			if mfFile.Name == "META-INF/MANIFEST.MF" {
				rc, err := mfFile.Open()
				if err != nil {
					return ModMetadata{}, err
				}
				manifestData, _ := io.ReadAll(rc)
				err = rc.Close()
				if err != nil {
					return ModMetadata{}, err
				}
				lines := strings.Split(string(manifestData), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "Implementation-Version:") {
						version = strings.TrimSpace(strings.TrimPrefix(line, "Implementation-Version:"))
						break
					}
				}
				break
			}
		}
	}
	name, ok := modEntry["displayName"].(string)
	if !ok || name == "" {
		name = modID
	}
	var depends []Dep
	if deps, ok := tomlData["dependencies"].(map[string]any); ok {
		modDepsAny, ok := deps[modID]
		if ok {
			if modDeps, ok := modDepsAny.([]any); ok {
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
	}

	return ModMetadata{
		Mod: Mod{
			ID:      modID,
			Version: version,
		},
		Name:    name,
		Depends: depends,
	}, nil
}

func getOldForgeMetadata(r *zip.Reader, f *zip.File) (ModMetadata, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered in getOldForgeMetadata: %v", r)
		}
	}()
	rc, err := f.Open()
	if err != nil {
		return ModMetadata{}, err
	}
	var data []map[string]any

	if err = json.NewDecoder(rc).Decode(&data); err != nil {
		return ModMetadata{}, err
	}

	if err = rc.Close(); err != nil {
		return ModMetadata{}, err
	}
	if len(data) == 0 {
		return ModMetadata{}, fmt.Errorf("no mod entries found in mcmod.info")
	}
	entry := data[0]
	modID, ok := entry["modid"].(string)
	if !ok {
		return ModMetadata{}, fmt.Errorf("modid not found in mcmod.info")
	}
	name, ok := entry["name"].(string)
	if !ok {
		name = modID
	}
	version, ok := entry["version"].(string)
	if !ok {
		version = "<not specified>"
	}
	if version == "${version}" {
		// extract from MANIFEST.MF
		for _, mfFile := range r.File {
			if mfFile.Name == "META-INF/MANIFEST.MF" {
				rc, err := mfFile.Open()
				if err != nil {
					return ModMetadata{}, err
				}
				manifestData, _ := io.ReadAll(rc)
				err = rc.Close()
				if err != nil {
					return ModMetadata{}, err
				}
				lines := strings.Split(string(manifestData), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "Implementation-Version:") {
						version = strings.TrimSpace(strings.TrimPrefix(line, "Implementation-Version:"))
						break
					}
				}
				break
			}
		}
	}
	if name == "" {
		name = modID
	}
	var depends []Dep

	if deps, ok := entry["dependencies"]; ok {
		if depsArr, ok := deps.([]any); ok {
			for _, dep := range depsArr {
				if s, ok := dep.(string); ok {
					depends = append(depends, Dep{
						ID:       s,
						Required: true,
					})
				}
			}
		}
	}

	if deps, ok := entry["requiredMods"]; ok {
		if depsArr, ok := deps.([]any); ok {
			for _, dep := range depsArr {
				if s, ok := dep.(string); ok {
					depends = append(depends, Dep{
						ID:       s,
						Required: true,
					})
				}
			}
		}
	}
	return ModMetadata{
		Mod: Mod{
			ID:      modID,
			Version: version,
		},
		Name:    name,
		Depends: depends,
	}, nil
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
func extractModMetadata(path string, r *zip.Reader) (ModMetadata, error) {
	log.Debugf("Extracting mod metadata from %s", path)
	var err error
	var meta ModMetadata
	for _, f := range r.File {
		err = nil
		switch f.Name {
		// Fabric
		case "fabric.mod.json":
			meta, err = getFabricMetadata(f)
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
		} else {
			break
		}
	}
	meta.Path = path
	return meta, err
}

func getModJarsFromBytes(name string, jarBytes []byte) map[string]*zip.Reader {
	r, err := zip.NewReader(bytes.NewReader(jarBytes), int64(len(jarBytes)))
	if err != nil {
		return nil
	}
	// Collect embedded jars
	var jars map[string]*zip.Reader
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".jar") {
			rc, err := f.Open()
			if err != nil {
				log.WithError(err).Error("Error opening jar file")
				continue
			}
			data, err := io.ReadAll(rc)
			_ = rc.Close()
			if err != nil {
				log.WithError(err).Error("Error reading jar file")
				continue
			}
			embeddedJars := getModJarsFromBytes(f.Name, data)
			for k, v := range embeddedJars {
				if jars == nil {
					jars = make(map[string]*zip.Reader)
				}
				jars[k] = v
			}
		}
	}
	if jars == nil {
		jars = make(map[string]*zip.Reader)
	}
	jars[name] = r
	return jars
}

func getModJars(jarPath string) map[string]*zip.Reader {
	data, err := os.ReadFile(jarPath)
	if err != nil {
		return nil
	}
	return getModJarsFromBytes(jarPath, data)
}

// Scan folder
func scanModFolder(folder string) (*Graph, error) {
	var jars map[string]*zip.Reader
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".jar") {
			return nil
		}
		modJars := getModJars(path)
		for p, r := range modJars {
			if jars == nil {
				jars = make(map[string]*zip.Reader)
			}
			jars[p] = r
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %d jars", len(jars))
	ignored := make(map[string]struct{})
	mods := make(map[string]ModMetadata)
	for path, r := range jars {
		info, err := extractModMetadata(path, r)
		if err != nil {
			log.WithError(err).WithField("path", path).Error("Error extracting mod metadata")
			continue
		}
		if strings.HasPrefix(info.Path, "META-INF") {
			log.Warnf("Ignoring embedded mod %+v", info)
			ignored[info.ID] = struct{}{}
			continue
		}
		if shouldIgnore(info.ID, ignored) {
			continue
		}
		var filtered []Dep
		for _, dep := range info.Depends {
			if shouldIgnore(dep.ID, ignored) {
				continue
			}
			filtered = append(filtered, dep)
		}
		info.Depends = filtered
		mods[info.ID] = info
	}
	log.Debugf("Extracted metadata for %d mods", len(mods))
	// Filter embedded mods from dependencies
	for k, mod := range mods {
		var filtered []Dep
		for _, dep := range mod.Depends {

			if depMod, exists := mods[dep.ID]; (exists && strings.HasPrefix(depMod.Path, "META-INF")) || shouldIgnore(dep.ID, ignored) {
				log.Warnf("Filtering embedded mod dependency %+v", dep)
				continue
			}
			filtered = append(filtered, dep)
		}
		mod.Depends = filtered
		mods[k] = mod
	}
	log.Debugf("Found %d mods", len(mods))
	return generateDependencyGraph(mods)
}

func generateDependencyGraph(mods map[string]ModMetadata) (*Graph, error) {
	graph := NewGraph()
	embeddings := make(map[string]struct{})
	nodes := make(map[string]Node)
	for _, mod := range mods {
		node := graph.AddNode(Node{
			ID:    mod.ID,
			Label: fmt.Sprintf("%s\n%s", mod.Name, mod.Version),
		})
		if strings.HasPrefix("META-INF", mod.Path) {
			log.Warnf("Marking %s as embedded", mod.ID)
			embeddings[mod.ID] = struct{}{}
		}
		nodes[mod.ID] = node
	}
	for _, mod := range mods {
		for _, dep := range mod.Depends {
			depNode, exists := nodes[dep.ID]
			if !exists {
				depNode = Node{
					ID:    dep.ID,
					Label: fmt.Sprintf("%s (missing)\nrequires: %s", dep.ID, dep.Compatibility),
				}
				if !dep.Required {
					depNode.Color = "yellow"
				} else {
					depNode.Color = "red"
				}
				nodes[dep.ID] = graph.AddNode(depNode)
			}
			graph.AddEdgeFromIDs(mod.ID, dep.ID)
		}
	}
	return graph, nil
}

func generateDependencyGraphSVG(ctx context.Context, modGraph *Graph, options GraphOptions) ([]byte, error) {

	g, graph, err := modGraph.Graphviz(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = graph.Close()
		_ = g.Close()
	}()

	g.SetLayout(options.Layout.Graphviz())
	var buf bytes.Buffer
	if err = g.Render(ctx, graph, graphviz.SVG, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
