package app

import (
	"bytes"
	"context"
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
	Nodes  []Node `json:"nodes"`
	Edges  []Edge `json:"links"`
	Layout Layout `json:"layout,omitempty"`
}

type Node struct {
	ID    string `json:"id"`
	Color string `json:"color,omitempty"`
	Label string `json:"name,omitempty"`
	Value int    `json:"val,omitempty"`
	Icon  string `json:"icon,omitempty"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make([]Node, 0),
		Edges: make([]Edge, 0),
	}
}

func (g *Graph) AddNode(node Node) Node {
	g.Nodes = append(g.Nodes, node)
	return node
}

func (g *Graph) AddEdgeFromIDs(sourceID, targetID string) {
	if sourceID == targetID {
		return
	}
	// Prevent duplicate edges
	for _, edge := range g.Edges {
		if (edge.Source == sourceID && edge.Target == targetID) ||
			(edge.Source == targetID && edge.Target == sourceID) {
			return
		}
	}
	// Prevent edges between non-existent nodes
	sourceExists := false
	targetExists := false
	for _, node := range g.Nodes {
		if node.ID == sourceID {
			sourceExists = true
		}
		if node.ID == targetID {
			targetExists = true
		}
	}
	if !sourceExists || !targetExists {
		return
	}
	// Add the edge
	g.Edges = append(g.Edges, Edge{
		Source: sourceID,
		Target: targetID,
	})
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
		gvNode.SetFillColor(node.Color)
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
