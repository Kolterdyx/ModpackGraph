import { Component, Input, Output, EventEmitter } from '@angular/core';

import { Button } from 'primeng/button';
@Component({
  template: `<p-button
      (onClick)="onClick.emit($event)"
      [disabled]="disabled"
      [type]="type"
      [icon]="icon"
      [loading]="loading"
      [severity]="severity"
      [label]="label"
    />`,
  imports: [Button],
  standalone: true,
  selector: 'app-button',
})
export class AppButtonComponent {
  @Output() onClick = new EventEmitter<Event>();
  @Input() disabled = false;
  @Input() type: 'button' | 'submit' = 'button';
  @Input() icon?: string;
  @Input() loading = false;
  @Input() severity: 'primary' | 'secondary' | 'success' | 'info' | 'warn' | 'danger' | 'help' | 'contrast' = 'primary';
  @Input() label?: string;
}
