import { Component } from '@angular/core';
import { ProgressSpinner } from 'primeng/progressspinner';

@Component({
  selector: 'app-loading-spinner',
  standalone: true,
  imports: [ProgressSpinner],
  templateUrl: './loading-spinner.component.html',
})
export class LoadingSpinnerComponent {}

