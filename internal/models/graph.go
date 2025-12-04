package models

// Graph represents a dependency graph for visualization
type Graph struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}

// Node represents a node in the dependency graph (a mod)
type Node struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Version string `json:"version"`
	Group   string `json:"group,omitempty"` // For visual grouping
}

// Edge represents an edge in the dependency graph (a dependency)
type Edge struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Required bool   `json:"required"`
	Label    string `json:"label,omitempty"` // Version constraint
}

// NewGraph creates a new empty graph
func NewGraph() *Graph {
	return &Graph{
		Nodes: make([]*Node, 0),
		Edges: make([]*Edge, 0),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(id, label, version string) {
	g.Nodes = append(g.Nodes, &Node{
		ID:      id,
		Label:   label,
		Version: version,
	})
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(from, to string, required bool, label string) {
	g.Edges = append(g.Edges, &Edge{
		From:     from,
		To:       to,
		Required: required,
		Label:    label,
	})
}
