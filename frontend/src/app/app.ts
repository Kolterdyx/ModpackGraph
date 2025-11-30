import { Component } from '@angular/core';
import { Button } from 'primeng/button';
import { Toast } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { GenerateDependencyGraphSVG } from '@wailsjs/go/app/App';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { DirectoryInput } from '@components/directory-input/directory-input';
import * as Models from '@wailsjs/go/models';
import { Select } from 'primeng/select';
import GraphOptions = Models.app.GraphOptions;
import Layout = Models.app.Layout;
import { SVGViewer } from '@components/svgviewer/svgviewer';

type Form<T> = {
  [K in keyof T]: FormControl<T[K] | null>;
};

@Component({
  selector: 'app-root',
  templateUrl: './app.html',
  imports: [
    Button,
    Toast,
    ReactiveFormsModule,
    DirectoryInput,
    Select,
    SVGViewer,
  ],
  providers: [
    MessageService,
  ],
  styleUrl: './app.scss'
})
export class App {

  protected formGroup?: FormGroup<Form<GraphOptions>>;

  protected layoutOptions: string[] = [];

  protected svgData?: string;

  constructor(
    private readonly messageService: MessageService,
  ) {
    for (let l in Layout) {
      this.layoutOptions.push(l)
    }
    this.formGroup = new FormGroup<Form<GraphOptions>>({
      path: new FormControl<string>('', [Validators.required]),
      layout: new FormControl<Layout>(Layout.fdp, [Validators.required]),
    });
  }
  protected async onGenerate() {
    this.messageService.add({severity: 'info', summary: 'Generating graph', detail: `Generating graph...`});
    try {
      const formValue = this.formGroup?.value;
      console.log(formValue);
      this.svgData = await GenerateDependencyGraphSVG(formValue as GraphOptions);
      this.messageService.add({severity: 'success', summary: 'Graph generated', detail: `Graph generated successfully.`});
    } catch (error) {
      this.messageService.add({severity: 'error', summary: 'Something went wrong.', detail: `Error: ${error}`});
      console.error("Error generating graph:", error);
    }
  }

}
