import { services } from '@wailsjs/go/models';

export class AnalysisReportModel {
  constructor(public report: services.AnalysisReport) {}

  getCriticalIssues(): number {
    return this.report.Summary?.CriticalConflicts || 0;
  }

  getWarnings(): number {
    return this.report.Summary?.WarningConflicts || 0;
  }

  getSummaryText(): string {
    const summary = this.report.Summary;
    if (!summary) return 'No analysis data available';

    const parts: string[] = [];
    parts.push(`${summary.TotalMods} mods`);
    if (summary.NewMods > 0) parts.push(`${summary.NewMods} new`);
    if (summary.UpdatedMods > 0) parts.push(`${summary.UpdatedMods} updated`);
    if (summary.TotalConflicts > 0) parts.push(`${summary.TotalConflicts} conflicts`);

    return parts.join(', ');
  }

  hasErrors(): boolean {
    const summary = this.report.Summary;
    return (summary?.CriticalConflicts || 0) > 0 || (summary?.MissingDependencies || 0) > 0;
  }

  getCacheHitRateFormatted(): string {
    const rate = this.report.Summary?.CacheHitRate || 0;
    return `${rate.toFixed(1)}%`;
  }

  getTotalConflicts(): number {
    return this.report.Summary?.TotalConflicts || 0;
  }

  getMissingDependencies(): number {
    return this.report.Summary?.MissingDependencies || 0;
  }

  getVersionConflicts(): number {
    return this.report.Summary?.VersionConflicts || 0;
  }
}

