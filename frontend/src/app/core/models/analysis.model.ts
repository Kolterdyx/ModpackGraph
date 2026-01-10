import { services } from '@wailsjs/go/models';

export class AnalysisReportModel {
  constructor(public report: services.AnalysisReport) {}

  getCriticalIssues(): number {
    return this.report.summary?.critical_conflicts || 0;
  }

  getWarnings(): number {
    return this.report.summary?.warning_conflicts || 0;
  }

  getSummaryText(): string {
    const summary = this.report.summary;
    if (!summary) return 'No analysis data available';

    const parts: string[] = [];
    parts.push(`${summary.total_mods} mods`);
    if (summary.new_mods > 0) parts.push(`${summary.new_mods} new`);
    if (summary.updated_mods > 0) parts.push(`${summary.updated_mods} updated`);
    if (summary.total_conflicts > 0) parts.push(`${summary.total_conflicts} conflicts`);

    return parts.join(', ');
  }

  hasErrors(): boolean {
    const summary = this.report.summary;
    return (summary?.critical_conflicts || 0) > 0 || (summary?.missing_dependencies || 0) > 0;
  }

  getCacheHitRateFormatted(): string {
    const rate = this.report.summary?.cache_hit_rate || 0;
    return `${rate.toFixed(1)}%`;
  }

  getTotalConflicts(): number {
    return this.report.summary?.total_conflicts || 0;
  }

  getMissingDependencies(): number {
    return this.report.summary?.missing_dependencies || 0;
  }

  getVersionConflicts(): number {
    return this.report.summary?.version_conflicts || 0;
  }
}

