import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { app } from '@wailsjs/go/models';
import ForceGraph, { NodeObject } from 'force-graph';
import { debounceTime, Subject } from 'rxjs';
import { ToggleButton } from 'primeng/togglebutton';
import { FormsModule } from '@angular/forms';
import Graph = app.Graph;
import Node = app.Node;

@Component({
  selector: 'app-interactive-two-tab',
  imports: [
    ToggleButton,
    FormsModule
  ],
  templateUrl: './interactive-two-tab.html',
  styleUrl: './interactive-two-tab.scss',
})
export class InteractiveTwoTab implements OnInit {

  @ViewChild("graph", {static: true}) graphElement!: ElementRef<HTMLDivElement>;

  @Input() set graphData(data: Graph | undefined) {
    if (!data) {
      return;
    }
    this.data = data;
    this.regenerate$.next()
  }

  private graph?: ForceGraph

  protected regenerate$: Subject<void> = new Subject();

  private data?: Graph;

  protected showIcons: boolean = false;

  constructor() {
  }

  private resizeCanvas() {
    for (let canvasEl of this.graphElement.nativeElement.getElementsByTagName("canvas")) {
      canvasEl.removeAttribute("width");
      canvasEl.removeAttribute("height");
      canvasEl.style.width = "100%";
      canvasEl.style.height = "100%";
    }
  }

  ngOnInit(): void {
    this.regenerate$
      .pipe(debounceTime(500))
      .subscribe(() => {
          this.graph = new ForceGraph(this.graphElement.nativeElement)
            .linkLabel('label')
            .linkWidth(1)
            .backgroundColor("#000")
            .linkVisibility(true)
            .linkColor(() => "#727272")
            .linkDirectionalArrowLength(6)

          if (this.showIcons) {

            this.graph
              .nodeCanvasObject((node: Node & NodeObject, ctx) => {
                if (!node.id || !node.x || !node.y) {
                  return;
                }
                if (!node?.icon) {
                  if (!node?.name) {
                    return;
                  }
                  const label = node.name;
                  const fontSize = 2;
                  ctx.font = `${fontSize}px Sans-Serif`;
                  const textWidth = ctx.measureText(label).width;
                  const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.3); // some padding

                  ctx.fillStyle = node?.color || 'white';
                  ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, bckgDimensions[0], bckgDimensions[1]);

                  ctx.textAlign = 'center';
                  ctx.textBaseline = 'middle';
                  ctx.fillStyle = 'black';
                  ctx.fillText(label, node.x, node.y);
                } else {
                  const size = 12
                  const img = new Image();
                  img.src = node.icon
                  ctx.drawImage(img, node.x - size / 2, node.y - size / 2, size, size);
                }
              })
              .nodePointerAreaPaint((node, color, ctx) => {
                if (!node.id || !node.x || !node.y) {
                  return;
                }
                const size = 12;
                ctx.fillStyle = color;
                ctx.fillRect(node.x - size / 2, node.y - size / 2, size, size); // draw square as pointer trap
              })
          }

          this.graph
            .cooldownTicks(100)
            .graphData(this.data ?? {nodes: [], links: []})

          this.graph.onEngineStop(() => this.graph?.zoomToFit(400));
        }
      )

  }

}
