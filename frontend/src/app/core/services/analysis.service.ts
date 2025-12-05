import { Injectable } from '@angular/core';
import { Observable, tap, finalize, filter } from 'rxjs';
import { services } from '@wailsjs/go/models';
import { WailsAppService } from '@core/services/wails/wails-app.service';
import { WailsEventsService } from '@core/services/wails/wails-events.service';
import { AnalysisStateService } from '@core/services/state/analysis-state.service';
import { ModpackStateService } from '@core/services/state/modpack-state.service';

@Injectable({
  providedIn: 'root',
})
export class AnalysisService {
  constructor(
    private wailsAppService: WailsAppService,
    private wailsEventsService: WailsEventsService,
    private analysisStateService: AnalysisStateService,
    private modpackStateService: ModpackStateService
  ) {
    // Set up progress tracking as a long-lived subscription
    // In zoneless mode, state updates will trigger change detection via signals or async pipes
    this.wailsEventsService.onProgress().pipe(
      filter(event => event.operation === 'analyze')
    ).subscribe((event) => {
      this.analysisStateService.updateAnalysisProgress(event.progress);
    });
  }

  analyzeModpack(path: string): Observable<services.AnalysisReport> {
    this.analysisStateService.startAnalysis();
    this.modpackStateService.selectModpack(path);

    return this.wailsAppService.analyzeModpack(path).pipe(
      tap((report) => {
        this.analysisStateService.setAnalysisReport(report);
        this.analysisStateService.completeAnalysis();
      }),
      finalize(() => {
        this.analysisStateService.completeAnalysis();
      })
    );
  }

  quickScan(path: string): Observable<services.ScanResult> {
    return this.wailsAppService.quickScan(path);
  }

  getAnalysisProgress$(): Observable<number> {
    return this.analysisStateService.getAnalysisProgress$();
  }
}

