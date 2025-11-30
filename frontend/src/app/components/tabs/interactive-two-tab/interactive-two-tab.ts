import { Component, Input } from '@angular/core';
import { app } from '@wailsjs/go/models';
import Graph = app.Graph;

@Component({
  selector: 'app-interactive-two-tab',
  imports: [],
  templateUrl: './interactive-two-tab.html',
  styleUrl: './interactive-two-tab.scss',
})
export class InteractiveTwoTab {
  @Input() graphData?: Graph;

}
