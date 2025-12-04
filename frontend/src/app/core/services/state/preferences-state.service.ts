import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { AppConfigService } from '@core/config/app-config.service';

interface PreferencesState {
  recentPaths: string[];
  language: string;
  theme: string;
}

@Injectable({
  providedIn: 'root',
})
export class PreferencesStateService {
  private readonly STORAGE_KEY = 'preferences-state';
  private readonly MAX_RECENT_PATHS: number;

  private recentPaths$ = new BehaviorSubject<string[]>([]);
  private language$ = new BehaviorSubject<string>('en');
  private theme$ = new BehaviorSubject<string>('light');

  constructor(private configService: AppConfigService) {
    this.MAX_RECENT_PATHS = configService.getMaxRecentPaths();
    this.loadPreferences();
  }

  private loadPreferences(): void {
    try {
      const stored = localStorage.getItem(this.STORAGE_KEY);
      if (stored) {
        const state: PreferencesState = JSON.parse(stored);
        this.recentPaths$.next(state.recentPaths || []);
        this.language$.next(state.language || 'en');
        this.theme$.next(state.theme || 'light');
      }
    } catch (error) {
      console.error('Failed to load preferences', error);
    }
  }

  private savePreferences(): void {
    try {
      const state: PreferencesState = {
        recentPaths: this.recentPaths$.value,
        language: this.language$.value,
        theme: this.theme$.value,
      };
      localStorage.setItem(this.STORAGE_KEY, JSON.stringify(state));
    } catch (error) {
      console.error('Failed to save preferences', error);
    }
  }

  addRecentPath(path: string): void {
    const current = this.recentPaths$.value;
    const filtered = current.filter((p) => p !== path);
    const updated = [path, ...filtered].slice(0, this.MAX_RECENT_PATHS);
    this.recentPaths$.next(updated);
    this.savePreferences();
  }

  getRecentPaths(): string[] {
    return this.recentPaths$.value;
  }

  getRecentPaths$(): Observable<string[]> {
    return this.recentPaths$.asObservable();
  }

  clearRecentPaths(): void {
    this.recentPaths$.next([]);
    this.savePreferences();
  }

  setLanguage(lang: string): void {
    this.language$.next(lang);
    this.savePreferences();
  }

  getLanguage$(): Observable<string> {
    return this.language$.asObservable();
  }

  setTheme(theme: string): void {
    this.theme$.next(theme);
    this.savePreferences();
  }

  getTheme$(): Observable<string> {
    return this.theme$.asObservable();
  }
}

