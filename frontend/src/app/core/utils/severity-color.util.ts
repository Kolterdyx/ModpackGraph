/**
 * Map conflict severity to PrimeNG severity
 */
export function getSeverityColor(severity: string): 'success' | 'info' | 'warn' | 'danger' {
  switch (severity.toLowerCase()) {
    case 'critical':
      return 'danger';
    case 'warning':
      return 'warn';
    case 'info':
      return 'info';
    default:
      return 'info';
  }
}

/**
 * Get icon for conflict type
 */
export function getConflictTypeIcon(type: string): string {
  switch (type.toLowerCase()) {
    case 'missing_dependency':
      return 'pi pi-exclamation-circle';
    case 'version_conflict':
      return 'pi pi-times-circle';
    case 'known_incompatible':
      return 'pi pi-ban';
    case 'circular_dependency':
      return 'pi pi-replay';
    default:
      return 'pi pi-info-circle';
  }
}

