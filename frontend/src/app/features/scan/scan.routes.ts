import { Routes } from '@angular/router';

export const SCAN_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/scan-results-page/scan-results-page.component').then(
        (m) => m.ScanResultsPageComponent
      ),
  },
];

