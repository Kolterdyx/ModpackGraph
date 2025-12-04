import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { services } from '@wailsjs/go/models';
import { BaseStateService } from '@core/services/state/base-state.service';
import { AppConfigService } from '@core/config/app-config.service';

interface ModpackState {
  selectedPath: string | null;
  scanResult: services.ScanResult | null;
  isLoading: boolean;
  scanProgress: number;
}

@Injectable({
  providedIn: 'root',
})
export class ModpackStateService extends BaseStateService<ModpackState> {
  private selectedPath$ = new BehaviorSubject<string | null>(null);
  private scanResult$ = new BehaviorSubject<services.ScanResult | null>(null);
  private isLoading$ = new BehaviorSubject<boolean>(false);
  private scanProgress$ = new BehaviorSubject<number>(0);

  constructor(configService: AppConfigService) {
    super(
      'modpack-state',
      {
        selectedPath: null,
        scanResult: null,
        isLoading: false,
        scanProgress: 0,
      },
      configService
    );

    // Initialize from loaded state
    const state = this.getState();
    this.selectedPath$.next(state.selectedPath);
    this.scanResult$.next(state.scanResult);
    this.isLoading$.next(state.isLoading);
    this.scanProgress$.next(state.scanProgress);
  }

  selectModpack(path: string): void {
    this.selectedPath$.next(path);
    this.updateState({ selectedPath: path });
  }

  updateScanResult(result: services.ScanResult): void {
    this.scanResult$.next(result);
    this.updateState({ scanResult: result });
  }

  setLoading(loading: boolean): void {
    this.isLoading$.next(loading);
    this.updateState({ isLoading: loading });
  }

  setScanProgress(progress: number): void {
    this.scanProgress$.next(progress);
    this.updateState({ scanProgress: progress });
  }

  clearSelection(): void {
    this.selectedPath$.next(null);
    this.scanResult$.next(null);
    this.isLoading$.next(false);
    this.scanProgress$.next(0);
    this.state$.next(this.initialState);
  }

  getSelectedPath$(): Observable<string | null> {
    return this.selectedPath$.asObservable();
  }

  getScanResult$(): Observable<services.ScanResult | null> {
    return this.scanResult$.asObservable();
  }

  getIsLoading$(): Observable<boolean> {
    return this.isLoading$.asObservable();
  }

  getScanProgress$(): Observable<number> {
    return this.scanProgress$.asObservable();
  }

  private updateState(partial: Partial<ModpackState>): void {
    this.state$.next({ ...this.state$.value, ...partial });
  }
}

