import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { services } from '@wailsjs/go/models';
import { BaseStateService } from '@core/services/state/base-state.service';
import { AppConfigService } from '@core/config/app-config.service';

interface AnalysisState {
  analysisReport: services.AnalysisReport | null;
  isAnalyzing: boolean;
  analysisProgress: number;
}

@Injectable({
  providedIn: 'root',
})
export class AnalysisStateService extends BaseStateService<AnalysisState> {
  private analysisReport$ = new BehaviorSubject<services.AnalysisReport | null>(null);
  private isAnalyzing$ = new BehaviorSubject<boolean>(false);
  private analysisProgress$ = new BehaviorSubject<number>(0);

  constructor(configService: AppConfigService) {
    super(
      'analysis-state',
      {
        analysisReport: null,
        isAnalyzing: false,
        analysisProgress: 0,
      },
      configService
    );

    // Initialize from loaded state
    const state = this.getState();
    this.analysisReport$.next(state.analysisReport);
    this.isAnalyzing$.next(state.isAnalyzing);
    this.analysisProgress$.next(state.analysisProgress);
  }

  setAnalysisReport(report: services.AnalysisReport): void {
    this.analysisReport$.next(report);
    this.updateState({ analysisReport: report });
  }

  startAnalysis(): void {
    this.isAnalyzing$.next(true);
    this.analysisProgress$.next(0);
    this.updateState({ isAnalyzing: true, analysisProgress: 0 });
  }

  updateAnalysisProgress(progress: number): void {
    this.analysisProgress$.next(progress);
    this.updateState({ analysisProgress: progress });
  }

  completeAnalysis(): void {
    this.isAnalyzing$.next(false);
    this.analysisProgress$.next(100);
    this.updateState({ isAnalyzing: false, analysisProgress: 100 });
  }

  clearAnalysis(): void {
    this.analysisReport$.next(null);
    this.isAnalyzing$.next(false);
    this.analysisProgress$.next(0);
    this.state$.next(this.initialState);
  }

  getAnalysisReport$(): Observable<services.AnalysisReport | null> {
    return this.analysisReport$.asObservable();
  }

  getIsAnalyzing$(): Observable<boolean> {
    return this.isAnalyzing$.asObservable();
  }

  getAnalysisProgress$(): Observable<number> {
    return this.analysisProgress$.asObservable();
  }

  private updateState(partial: Partial<AnalysisState>): void {
    this.state$.next({ ...this.state$.value, ...partial });
  }
}

