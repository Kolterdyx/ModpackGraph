package app

import (
	"encoding/json"
	"fmt"
)

type GraphGenerationOptions struct {
	Path string `json:"path,omitempty"`
}

type Graph struct {
	Nodes map[string]*Node `json:"nodes" ts_type:"Node[]"`
	Edges map[string]*Edge `json:"links" ts_type:"Edge[]"`
}

func (g *Graph) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Nodes []Node `json:"nodes" ts_type:"Node[]"`
		Edges []Edge `json:"links" ts_type:"Edge[]"`
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
		Nodes: nodes,
		Edges: edges,
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

func (g *Graph) GetNode(id string) (*Node, bool) {
	node, exists := g.Nodes[id]
	return node, exists
}

func (g *Graph) GetEdge(sourceID, targetID string) (*Edge, bool) {
	edge, exists := g.Edges[fmt.Sprintf("%s->%s", sourceID, targetID)]
	return edge, exists
}
