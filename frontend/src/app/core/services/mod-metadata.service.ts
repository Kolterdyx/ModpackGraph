import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { tap } from 'rxjs/operators';
import { models } from '@wailsjs/go/models';
import { WailsAppService } from '@core/services/wails/wails-app.service';

@Injectable({
  providedIn: 'root',
})
export class ModMetadataService {
  private cache = new Map<string, models.ModMetadata>();

  constructor(private wailsAppService: WailsAppService) {}

  getModMetadata(modId: string): Observable<models.ModMetadata> {
    // Check cache first
    if (this.cache.has(modId)) {
      return of(this.cache.get(modId)!);
    }

    // Fetch from backend and cache
    return this.wailsAppService.getModMetadata(modId).pipe(
      tap((metadata) => {
        this.cache.set(modId, metadata);
      })
    );
  }

  clearCache(): void {
    this.cache.clear();
  }
}

