import { Component, Input } from '@angular/core';
import { Button } from "primeng/button";
import { SVGViewer } from "@components/svgviewer/svgviewer";
import { GenerateDependencyGraphSVG } from '@wailsjs/go/app/App';
import { MessageService } from 'primeng/api';
import { app } from '@wailsjs/go/models';
import GraphOptions = app.GraphOptions;
import { Nullable } from '@/app/models/form';

@Component({
  selector: 'app-svg-graph-tab',
  imports: [
    Button,
    SVGViewer
  ],
  templateUrl: './svg-graph-tab.html',
  styleUrl: './svg-graph-tab.scss',
})
export class SvgGraphTab {

  @Input() graphOptions?: Partial<Nullable<GraphOptions>>;

  protected svgData?: string;

  constructor(
    private readonly messageService: MessageService,
  ) {
  }

  protected async onGenerate() {
    if (!this.graphOptions) {
      return;
    }
    this.messageService.add({severity: 'info', summary: 'Generating graph', detail: `Generating graph...`});
    try {
      this.svgData = await GenerateDependencyGraphSVG(this.graphOptions as GraphOptions);
      this.messageService.add({severity: 'success', summary: 'Graph generated', detail: `Graph generated successfully.`});
    } catch (error) {
      this.messageService.add({severity: 'error', summary: 'Something went wrong.', detail: `Error: ${error}`});
      console.error("Error generating graph:", error);
    }
  }
}
