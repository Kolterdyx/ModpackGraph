import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { app } from '@wailsjs/go/models';
import ForceGraph, { NodeObject } from 'force-graph';
import { debounceTime, Subject } from 'rxjs';
import { ToggleButton } from 'primeng/togglebutton';
import { FormsModule } from '@angular/forms';
import { converter, parse } from 'culori';
import Graph = app.Graph;
import Node = app.Node;

const toOklch = converter<'oklch'>('oklch');

function darker(color: string, amount: number = 0.1): string {
  const c = toOklch(parse(color));
  if (!c) {
    return color;
  }
  return `oklch(${Math.max(0, c.l - amount)} ${c.c} ${c.h})`;
}

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
    for (const node of data.nodes) {
      if (node.id && node.icon) {
        this.images[node.id] = new Image()
        this.images[node.id].src = node.icon;
      }
    }
    this.regenerate$.next()
  }

  private images: { [key: string]: HTMLImageElement } = {};

  private graph?: ForceGraph

  protected regenerate$: Subject<void> = new Subject();

  private data?: Graph;

  protected showIcons: boolean = false;

  private resizeObserver?: ResizeObserver;

  constructor() {
  }

  ngOnInit(): void {

    this.resizeObserver = new ResizeObserver((entries, observer) => {
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
          this.graph = new ForceGraph(this.graphElement.nativeElement)
            .width(rect.width)
            .height(rect.height)
            .linkLabel('label')
            .linkWidth(1)
            .backgroundColor("#000")
            .linkVisibility(true)
            .linkColor(() => "#727272")
            .linkDirectionalArrowLength(6)

          if (this.showIcons) {
            const size = 12;
            this.graph
              .nodeCanvasObject((node: Node & NodeObject, ctx) => {
                if (!node.id || !node.x || !node.y) {
                  return;
                }
                if (!node?.icon) {
                  // draw circle
                  ctx.beginPath()
                  ctx.arc(node.x, node.y, size / 2, 0, 2 * Math.PI, false);
                  ctx.fillStyle = node?.color ?? "#fff";
                  ctx.fill();
                  ctx.lineWidth = size / 15;
                  ctx.strokeStyle = darker(node.color ?? 'white');
                  ctx.stroke();
                } else {
                  ctx.drawImage(this.images[node.id], node.x - size / 2, node.y - size / 2, size, size);
                }
              })
              .nodePointerAreaPaint((node, color, ctx) => {
                if (!node.id || !node.x || !node.y) {
                  return;
                }
                ctx.fillStyle = color;
                ctx.fillRect(node.x - size / 2, node.y - size / 2, size, size); // draw square as pointer trap
              })
          }
          this.graph.graphData(this.data ?? {nodes: [], links: []})
        }
      )
    this.regenerate$.next()
  }

}
