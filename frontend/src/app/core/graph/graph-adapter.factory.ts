import { Injectable } from '@angular/core';
import { GraphAdapter } from '@core/graph/graph-adapter.interface';

@Injectable({
  providedIn: 'root',
})
export class GraphAdapterFactory {
  private adapters = new Map<string, GraphAdapter>();

  async createAdapter(type: '2d' | '3d' | 'list'): Promise<GraphAdapter> {
    if (this.adapters.has(type)) {
      return this.adapters.get(type)!;
    }

    let adapter: GraphAdapter;

    switch (type) {
      case '2d': {
        const { ForceGraph2DAdapter } = await import('./force-graph-2d.adapter');
        adapter = new ForceGraph2DAdapter();
        break;
      }
      case '3d': {
        const { ForceGraph3DAdapter } = await import('./force-graph-3d.adapter');
        adapter = new ForceGraph3DAdapter();
        break;
      }
      case 'list': {
        const { GraphListAdapter } = await import('./graph-list.adapter');
        adapter = new GraphListAdapter();
        break;
      }
      default:
        throw new Error(`Unknown adapter type: ${type}`);
    }

    this.adapters.set(type, adapter);
    return adapter;
  }

  destroyAdapter(type: string): void {
    const adapter = this.adapters.get(type);
    if (adapter) {
      adapter.destroy();
      this.adapters.delete(type);
    }
  }

  destroyAll(): void {
    this.adapters.forEach((adapter) => adapter.destroy());
    this.adapters.clear();
  }
}

