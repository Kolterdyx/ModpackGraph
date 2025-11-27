import { Component, ElementRef, ViewChild } from '@angular/core';
import { Button } from 'primeng/button';
import { Toast } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { GenerateDependencyGraph } from '@wailsjs/go/app/App';
import SvgPanZoom from 'svg-pan-zoom';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { DirectoryInput } from '@components/directory-input/directory-input';
import { GraphGenerationData } from '@/app/models/graph-generation-data';

@Component({
  selector: 'app-root',
  templateUrl: './app.html',
  imports: [
    Button,
    Toast,
    ReactiveFormsModule,
    DirectoryInput,
  ],
  providers: [
    MessageService,
  ],
  styleUrl: './app.scss'
})
export class App {

  @ViewChild('svgContainer', {static: true}) svgContainer!: ElementRef<HTMLDivElement>;

  protected svgData?: string;

  private panZoomInstance: any;

  protected formGroup?: FormGroup;

  constructor(
    private readonly messageService: MessageService,
  ) {
    this.formGroup = new FormGroup({
      path: new FormControl(''),
    });
  }

  protected async onGenerate() {
    this.messageService.add({severity: 'info', summary: 'Generating graph', detail: `Generating graph...`});
    try {
      const formData = this.formGroup?.value as GraphGenerationData;
      this.svgData = await GenerateDependencyGraph(formData.path);
      this.replaceSVGElement();
      this.messageService.add({severity: 'success', summary: 'Graph generated', detail: `Graph generated successfully.`});
    } catch (error) {
      this.messageService.add({severity: 'error', summary: 'Something went wrong.', detail: `Error: ${error}`});
      console.error("Error generating graph:", error);
    }
  }

  private replaceSVGElement() {
    if (!this.svgData) return;
    this.svgContainer.nativeElement.innerHTML = this.svgData;
    const svgEl = this.svgContainer.nativeElement.querySelector('svg');
    if (svgEl) {
      if (this.panZoomInstance) this.panZoomInstance.destroy();
      this.panZoomInstance = SvgPanZoom(svgEl, {
        zoomEnabled: true,
        controlIconsEnabled: true,
        fit: true,
        center: true,
        minZoom: 0.1,
        maxZoom: 10,
      });
    }
  }

}
