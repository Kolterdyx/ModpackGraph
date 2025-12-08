import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-scan-results-page',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './scan-results-page.component.html',
})
export class ScanResultsPageComponent {
  path$: Observable<string>;

  constructor(private route: ActivatedRoute) {
    this.path$ = this.route.queryParams.pipe(
      map((params) => decodeURIComponent(params['path'] || ''))
    );
  }
}

