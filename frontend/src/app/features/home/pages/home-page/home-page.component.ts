import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TreeTableModule } from 'primeng/treetable';
import { TreeNode } from 'primeng/api';
import { AnalysisService } from '@core/services/analysis.service';

@Component({
  selector: 'app-home-page',
  standalone: true,
  imports: [CommonModule, TreeTableModule],
  templateUrl: './home-page.component.html',
})
export class HomePageComponent implements OnInit {
  protected modTree: TreeNode<any>[] = [];

  constructor(
    private readonly analysisService: AnalysisService,
  ) {
  }

  ngOnInit(): void {
    this.analysisService.quickScan("/home/kolterdyx/GolandProjects/ModpackGraph/mods")
      .subscribe((result) => {
        result.mods.forEach(mod => {
          this.modTree.push({
            data: {
              name: mod.name,
              version: mod.version,
              id: mod.id,
            },
            children: mod.dependencies.map(dep => ({
              data: {
                name: dep.mod_id,
                version: dep.version_range,
                id: dep.mod_id,
              }
            })),
          })
        })
      });
  }
}

