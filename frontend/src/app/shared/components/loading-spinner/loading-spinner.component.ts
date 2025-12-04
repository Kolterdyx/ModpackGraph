import { Component } from '@angular/core';
import { ProgressSpinner } from 'primeng/progressspinner';

@Component({
  selector: 'app-loading-spinner',
  standalone: true,
  imports: [ProgressSpinner],
  template: `
    <div class="flex justify-center items-center p-4">
      <p-progressSpinner />
    </div>
  `,
})
export class LoadingSpinnerComponent {}

