import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { app } from '@wailsjs/go/models';
import Graph = app.Graph;
import { debounceTime, Subject } from 'rxjs';
import ForceGraph3D, { ForceGraph3DInstance } from '3d-force-graph';
import { ToggleButton } from 'primeng/togglebutton';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-interactive-three-tab',
  imports: [
    ToggleButton,
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

  protected showIcons: boolean = false;

  ngOnInit(): void {
    this.regenerate$
      .pipe(debounceTime(500))
      .subscribe(() => {
          this.graph = new ForceGraph3D(this.graphElement.nativeElement)
            .linkLabel('label')
            .linkWidth(1)
            .backgroundColor("#000")
            .linkVisibility(true)
            .linkColor(() => "#727272")
            .linkDirectionalArrowLength(6)

          // if (this.showIcons) {
          //
          //   this.graph
          //     .nodeThreeObject((node: Node & NodeObject, ctx) => {
          //       if (!node.id || !node.x || !node.y) {
          //         return;
          //       }
          //       if (!node?.icon) {
          //         if (!node?.name) {
          //           return;
          //         }
          //         const label = node.name;
          //         const fontSize = 2;
          //         ctx.font = `${fontSize}px Sans-Serif`;
          //         const textWidth = ctx.measureText(label).width;
          //         const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.3); // some padding
          //
          //         ctx.fillStyle = node?.color || 'white';
          //         ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, bckgDimensions[0], bckgDimensions[1]);
          //
          //         ctx.textAlign = 'center';
          //         ctx.textBaseline = 'middle';
          //         ctx.fillStyle = 'black';
          //         ctx.fillText(label, node.x, node.y);
          //       } else {
          //         const size = 12
          //         const img = new Image();
          //         img.src = node.icon
          //         ctx.drawImage(img, node.x - size / 2, node.y - size / 2, size, size);
          //       }
          //     })
          //     .nodePointerAreaPaint((node, color, ctx) => {
          //       if (!node.id || !node.x || !node.y) {
          //         return;
          //       }
          //       const size = 12;
          //       ctx.fillStyle = color;
          //       ctx.fillRect(node.x - size / 2, node.y - size / 2, size, size); // draw square as pointer trap
          //     })
          // }

          this.graph.graphData(this.data ?? {nodes: [], links: []})
            .cooldownTicks(100)

          this.graph.onEngineStop(()=> this.graph?.zoomToFit(400))
        }
      )

  }

}
