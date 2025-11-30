import { Component, Input } from '@angular/core';
import { app } from '@wailsjs/go/models';
import Graph = app.Graph;

@Component({
  selector: 'app-interactive-three-tab',
  imports: [],
  templateUrl: './interactive-three-tab.html',
  styleUrl: './interactive-three-tab.scss',
})
export class InteractiveThreeTab {
  @Input() graphData?: Graph;

}
