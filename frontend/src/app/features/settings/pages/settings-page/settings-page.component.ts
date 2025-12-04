import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-settings-page',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-4">
      <h1 class="text-2xl font-bold mb-4">Settings</h1>

      <!-- TODO: Preferences form (language, theme) -->
      <!-- TODO: Conflict rules manager (add/delete rules) -->
      <!-- TODO: Cache management (clear cache button) -->
    </div>
  `,
})
export class SettingsPageComponent {}

