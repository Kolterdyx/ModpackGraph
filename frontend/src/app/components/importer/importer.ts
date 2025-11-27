import { Component, ElementRef, ViewChild } from '@angular/core';
import { Button } from 'primeng/button';
import { GenerateDependencyGraph, OpenSelectFolderDialog } from '@wailsjs/go/app/App';
import { MessageService } from 'primeng/api';
import { Toast } from 'primeng/toast';
import { Tooltip } from 'primeng/tooltip';
import SvgPanZoom from 'svg-pan-zoom';

@Component({
  selector: 'app-importer',
  templateUrl: './importer.html',
  styleUrl: './importer.scss',
  imports: [
    Button,
    Toast,
    Tooltip,
  ],
  providers: [MessageService],
})
export class Importer {

  @ViewChild('svgContainer', {static: true}) svgContainer!: ElementRef<HTMLDivElement>;

  constructor(
    private readonly messageService: MessageService,
  ) {
  }

  protected svgData?: string;

  protected selectedPath: string = '';

  private panZoomInstance: any;

  protected async onGenerate() {
    this.messageService.add({severity: 'info', summary: 'Generating graph', detail: `Generating graph...`});
    try {
      if (!this.selectedPath) {
        this.messageService.add({severity: 'error', summary: 'No folder selected', detail: 'Please select a folder first.'});
        return;
      }
      this.svgData = await GenerateDependencyGraph(this.selectedPath);

      // Clear previous SVG if exists
      this.svgContainer.nativeElement.innerHTML = '';

      // Insert inline SVG
      this.svgContainer.nativeElement.innerHTML = this.svgData;

      // Initialize svg-pan-zoom
      const svgEl = this.svgContainer.nativeElement.querySelector('svg');
      if (svgEl) {
        // Destroy previous instance if exists
        if (this.panZoomInstance) this.panZoomInstance.destroy();

        this.panZoomInstance = SvgPanZoom(svgEl, {
          zoomEnabled: true,
          controlIconsEnabled: true,
          fit: true,
          center: true,
          minZoom: 0.5,
          maxZoom: 10,
        });
      }

      this.messageService.add({severity: 'info', summary: 'Graph generated', detail: `Graph generated successfully.`});
    } catch (error) {
      this.messageService.add({severity: 'error', summary: 'Something went wrong.', detail: `Error: ${error}`});
      console.error("Error generating graph:", error);
    }
  }

  protected async onSelectFolder() {
    try {
      this.selectedPath = await OpenSelectFolderDialog(this.selectedPath);
      this.messageService.add({severity: 'info', summary: 'Folder selected', detail: `Selected folder: ${this.selectedPath}`});
    } catch (error) {
      this.messageService.add({severity: 'error', summary: 'Something went wrong.', detail: `Error: ${error}`});
      console.error("Error selecting folder:", error);
    }
  }
}
