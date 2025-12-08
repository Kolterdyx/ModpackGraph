import { Injectable } from '@angular/core';
import { Environment, EnvironmentInfo } from '@wailsjs/runtime/runtime';
import { from, Observable, throwError, TimeoutError } from 'rxjs';
import { catchError, retry, timeout } from 'rxjs/operators';
import { AppConfigService } from '@core/config/app-config.service';
import * as WailsApp from '@wailsjs/go/app/App';

@Injectable({
  providedIn: 'root',
})
export class WailsRuntimeService {

  constructor(private configService: AppConfigService) {}

  getEnvironmentInfo(): Observable<EnvironmentInfo> {
    return this.handleObservable(Environment());
  }

  openDirectoryDialog(): Observable<string> {
    return this.handleObservable(WailsApp.OpenDirectoryDialog());
  }

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
  // Add more runtime utilities as needed
}

