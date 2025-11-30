import { Component, Input, OnInit } from '@angular/core';
import { SVGViewer } from "@components/tabs/svg-graph-tab/svgviewer/svgviewer";
import { GenerateDependencyGraphSVG } from '@wailsjs/go/app/App';
import { app } from '@wailsjs/go/models';
import { debounceTime, Subject } from 'rxjs';
import Graph = app.Graph;

@Component({
  selector: 'app-svg-graph-tab',
  imports: [
    SVGViewer
  ],
  templateUrl: './svg-graph-tab.html',
  styleUrl: './svg-graph-tab.scss',
})
export class SvgGraphTab implements OnInit {

  @Input() set graphData(graph: Graph | undefined) {
    if (!graph) {
      return;
    }
    this.graphData$.next(graph)
  }

  protected svgData?: string;

  private graphData$: Subject<Graph> = new Subject();

  ngOnInit(): void {
    this.graphData$
      .pipe(debounceTime(500))
      .subscribe(graphData => {
        console.log(graphData);
        GenerateDependencyGraphSVG(graphData)
          .then((svgData) => {
            console.log('svg generated');
            this.svgData = svgData;
          })
      });
  }

}
