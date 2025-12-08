import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TreeTableModule } from 'primeng/treetable';
import { TreeNode } from 'primeng/api';
import { WailsAppService } from '@core/services/wails';

@Component({
  selector: 'app-home-page',
  standalone: true,
  imports: [CommonModule, TreeTableModule],
  templateUrl: './home-page.component.html',
})
export class HomePageComponent implements OnInit {
  protected modTree: TreeNode<any>[] = [];

  constructor(
    private readonly wailsAppService: WailsAppService
  ) {}

  ngOnInit(): void {
    this.wailsAppService.scanModpack("/home/kolterdyx/GolandProjects/ModpackGraph/mods").subscribe((data) => {
      for (const mod of data.mods) {
        this.modTree.push({
          data: mod,
          children: [],
        });
      }
    })
  }

}

