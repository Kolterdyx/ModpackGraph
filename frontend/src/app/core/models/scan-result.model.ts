import { services } from '@wailsjs/go/models';

export class ScanResultModel {
  constructor(public result: services.ScanResult) {}

  hasChanges(): boolean {
    return (
      this.result.NewMods.length > 0 ||
      this.result.UpdatedMods.length > 0 ||
      this.result.RemovedMods.length > 0
    );
  }

  getChangesSummary(): string {
    const parts: string[] = [];

    if (this.result.NewMods.length > 0) {
      parts.push(`${this.result.NewMods.length} new`);
    }
    if (this.result.UpdatedMods.length > 0) {
      parts.push(`${this.result.UpdatedMods.length} updated`);
    }
    if (this.result.RemovedMods.length > 0) {
      parts.push(`${this.result.RemovedMods.length} removed`);
    }

    return parts.length > 0 ? parts.join(', ') : 'No changes';
  }

  getCacheStatistics(): { hits: number; misses: number; total: number; hitRate: number } {
    const hits = this.result.CacheHits;
    const misses = this.result.CacheMisses;
    const total = hits + misses;
    const hitRate = total > 0 ? (hits / total) * 100 : 0;

    return { hits, misses, total, hitRate };
  }

  getTotalMods(): number {
    return this.result.Mods.length;
  }
}

