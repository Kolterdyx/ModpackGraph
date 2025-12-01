import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { app } from '@wailsjs/go/models';
import ForceGraph, { LinkObject, NodeObject } from 'force-graph';
import { debounceTime, Subject } from 'rxjs';
import { FormsModule } from '@angular/forms';
import { GraphDisplayOptions } from '@/app/models/graph-display-options';
import Graph = app.Graph;
import Node = app.Node;
import Edge = app.Edge;

@Component({
  selector: 'app-interactive-two-tab',
  imports: [
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
    this.nodeMap = {};
    this.images = {};
    for (const node of data.nodes) {
      this.nodeMap[node.id!] = node;
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

  private resizeObserver?: ResizeObserver;

  private nodeMap: { [key: string]: Node } = {};

  @Input() set options(displayOptions: GraphDisplayOptions) {
    this.displayOptions = displayOptions;
    this.regenerate$.next()
  }

  private displayOptions?: GraphDisplayOptions;

  constructor() {
  }

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
          this.graph = new ForceGraph(this.graphElement.nativeElement)
            .width(rect.width)
            .height(rect.height)
            .d3AlphaDecay(0.1)
            .linkLabel((link: Pick<Edge, 'label' | 'required'> & LinkObject) => {
              return link.required ? $localize`Required: ${link.label}` : $localize`Optional: ${link.label}`;
            })
            .linkWidth(1)
            .backgroundColor("#000")
            .linkVisibility(true)
            .linkColor((link: Pick<Edge, 'label' | 'required'> & LinkObject) => {
              const target = link.target as NodeObject
              if (this.nodeMap[target?.id ?? '']?.present) {
                return "#727272";
              }
              return link.required ? "#ff0000" : "#ffcc00";
            })
            .linkDirectionalArrowLength(6)

          if (this.displayOptions?.showIcons) {
            const size = 10;
            const strokeSize = size / 15;
            this.graph
              .nodeCanvasObject((node: Node & NodeObject, ctx) => {
                if (!node.id || !node.x || !node.y) {
                  return;
                }
                if (!node?.icon) {
                  // draw circle
                  ctx.beginPath()
                  ctx.arc(node.x, node.y, size / 2, 0, 2 * Math.PI, false);
                  ctx.fillStyle = "#cfcfcf";
                  ctx.fill();
                  ctx.lineWidth = strokeSize;
                  ctx.strokeStyle = "#8e8e8e";
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
                const s = strokeSize * 2 + size;
                ctx.fillRect(node.x - s / 2, node.y - s / 2, s, s); // draw square as pointer trap
              })
          }
          this.graph.graphData(this.data ?? {nodes: [], links: []})
        }
      )
    this.regenerate$.next()
  }

}
