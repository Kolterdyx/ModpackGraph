import { GraphAdapter } from '@core/graph/graph-adapter.interface';

export class ForceGraph2DAdapter implements GraphAdapter {
  private graph: any;
  private clickCallback?: (nodeId: string) => void;

  async init(container: HTMLElement, data: any, options?: any): Promise<void> {
    // Dynamically import force-graph
    const ForceGraph = (await import('force-graph')).default;

    this.graph = new ForceGraph(container)
      .graphData(data)
      .nodeLabel('label')
      .nodeAutoColorBy('group')
      .linkDirectionalArrowLength(6)
      .linkDirectionalArrowRelPos(1)
      .onNodeClick((node: any) => {
        if (this.clickCallback) {
          this.clickCallback(node.id);
        }
      });

    if (options) {
      // Apply additional options
    }
  }

  update(data: any): void {
    if (this.graph) {
      this.graph.graphData(data);
    }
  }

  destroy(): void {
    if (this.graph) {
      this.graph._destructor();
      this.graph = null;
    }
  }

  setNodeColor(id: string, color: string): void {
    if (this.graph) {
      this.graph.nodeColor((node: any) => (node.id === id ? color : null));
    }
  }

  highlightNode(id: string): void {
    if (this.graph) {
      this.graph.centerAt(0, 0, 1000);
      // Implement node highlighting logic
    }
  }

  getNodePosition(id: string): { x: number; y: number } | null {
    // Implement position retrieval
    return null;
  }

  onNodeClick(callback: (nodeId: string) => void): void {
    this.clickCallback = callback;
  }
}

