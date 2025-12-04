import { Injectable } from '@angular/core';
import { Observable, tap, finalize } from 'rxjs';
import { services } from '@wailsjs/go/models';
import { WailsAppService } from '@core/services/wails/wails-app.service';
import { ModpackStateService } from '@core/services/state/modpack-state.service';
import { PreferencesStateService } from '@core/services/state/preferences-state.service';

@Injectable({
  providedIn: 'root',
})
export class ModpackService {
  constructor(
    private wailsAppService: WailsAppService,
    private modpackStateService: ModpackStateService,
    private preferencesStateService: PreferencesStateService
  ) {}

  scanModpack(path: string): Observable<services.ScanResult> {
    this.modpackStateService.selectModpack(path);
    this.modpackStateService.setLoading(true);
    this.modpackStateService.setScanProgress(0);

    return this.wailsAppService.scanModpack(path).pipe(
      tap((result) => {
        this.modpackStateService.updateScanResult(result);
        this.preferencesStateService.addRecentPath(path);
        this.modpackStateService.setScanProgress(100);
      }),
      finalize(() => {
        this.modpackStateService.setLoading(false);
      })
    );
  }

  refreshModpack(path: string): Observable<services.ScanResult> {
    this.modpackStateService.setLoading(true);
    this.modpackStateService.setScanProgress(0);

    return this.wailsAppService.refreshModpack(path).pipe(
      tap((result) => {
        this.modpackStateService.updateScanResult(result);
        this.modpackStateService.setScanProgress(100);
      }),
      finalize(() => {
        this.modpackStateService.setLoading(false);
      })
    );
  }

  getModpackStatus(path: string): Observable<services.ScanResult> {
    return this.wailsAppService.getModpackStatus(path);
  }

  quickScan(path: string): Observable<services.ScanResult> {
    return this.wailsAppService.quickScan(path).pipe(
      tap((result) => {
        this.modpackStateService.updateScanResult(result);
      })
    );
  }
}

