import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';

@Component({
  selector: 'app-analysis-page',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-4">
      <h1 class="text-2xl font-bold mb-4">Analysis</h1>
      <p class="mb-4">Path: {{ path }}</p>

      <!-- TODO: Add AppTabs with three tabs -->
      <!-- TODO: Dependencies tab with dependency-view component -->
      <!-- TODO: Conflicts tab with conflict-view component -->
      <!-- TODO: Graph tab with graph-view component (2D/3D/List toggle) -->
    </div>
  `,
})
export class AnalysisPageComponent implements OnInit {
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

