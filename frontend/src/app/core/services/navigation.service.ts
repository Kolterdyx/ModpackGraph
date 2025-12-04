import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { encodePath, decodePath } from '@core/utils/path-encoder.util';

@Injectable({
  providedIn: 'root',
})
export class NavigationService {
  constructor(private router: Router) {}

  navigateToHome(): void {
    this.router.navigate(['/home']);
  }

  navigateToScan(path: string): void {
    this.router.navigate(['/scan'], {
      queryParams: { path: encodePath(path) },
    });
  }

  navigateToAnalysis(path: string): void {
    this.router.navigate(['/analysis'], {
      queryParams: { path: encodePath(path) },
    });
  }

  navigateToSettings(): void {
    this.router.navigate(['/settings']);
  }

  navigateToAbout(): void {
    this.router.navigate(['/about']);
  }

  decodePathParam(encodedPath: string): string {
    return decodePath(encodedPath);
  }
}

