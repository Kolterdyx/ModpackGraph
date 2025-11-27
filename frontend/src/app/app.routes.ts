import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: "",
    loadComponent: () => import('@components/importer/importer').then(m => m.Importer)
  },
];
