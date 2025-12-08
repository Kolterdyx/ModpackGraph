import { services } from '@wailsjs/go/models';

export class ScanResultModel {
  constructor(public result: services.ScanResult) {}

  hasChanges(): boolean {
    return (
      this.result.new_mods.length > 0 ||
      this.result.updated_mods.length > 0 ||
      this.result.removed_mods.length > 0
    );
  }

  getChangesSummary(): string {
    const parts: string[] = [];

    if (this.result.new_mods.length > 0) {
      parts.push(`${this.result.new_mods.length} new`);
    }
    if (this.result.updated_mods.length > 0) {
      parts.push(`${this.result.updated_mods.length} updated`);
    }
    if (this.result.removed_mods.length > 0) {
      parts.push(`${this.result.removed_mods.length} removed`);
    }

    return parts.length > 0 ? parts.join(', ') : 'No changes';
  }

  getCacheStatistics(): { hits: number; misses: number; total: number; hitRate: number } {
    const hits = this.result.cache_hits;
    const misses = this.result.cache_misses;
    const total = hits + misses;
    const hitRate = total > 0 ? (hits / total) * 100 : 0;

    return { hits, misses, total, hitRate };
  }

  getTotalMods(): number {
    return this.result.mods.length;
  }
}

