import { Component, Input } from '@angular/core';
import { Message } from 'primeng/message';

@Component({
  selector: 'app-error-message',
  standalone: true,
  imports: [Message],
  templateUrl: './error-message.component.html',
})
export class ErrorMessageComponent {
  @Input() message = 'An error occurred';
}

