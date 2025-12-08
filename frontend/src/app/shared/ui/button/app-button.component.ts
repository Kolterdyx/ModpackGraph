import { Component, Input, Output, EventEmitter } from '@angular/core';

import { Button } from 'primeng/button';
@Component({
  templateUrl: './app-button.component.html',
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
