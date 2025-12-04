import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';

@Component({
  selector: 'app-scan-results-page',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-4">
      <h1 class="text-2xl font-bold mb-4">Scan Results</h1>
      <p class="mb-2">Path: {{ path }}</p>

      <!-- TODO: Display scan results -->
      <!-- TODO: Show mods list with AppDataTable -->
      <!-- TODO: Cache statistics cards -->
      <!-- TODO: New/updated/removed mods sections -->
    </div>
  `,
})
export class ScanResultsPageComponent implements OnInit {
  path: string = '';

  constructor(private route: ActivatedRoute) {}

  ngOnInit(): void {
    this.route.queryParams
      .pipe(map((params) => decodeURIComponent(params['path'] || '')))
      .subscribe((path) => {
        this.path = path;
      });
  }
}

