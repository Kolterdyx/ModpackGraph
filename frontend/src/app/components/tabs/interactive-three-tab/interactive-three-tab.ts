import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { app } from '@wailsjs/go/models';
import { debounceTime, Subject } from 'rxjs';
import ForceGraph3D, { ForceGraph3DInstance } from '3d-force-graph';
import { FormsModule } from '@angular/forms';
import { DisplayOptions } from '@/app/models/display-options';
import Graph = app.Graph;
import { LinkObject, NodeObject } from 'force-graph';
import Edge = app.Edge;
import Node = app.Node;

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

  private nodeMap: { [key: string]: Node } = {};

  @Input() set graphData(data: Graph | undefined) {
    if (!data) {
      return;
    }
    this.data = data;
    this.nodeMap = {};
    for (const node of data.nodes) {
      this.nodeMap[node.id!] = node;
    }
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
          .d3AlphaDecay(this.displayOptions?.alphaDecay ?? 0.0228)
          .d3VelocityDecay(this.displayOptions?.velocityDecay ?? 0.4)
          .linkLabel((link: Pick<Edge, 'label' | 'required'> & LinkObject) => {
            return link.required ? $localize`Required: ${link.label}` : $localize`Optional: ${link.label}`;
          })
          .linkWidth(1)
          .backgroundColor("#000")
          .linkVisibility(true)
          .linkColor((link: Pick<Edge, 'label' | 'required'> & LinkObject) => {
            const target = link.target as NodeObject
            console.log(target.id, this.nodeMap[target?.id ?? ''], this.nodeMap)
            if (this.nodeMap[target?.id ?? '']?.present) {
              return "#727272";
            }
            return link.required ? "#ff0000" : "#ffcc00";
          })
          .linkDirectionalArrowLength(6)

        this.graph.graphData(this.data ?? {nodes: [], links: []})

      })
    this.regenerate$.next()
  }
}
