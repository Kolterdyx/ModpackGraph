import { Injectable } from '@angular/core';
import { from, Observable, throwError, TimeoutError } from 'rxjs';
import { catchError, retry, timeout } from 'rxjs/operators';
import * as WailsApp from '@wailsjs/go/app/App';
import { models, services } from '@wailsjs/go/models';
import { AppConfigService } from '@core/config/app-config.service';

@Injectable({
  providedIn: 'root',
})
export class WailsAppService {
  constructor(private configService: AppConfigService) {}

  private handleObservable<T>(promise: Promise<T>): Observable<T> {
    return from(promise).pipe(
      timeout(this.configService.getWailsTimeout()),
      retry(this.configService.getWailsRetryAttempts()),
      catchError((error) => {
        if (error instanceof TimeoutError) {
          console.error('Wails API call timed out');
        }
        return throwError(() => error);
      })
    );
  }

  scanModpack(path: string): Observable<services.ScanResult> {
    return this.handleObservable(WailsApp.ScanModpack(path));
  }

  analyzeModpack(path: string): Observable<services.AnalysisReport> {
    return this.handleObservable(WailsApp.AnalyzeModpack(path));
  }

  getDependencyGraph(path: string): Observable<models.Graph> {
    return this.handleObservable(WailsApp.GetDependencyGraph(path));
  }

  getModpackStatus(path: string): Observable<services.ScanResult> {
    return this.handleObservable(WailsApp.GetModpackStatus(path));
  }

  getModMetadata(modId: string): Observable<models.ModMetadata> {
    return this.handleObservable(WailsApp.GetModMetadata(modId));
  }

  refreshModpack(path: string): Observable<services.ScanResult> {
    return this.handleObservable(WailsApp.RefreshModpack(path));
  }

  quickScan(path: string): Observable<services.ScanResult> {
    return this.handleObservable(WailsApp.QuickScan(path));
  }

  getConflictRules(): Observable<models.ConflictRule[]> {
    return this.handleObservable(WailsApp.GetConflictRules());
  }

  addConflictRule(
    modIdA: string,
    modIdB: string,
    conflictType: string,
    description: string,
    severity: string
  ): Observable<void> {
    return this.handleObservable(
      WailsApp.AddConflictRule(modIdA, modIdB, conflictType, description, severity)
    );
  }

  deleteConflictRule(id: number): Observable<void> {
    return this.handleObservable(WailsApp.DeleteConflictRule(id));
  }
}

