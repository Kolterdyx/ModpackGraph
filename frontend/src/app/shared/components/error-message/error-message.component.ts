import { Component, Input } from '@angular/core';
import { Message } from 'primeng/message';

@Component({
  selector: 'app-error-message',
  standalone: true,
  imports: [Message],
  template: `
    <p-message
      severity="error"
      [text]="message"
    />
  `,
})
export class ErrorMessageComponent {
  @Input() message = 'An error occurred';
}

