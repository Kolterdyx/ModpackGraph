import { models } from '@wailsjs/go/models';
import { formatDate, formatRelativeTime } from '@core/utils/format-date.util';

export class ModpackModel {
  constructor(public modpack: models.Modpack) {
  }

  getFormattedScanTime(): string {
    return formatDate(this.modpack.last_scanned);
  }

  getRelativeScanTime(): string {
    return formatRelativeTime(this.modpack.last_scanned);
  }

  isEmpty(): boolean {
    return this.modpack.mod_count === 0;
  }

  getModCount(): number {
    return this.modpack.mod_count;
  }

  hasBeenScanned(): boolean {
    return !!this.modpack.last_scanned;
  }

  getName(): string {
    return this.modpack.name;
  }

  getPath(): string {
    return this.modpack.path;
  }
}

