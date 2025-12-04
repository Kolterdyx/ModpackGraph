import { Component, Input } from '@angular/core';
import { Card } from 'primeng/card';

@Component({
  selector: 'app-card',
  standalone: true,
  imports: [Card],
  template: `
    <p-card [header]="header" [subheader]="subheader">
      <ng-content />
    </p-card>
  `,
})
export class AppCardComponent {
  @Input() header?: string;
  @Input() subheader?: string;
}

