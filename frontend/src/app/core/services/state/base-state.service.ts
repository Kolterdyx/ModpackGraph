import { BehaviorSubject, Observable } from 'rxjs';
import { AppConfigService } from '@core/config/app-config.service';

export abstract class BaseStateService<T> {
  protected state$: BehaviorSubject<T>;

  constructor(
    protected storageKey: string,
    protected initialState: T,
    protected configService: AppConfigService
  ) {
    const loadedState = this.loadFromStorage();
    this.state$ = new BehaviorSubject<T>(loadedState ?? initialState);

    // Subscribe to state changes and persist if enabled
    this.state$.subscribe((state) => {
      if (this.shouldPersist()) {
        this.saveToStorage(state);
      }
    });
  }

  protected shouldPersist(): boolean {
    return this.configService.isFeatureEnabled('persistState');
  }

  protected saveToStorage(state: T): void {
    try {
      localStorage.setItem(this.storageKey, JSON.stringify(state));
    } catch (error) {
      console.error(`Failed to save state to localStorage: ${this.storageKey}`, error);
    }
  }

  protected loadFromStorage(): T | null {
    try {
      const stored = localStorage.getItem(this.storageKey);
      if (stored) {
        return JSON.parse(stored) as T;
      }
    } catch (error) {
      console.error(`Failed to load state from localStorage: ${this.storageKey}`, error);
    }
    return null;
  }

  clearStorage(): void {
    try {
      localStorage.removeItem(this.storageKey);
    } catch (error) {
      console.error(`Failed to clear storage: ${this.storageKey}`, error);
    }
  }

  getState(): T {
    return this.state$.value;
  }

  getState$(): Observable<T> {
    return this.state$.asObservable();
  }
}

