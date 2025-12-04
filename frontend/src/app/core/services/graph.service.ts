import { Injectable } from '@angular/core';
import { models } from '@wailsjs/go/models';
import { DisplayOptions } from '@core/services/state/graph-state.service';

@Injectable({
  providedIn: 'root',
})
export class GraphService {
  buildGraphData(graph: models.Graph, format: '2d' | '3d'): any {
    if (!graph) {
      return null;
    }

    // Transform backend graph to library-specific format
    const nodes = graph.nodes.map((node) => ({
      id: node.id,
      label: node.label,
      version: node.version,
      group: node.group,
    }));

    const links = graph.edges.map((edge) => ({
      source: edge.from,
      target: edge.to,
      required: edge.required,
      label: edge.label,
    }));

    return { nodes, links };
  }

  applyFilters(graph: models.Graph, options: DisplayOptions): models.Graph {
    if (!graph) {
      return graph;
    }

    // Filter logic based on display options
    const filteredNodes = graph.nodes.filter((node) => {
      // Apply filter logic here
      return true;
    });

    const filteredEdges = graph.edges.filter((edge) => {
      if (!options.showRequired && edge.required) {
        return false;
      }
      if (!options.showOptional && !edge.required) {
        return false;
      }
      return true;
    });

    return {
      nodes: filteredNodes,
      edges: filteredEdges,
    };
  }
}

