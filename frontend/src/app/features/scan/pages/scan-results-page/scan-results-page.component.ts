import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-scan-results-page',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-4">
      <h1 class="text-2xl font-bold mb-4">Scan Results</h1>
      <p class="mb-2">Path: {{ path$ | async }}</p>

      <!-- TODO: Display scan results -->
      <!-- TODO: Show mods list with AppDataTable -->
      <!-- TODO: Cache statistics cards -->
      <!-- TODO: New/updated/removed mods sections -->
    </div>
  `,
})
export class ScanResultsPageComponent {
  path$: Observable<string>;

  constructor(private route: ActivatedRoute) {
    this.path$ = this.route.queryParams.pipe(
      map((params) => decodeURIComponent(params['path'] || ''))
    );
  }
}

