import { models } from '@wailsjs/go/models';
import { formatVersion } from '@core/utils/format-version.util';

export class ModModel {
  constructor(public metadata: models.ModMetadata) {}

  getDisplayName(): string {
    return this.metadata.name || this.metadata.id;
  }

  getStatusBadge(): 'success' | 'info' | 'warn' | 'danger' {
    // Determine status based on metadata
    return 'info';
  }

  hasIcon(): boolean {
    return !!this.metadata.icon_data;
  }

  formatAuthors(): string {
    return this.metadata.authors?.join(', ') || 'Unknown';
  }

  isClientOnly(): boolean {
    return this.metadata.environment === 'client';
  }

  isServerOnly(): boolean {
    return this.metadata.environment === 'server';
  }

  getLoaderName(): string {
    switch (this.metadata.loader_type) {
      case 'fabric':
        return 'Fabric';
      case 'forge_modern':
        return 'Forge (Modern)';
      case 'forge_legacy':
        return 'Forge (Legacy)';
      case 'neoforge':
        return 'NeoForge';
      default:
        return 'Unknown';
    }
  }

  getFormattedVersion(): string {
    return formatVersion(this.metadata.version);
  }
}

