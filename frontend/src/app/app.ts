import { Component, computed, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { Toast } from 'primeng/toast';
import { Button } from 'primeng/button';
import { WailsEventsService } from '@core/services/wails';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { ConfirmService } from '@services/confirm.service';
import { switchMap } from 'rxjs';
import { Tooltip } from 'primeng/tooltip';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    Toast,
    Button,
    ConfirmDialog,
    ConfirmDialog,
    Tooltip
  ],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App {

  protected isDarkTheme = signal<boolean>(true);
  protected themeChangeTooltip = computed(() => this.isDarkTheme() ? $localize`Switch to Light Theme` : $localize`Switch to Dark Theme`)

  constructor(
    eventService: WailsEventsService,
    confirmService: ConfirmService,
  ) {
    this.setTheme();
    eventService
      .on('on_before_close')
      .pipe(switchMap(() => confirmService.confirm({
        message: $localize`Are you sure you want to exit the application?`,
        header: $localize`Exit Confirmation`,
        acceptLabel: $localize`Exit`,
        rejectLabel: $localize`Cancel`,
      })))
      .subscribe((confirmed) => {
        if (confirmed !== undefined) {
          eventService.emit('app_close_response', !confirmed); // Send false to close the app, true to cancel
        }
      });
  }

  protected toggleTheme(): void {
    this.isDarkTheme.set(!this.isDarkTheme());
    this.setTheme();
  }

  private setTheme(): void {
    if (this.isDarkTheme()) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }
}
