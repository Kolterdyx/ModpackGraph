import { Injectable } from '@angular/core';
import { environment } from './environment';

@Injectable({
  providedIn: 'root',
})
export class AppConfigService {
  private readonly config = environment;

  isFeatureEnabled(flag: keyof typeof environment.featureFlags): boolean {
    return this.config.featureFlags[flag];
  }

  getConfig() {
    return this.config;
  }

  isProduction(): boolean {
    return this.config.production;
  }

  getWailsTimeout(): number {
    return this.config.wails.timeout;
  }

  getWailsRetryAttempts(): number {
    return this.config.wails.retryAttempts;
  }

  getMaxRecentPaths(): number {
    return this.config.cache.maxRecentPaths;
  }
}

