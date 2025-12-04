import { models } from '@wailsjs/go/models';
import { getSeverityColor, getConflictTypeIcon } from '@core/utils/severity-color.util';

export class ConflictModel {
  constructor(public conflict: models.Conflict) {}

  getSeverityColor(): 'success' | 'info' | 'warn' | 'danger' {
    return getSeverityColor(this.conflict.severity);
  }

  getTypeIcon(): string {
    return getConflictTypeIcon(this.conflict.type);
  }

  getAffectedModNames(): string[] {
    return this.conflict.affected_mods || [];
  }

  getDescriptionFormatted(): string {
    return this.conflict.description || 'No description available';
  }

  isCritical(): boolean {
    return this.conflict.severity.toLowerCase() === 'critical';
  }

  isWarning(): boolean {
    return this.conflict.severity.toLowerCase() === 'warning';
  }
}

export class ConflictRuleModel {
  constructor(public rule: models.ConflictRule) {}

  getSeverityColor(): 'success' | 'info' | 'warn' | 'danger' {
    return getSeverityColor(this.rule.severity);
  }

  getTypeIcon(): string {
    return getConflictTypeIcon(this.rule.conflict_type);
  }

  getDescription(): string {
    return this.rule.description;
  }
}

