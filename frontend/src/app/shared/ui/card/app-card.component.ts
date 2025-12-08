import { Component, Input } from '@angular/core';
import { Card } from 'primeng/card';

@Component({
  selector: 'app-card',
  standalone: true,
  imports: [Card],
  templateUrl: './app-card.component.html',
})
export class AppCardComponent {
  @Input() header?: string;
  @Input() subheader?: string;
}

