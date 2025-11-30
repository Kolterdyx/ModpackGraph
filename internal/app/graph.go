package app

import (
	"context"
	"fmt"
	"github.com/goccy/go-graphviz"
)

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"links"`
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

func (g *Graph) Graphviz(ctx context.Context) (*graphviz.Graphviz, *graphviz.Graph, error) {
	gv, err := graphviz.New(ctx)
	if err != nil {
		return nil, nil, err
	}
	graph, _ := gv.Graph()
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
				return nil, nil, err
			}
			gvEdge.SetDir(graphviz.ForwardDir)
		}
	}
	return gv, graph, nil
}
