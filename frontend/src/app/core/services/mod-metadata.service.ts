import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
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
      return new Observable((observer) => {
        observer.next(this.cache.get(modId)!);
        observer.complete();
      });
    }

    // Fetch from backend and cache
    return new Observable((observer) => {
      this.wailsAppService.getModMetadata(modId).subscribe({
        next: (metadata) => {
          this.cache.set(modId, metadata);
          observer.next(metadata);
          observer.complete();
        },
        error: (error) => observer.error(error),
      });
    });
  }

  clearCache(): void {
    this.cache.clear();
  }
}

