import { Component, Input, OnChanges, OnInit } from '@angular/core';
import { app } from '@wailsjs/go/models';
import { GraphDisplayOptions, ListDisplayOptions } from '@/app/models/graph-display-options';
import Graph = app.Graph;
import { DataView } from 'primeng/dataview';
import { ScrollPanel } from 'primeng/scrollpanel';
import { Tag } from 'primeng/tag';

interface Mod {
  id: string;
  name: string;
  presentVersion: string;
  present: boolean;
  required: boolean;
  requiredVersion: string;
  iconURL?: string;
}

@Component({
  selector: 'app-list-tab',
  imports: [
    DataView,
    ScrollPanel,
    Tag
  ],
  templateUrl: './list-tab.html',
  styleUrl: './list-tab.scss',
})
export class ListTab implements OnChanges, OnInit {
  @Input() graphData?: Graph
  @Input() options?: ListDisplayOptions

  protected mods: Mod[] = [];

  ngOnInit() {
    this.ngOnChanges();
  }

  ngOnChanges() {
    if (!this.graphData) {
      return;
    }
    this.mods = []
    for (const node of this.graphData.nodes) {
      const isRequired = this.graphData.links.some(edge => edge.target === node.id && edge.required);
      const isPresent = node?.present ?? false;
      this.mods.push({
        id: (node.id ?? '').toString(),
        name: node.name ?? (node.id ?? '').toString(),
        presentVersion: node.presentVersion ?? '',
        present: isPresent,
        required: isRequired,
        requiredVersion: node.requiredVersion ?? '',
        iconURL: node.icon,
      });
      this.mods.sort((a, b) => {
        // Missing required mods first
        if (a.required && !a.present && !(b.required && !b.present)) {
          return -1;
        }
        if (b.required && !b.present && !(a.required && !a.present)) {
          return 1;
        }
        // Then missing optional mods
        if (!a.present && b.present) {
          return -1;
        }
        if (!b.present && a.present) {
          return 1;
        }
        // Then by name
        return a.name.localeCompare(b.name);
      } );
    }
  }

  protected getStatusSeverity(id: string) {
    const mod = this.mods.find(m => m.id === id);
    if (!mod) {
      return 'warn';
    }
    if (mod.required && !mod.present) {
      return 'danger';
    }
    if (mod.present) {
      return 'success';
    }
    return 'warn';
  }

  protected readonly $localize = $localize;
}
