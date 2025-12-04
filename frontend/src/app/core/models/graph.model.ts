import { models } from '@wailsjs/go/models';

export class GraphModel {
  constructor(public graph: models.Graph) {}

  getFilteredGraph(filter: (node: models.Node) => boolean): models.Graph {
    const filteredNodes = this.graph.nodes.filter(filter);
    const nodeIds = new Set(filteredNodes.map((n) => n.id));

    const filteredEdges = this.graph.edges.filter(
      (edge) => nodeIds.has(edge.from) && nodeIds.has(edge.to)
    );

    return {
      nodes: filteredNodes,
      edges: filteredEdges,
    };
  }

  getNodeById(id: string): models.Node | undefined {
    return this.graph.nodes.find((n) => n.id === id);
  }

  getNodesByGroup(group: string): models.Node[] {
    return this.graph.nodes.filter((n) => n.group === group);
  }

  getConnectedNodes(nodeId: string): models.Node[] {
    const connectedIds = new Set<string>();

    this.graph.edges.forEach((edge) => {
      if (edge.from === nodeId) {
        connectedIds.add(edge.to);
      }
      if (edge.to === nodeId) {
        connectedIds.add(edge.from);
      }
    });

    return this.graph.nodes.filter((n) => connectedIds.has(n.id));
  }

  getNodeCount(): number {
    return this.graph.nodes.length;
  }

  getEdgeCount(): number {
    return this.graph.edges.length;
  }
}

