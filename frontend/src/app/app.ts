import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { Toast } from 'primeng/toast';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, Toast],
  template: `
    <div class="app-shell flex h-screen">
      <!-- TODO: Add navigation sidebar -->
      <main class="flex-1 overflow-auto">
        <router-outlet />
      </main>
      <p-toast />
    </div>
  `,
  styleUrl: './app.scss',
})
export class App {}
