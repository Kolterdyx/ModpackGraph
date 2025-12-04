import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { models } from '@wailsjs/go/models';
import { WailsAppService } from '@core/services/wails/wails-app.service';

@Injectable({
  providedIn: 'root',
})
export class ConflictRulesService {
  constructor(private wailsAppService: WailsAppService) {}

  getRules(): Observable<models.ConflictRule[]> {
    return this.wailsAppService.getConflictRules();
  }

  addRule(
    modIdA: string,
    modIdB: string,
    conflictType: string,
    description: string,
    severity: string
  ): Observable<void> {
    return this.wailsAppService.addConflictRule(
      modIdA,
      modIdB,
      conflictType,
      description,
      severity
    );
  }

  deleteRule(id: number): Observable<void> {
    return this.wailsAppService.deleteConflictRule(id);
  }
}

