import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { map } from 'rxjs/operators';
import { AnalysisStateService } from '@core/services/state/analysis-state.service';

export const analysisLoadedGuard: CanActivateFn = (route, state) => {
  const analysisStateService = inject(AnalysisStateService);
  const router = inject(Router);

  return analysisStateService.getAnalysisReport$().pipe(
    map((report) => {
      if (!report) {
        router.navigate(['/home']);
        return false;
      }
      return true;
    })
  );
};

