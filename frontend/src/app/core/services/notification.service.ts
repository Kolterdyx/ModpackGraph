import { Injectable } from '@angular/core';
import { MessageService } from 'primeng/api';

@Injectable({
  providedIn: 'root',
})
export class NotificationService {
  constructor(private messageService: MessageService) {}

  success(message: string, detail?: string): void {
    this.messageService.add({
      severity: 'success',
      summary: message,
      detail: detail,
      life: 3000,
    });
  }

  error(message: string, detail?: string): void {
    this.messageService.add({
      severity: 'error',
      summary: message,
      detail: detail,
      life: 5000,
    });
  }

  warn(message: string, detail?: string): void {
    this.messageService.add({
      severity: 'warn',
      summary: message,
      detail: detail,
      life: 4000,
    });
  }

  info(message: string, detail?: string): void {
    this.messageService.add({
      severity: 'info',
      summary: message,
      detail: detail,
      life: 3000,
    });
  }

  showProgress(message: string, progress: number): void {
    this.messageService.add({
      severity: 'info',
      summary: message,
      detail: `${progress}%`,
      life: 1000,
    });
  }
}

