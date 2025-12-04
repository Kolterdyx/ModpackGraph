import { GraphAdapter } from '@core/graph/graph-adapter.interface';

export class GraphListAdapter implements GraphAdapter {
  private data: any;
  private container?: HTMLElement;

  init(container: HTMLElement, data: any, options?: any): void {
    this.container = container;
    this.data = data;
    // List view doesn't need initialization like graph libraries
  }

  update(data: any): void {
    this.data = data;
  }

  destroy(): void {
    this.data = null;
    this.container = undefined;
  }

  setNodeColor(id: string, color: string): void {
    // Not applicable for list view
  }

  highlightNode(id: string): void {
    // Not applicable for list view
  }

  getNodePosition(id: string): { x: number; y: number } | null {
    return null;
  }

  onNodeClick(callback: (nodeId: string) => void): void {
    // Handled by the list component itself
  }

  getData(): any {
    return this.data;
  }
}

