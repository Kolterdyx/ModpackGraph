import { Component } from '@angular/core';
import { Toast } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { FormBuilder, FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { DirectoryInput } from '@components/directory-input/directory-input';
import { Select } from 'primeng/select';
import { SvgGraphTab } from '@components/tabs/svg-graph-tab/svg-graph-tab';
import { Form } from '@/app/models/form';
import { SelectButton } from 'primeng/selectbutton';
import { InteractiveTwoTab } from '@components/tabs/interactive-two-tab/interactive-two-tab';
import { InteractiveThreeTab } from '@components/tabs/interactive-three-tab/interactive-three-tab';
import { GenerateDependencyGraph } from '@wailsjs/go/app/App';
import { app } from '@wailsjs/go/models';
import Graph = app.Graph;
import GraphGenerationOptions = app.GraphGenerationOptions;
import Layout = app.Layout;
import { Button } from 'primeng/button';

interface SelectValue {
  label: string;
  value: string;
  icon?: string;
}

@Component({
  selector: 'app-root',
  templateUrl: './app.html',
  imports: [
    Toast,
    ReactiveFormsModule,
    DirectoryInput,
    Select,
    SvgGraphTab,
    SelectButton,
    FormsModule,
    InteractiveTwoTab,
    InteractiveThreeTab,
    Button,
  ],
  providers: [
    MessageService,
  ],
  styleUrl: './app.scss'
})
export class App {

  protected formGroup?: FormGroup<Form<GraphGenerationOptions>>;
  protected layoutOptions: string[] = [];
  protected currentTab: string = '2Di';
  protected items: SelectValue[] = [
    {
      label: 'Graphviz',
      value: 'graphviz',
      icon: 'pi pi-image',
    },
    {
      label: '2D Interactive',
      value: '2Di',
      icon: 'pi pi-stop',
    },
    {
      label: '3D Interactive',
      value: '3Di',
      icon: 'pi pi-box',
    }
  ];
  protected graphData?: Graph;

  constructor(
    private readonly fb: FormBuilder,
    private readonly messageService: MessageService,
  ) {
    for (let l in Layout) {
      this.layoutOptions.push(l)
    }
    this.formGroup = this.fb.group<Form<GraphGenerationOptions>>({
      path: new FormControl<string>('', [Validators.required]),
      layout: new FormControl<Layout>(Layout.fdp, [Validators.required]),
    })
  }


  protected async onGenerate() {
    const graphOptions = this.formGroup?.value
    if (!graphOptions) {
      return;
    }
    this.messageService.add({severity: 'info', summary: 'Generating graph', detail: `Generating graph...`});
    try {
      this.graphData = await GenerateDependencyGraph(graphOptions as GraphGenerationOptions);
      console.log(this.graphData);
      this.messageService.add({severity: 'success', summary: 'Graph generated', detail: `Graph generated successfully.`});
    } catch (error) {
      this.messageService.add({severity: 'error', summary: 'Something went wrong.', detail: `Error: ${error}`});
      console.error("Error generating graph:", error);
    }
  }
}
