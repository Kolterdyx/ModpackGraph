import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { app } from '@wailsjs/go/models';
import { debounceTime, Subject } from 'rxjs';
import ForceGraph3D, { ForceGraph3DInstance } from '3d-force-graph';
import { FormsModule } from '@angular/forms';
import { DisplayOptions } from '@/app/models/display-options';
import Graph = app.Graph;

@Component({
  selector: 'app-interactive-three-tab',
  imports: [
    FormsModule
  ],
  templateUrl: './interactive-three-tab.html',
  styleUrl: './interactive-three-tab.scss',
})
export class InteractiveThreeTab implements OnInit {

  @ViewChild("graph", {static: true}) graphElement!: ElementRef<HTMLDivElement>;

  @Input() set graphData(data: Graph | undefined) {
    if (!data) {
      return;
    }
    this.data = data;
    this.regenerate$.next()
  }

  private graph?: ForceGraph3DInstance

  protected regenerate$: Subject<void> = new Subject();

  private data?: Graph;

  @Input() set options(displayOptions: DisplayOptions) {
    this.displayOptions = displayOptions;
    this.regenerate$.next()
  }

  private displayOptions?: DisplayOptions;

  private resizeObserver?: ResizeObserver;


  ngOnInit(): void {

    this.resizeObserver = new ResizeObserver(() => {
      if (this.graph) {
        const rect = this.graphElement.nativeElement.getBoundingClientRect();
        this.graph.width(rect.width).height(rect.height);
      }
    });

    this.resizeObserver.observe(this.graphElement.nativeElement);
    this.regenerate$
      .pipe(debounceTime(500))
      .subscribe(() => {
        const rect = this.graphElement.nativeElement.getBoundingClientRect();
        this.graph = new ForceGraph3D(this.graphElement.nativeElement)
          .width(rect.width)
          .height(rect.height)
          .linkLabel('label')
          .linkWidth(1)
          .backgroundColor("#000")
          .linkVisibility(true)
          .linkColor(() => "#727272")
          .linkDirectionalArrowLength(6)

        this.graph.graphData(this.data ?? {nodes: [], links: []})

      })
    this.regenerate$.next()
  }
}
