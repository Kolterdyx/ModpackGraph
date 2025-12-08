import { Injectable } from '@angular/core';
import { ConfirmationService } from 'primeng/api';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ConfirmService {

  constructor(
    private readonly confirmService: ConfirmationService,
  ) {
  }

  public confirm(options: {
    message: string;
    header?: string;
    acceptLabel?: string;
    rejectLabel?: string;
  }): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      this.confirmService.confirm({
        message: options.message,
        header: options.header || $localize`Confirmation`,
        acceptLabel: options.acceptLabel || $localize`Yes`,
        acceptButtonProps: {
          // outlined: true,
          severity: 'danger',
        },
        rejectLabel: options.rejectLabel || $localize`No`,
        rejectButtonProps: {
          outlined: true,
          severity: 'secondary',
        },
        dismissableMask: true,
        accept: () => {
          observer.next(true);
          observer.complete();
        },
        reject: () => {
          observer.next(false);
          observer.complete();
        },
      });
    });
  }
}
