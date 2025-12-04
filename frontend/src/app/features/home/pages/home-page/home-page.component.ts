import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home-page',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-4">
      <h1 class="text-3xl font-bold mb-4">ModpackGraph</h1>
      <p class="mb-4">Select a modpack directory to begin analysis.</p>

      <!-- TODO: Add directory selection UI -->
      <!-- TODO: Add recent paths list -->
      <!-- TODO: Add action buttons (Scan, Analyze) -->
    </div>
  `,
})
export class HomePageComponent {}

