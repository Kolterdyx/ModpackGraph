import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    redirectTo: '/home',
    pathMatch: 'full',
  },
  {
    path: 'home',
    loadChildren: () => import('./features/home/home.routes').then((m) => m.HOME_ROUTES),
  },
  {
    path: 'scan',
    loadChildren: () => import('./features/scan/scan.routes').then((m) => m.SCAN_ROUTES),
  },
  {
    path: 'analysis',
    loadChildren: () => import('./features/analysis/analysis.routes').then((m) => m.ANALYSIS_ROUTES),
  },
  {
    path: 'settings',
    loadChildren: () => import('./features/settings/settings.routes').then((m) => m.SETTINGS_ROUTES),
  },
  {
    path: 'about',
    loadChildren: () => import('./features/about/about.routes').then((m) => m.ABOUT_ROUTES),
  },
  {
    path: '**',
    redirectTo: '/home',
  },
];

