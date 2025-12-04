export interface GraphAdapter {
  onNodeClick(callback: (nodeId: string) => void): void;
  getNodePosition(id: string): { x: number; y: number } | null;
  highlightNode(id: string): void;
  setNodeColor(id: string, color: string): void;
  destroy(): void;
  update(data: any): void;
  init(container: HTMLElement, data: any, options?: any): void;
}
