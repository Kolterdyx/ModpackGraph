package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/goccy/go-graphviz"
	log "github.com/sirupsen/logrus"
)

type Layout graphviz.Layout

var AllLayouts = []Layout{
	Layout(graphviz.CIRCO),
	Layout(graphviz.DOT),
	Layout(graphviz.FDP),
	Layout(graphviz.NEATO),
	Layout(graphviz.NOP),
	Layout(graphviz.NOP1),
	Layout(graphviz.NOP2),
	Layout(graphviz.OSAGE),
	Layout(graphviz.PATCHWORK),
	Layout(graphviz.SFDP),
	Layout(graphviz.TWOPI),
}

func (l Layout) TSName() string {
	return string(l)
}

func (l Layout) Graphviz() graphviz.Layout {
	return graphviz.Layout(l)
}

type GraphGenerationOptions struct {
	Path   string `json:"path,omitempty"`
	Layout Layout `json:"layout,omitempty"`
}

type Graph struct {
	Nodes  map[string]*Node `json:"nodes" ts_type:"Node[]"`
	Edges  map[string]*Edge `json:"links" ts_type:"Edge[]"`
	Layout Layout           `json:"layout,omitempty"`
}

func (g *Graph) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Nodes  []Node `json:"nodes" ts_type:"Node[]"`
		Edges  []Edge `json:"links" ts_type:"Edge[]"`
		Layout Layout `json:"layout,omitempty"`
	}
	nodes := make([]Node, 0, len(g.Nodes))
	for _, node := range g.Nodes {
		nodes = append(nodes, *node)
	}
	edges := make([]Edge, 0, len(g.Edges))
	for _, edge := range g.Edges {
		edges = append(edges, *edge)
	}
	return json.Marshal(&Alias{
		Nodes:  nodes,
		Edges:  edges,
		Layout: g.Layout,
	})
}

type Node struct {
	ID              string `json:"id,omitempty" ts_type:"string | number"`
	Label           string `json:"name,omitempty"`
	Icon            string `json:"icon,omitempty"`
	Present         bool   `json:"present,omitempty"`
	PresentVersion  string `json:"presentVersion,omitempty"`
	RequiredVersion Compat `json:"requiredVersion,omitempty" ts_type:"string"`
}

type Edge struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Label    string `json:"label,omitempty"`
	Required bool   `json:"required,omitempty"`
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string]*Edge),
	}
}

func (g *Graph) AddNode(node Node) *Node {
	g.Nodes[node.ID] = &node
	return g.Nodes[node.ID]
}

func (g *Graph) AddEdgeFromIDs(edge Edge) {
	if edge.Source == "" || edge.Target == "" {
		return
	}
	if edge.Source == edge.Target {
		return
	}
	// Prevent duplicate edges
	for _, e := range g.Edges {
		if e.Source == edge.Source && e.Target == edge.Target {
			return
		}
	}
	// Prevent edges between non-existent nodes
	sourceExists := false
	targetExists := false
	for _, node := range g.Nodes {
		if node.ID == edge.Source {
			sourceExists = true
		}
		if node.ID == edge.Target {
			targetExists = true
		}
	}
	if !sourceExists || !targetExists {
		return
	}
	// Add the edge
	g.Edges[fmt.Sprintf("%s->%s", edge.Source, edge.Target)] = &edge
}

func (g *Graph) Graphviz(ctx context.Context) (string, error) {
	gv, err := graphviz.New(ctx)
	if err != nil {
		return "", err
	}
	graph, err := gv.Graph()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = graph.Close()
		_ = gv.Close()
	}()
	nodeMap := make(map[string]*graphviz.Node)

	for _, node := range g.Nodes {
		gvNode, _ := graph.CreateNodeByName(node.ID)
		gvNode.SetLabel(node.Label)
		gvNode.SetStyle(graphviz.FilledNodeStyle)
		gvNode.SetShape(graphviz.BoxShape)
		gvNode.SetID(node.ID)
		nodeMap[node.ID] = gvNode
	}
	for _, edge := range g.Edges {
		sourceNode, sourceExists := nodeMap[edge.Source]
		targetNode, targetExists := nodeMap[edge.Target]
		if sourceExists && targetExists {
			gvEdge, err := graph.CreateEdgeByName(fmt.Sprintf("%s -> %s", edge.Source, edge.Target), sourceNode, targetNode)
			if err != nil {
				return "", err
			}
			gvEdge.SetDir(graphviz.ForwardDir)
		}
	}

	gv.SetLayout(g.Layout.Graphviz())
	var buf bytes.Buffer
	if err = gv.Render(ctx, graph, graphviz.SVG, &buf); err != nil {
		return "", err
	}
	log.Debug("Generated SVG data")
	return buf.String(), nil
}

func (g *Graph) GetNode(id string) (*Node, bool) {
	node, exists := g.Nodes[id]
	return node, exists
}

func (g *Graph) GetEdge(sourceID, targetID string) (*Edge, bool) {
	edge, exists := g.Edges[fmt.Sprintf("%s->%s", sourceID, targetID)]
	return edge, exists
}
