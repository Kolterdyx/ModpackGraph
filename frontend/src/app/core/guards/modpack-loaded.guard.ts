import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { map } from 'rxjs/operators';
import { ModpackStateService } from '@core/services/state/modpack-state.service';

export const modpackLoadedGuard: CanActivateFn = (route, state) => {
  const modpackStateService = inject(ModpackStateService);
  const router = inject(Router);

  return modpackStateService.getSelectedPath$().pipe(
    map((path) => {
      if (!path) {
        router.navigate(['/home']);
        return false;
      }
      return true;
    })
  );
};

