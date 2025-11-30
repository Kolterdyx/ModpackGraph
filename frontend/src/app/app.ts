import { Component } from '@angular/core';
import { Button } from 'primeng/button';
import { Toast } from 'primeng/toast';
import { MenuItem, MenuItemCommandEvent, MessageService } from 'primeng/api';
import { GenerateDependencyGraphSVG } from '@wailsjs/go/app/App';
import { FormBuilder, FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { DirectoryInput } from '@components/directory-input/directory-input';
import * as Models from '@wailsjs/go/models';
import { Select } from 'primeng/select';
import GraphOptions = Models.app.GraphOptions;
import Layout = Models.app.Layout;
import { SvgGraphTab } from '@components/svg-graph-tab/svg-graph-tab';
import { Form } from '@/app/models/form';
import { Menubar } from 'primeng/menubar';

@Component({
  selector: 'app-root',
  templateUrl: './app.html',
  imports: [
    Toast,
    ReactiveFormsModule,
    DirectoryInput,
    Select,
    SvgGraphTab,
    Menubar,
  ],
  providers: [
    MessageService,
  ],
  styleUrl: './app.scss'
})
export class App {

  protected formGroup?: FormGroup<Form<GraphOptions>>;

  protected layoutOptions: string[] = [];

  protected currentTab: string = 'graphviz';
  protected items: MenuItem[] = [];

  constructor(
    private readonly fb: FormBuilder,
  ) {
    for (let l in Layout) {
      this.layoutOptions.push(l)
    }
    this.formGroup = this.fb.group<Form<GraphOptions>>({
        path: new FormControl<string>('', [Validators.required]),
        layout: new FormControl<Layout>(Layout.fdp, [Validators.required]),
    })
    this.items = [
      {
        label: 'Graph',
        command: () => {
          this.currentTab = 'graphviz';
        }
      },
      {
        label: 'Menu',
        command: () => {
          this.currentTab = 'menu';
        }
      }
    ]
  }

}
