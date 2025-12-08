import { Component, computed, inject, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { Toast } from 'primeng/toast';
import { Button } from 'primeng/button';
import { WailsEventsService } from '@core/services/wails';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { ConfirmService } from '@core/services/confirm.service';
import { switchMap } from 'rxjs';
import { Tooltip } from 'primeng/tooltip';
import { ThemeService } from '@core/services/theme.service';
import { toSignal } from '@angular/core/rxjs-interop';

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
export class App implements OnInit {

  private readonly themeService = inject(ThemeService);
  protected readonly isDarkTheme = toSignal(this.themeService.getIsDarkTheme(), {initialValue: false});
  protected readonly themeChangeTooltip = computed(() => this.isDarkTheme() ? $localize`Switch to Light Theme` : $localize`Switch to Dark Theme`);

  constructor(
    private readonly eventService: WailsEventsService,
    private readonly confirmService: ConfirmService,
  ) {
  }

  ngOnInit(): void {
    this.themeService.getIsDarkTheme().subscribe((isDarkTheme) => {
      this.setTheme(isDarkTheme);
    });
    this.eventService
      .on('on_before_close')
      .pipe(switchMap(() => this.confirmService.confirm({
        message: $localize`Are you sure you want to exit the application?`,
        header: $localize`Exit Confirmation`,
        acceptLabel: $localize`Exit`,
        rejectLabel: $localize`Cancel`,
      })))
      .subscribe((confirmed) => {
        if (confirmed !== undefined) {
          this.eventService.emit('app_close_response', !confirmed); // Send false to close the app, true to cancel
        }
      });
  }

  protected toggleTheme(): void {
    this.themeService.toggleTheme()
  }

  private setTheme(isDarkTheme: boolean): void {
    if (isDarkTheme) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }
}
