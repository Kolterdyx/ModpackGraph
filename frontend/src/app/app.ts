import { Component } from '@angular/core';
import { Toast } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { FormBuilder, FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { DirectoryInput } from '@components/directory-input/directory-input';
import { Form } from '@/app/models/form';
import { SelectButton } from 'primeng/selectbutton';
import { InteractiveTwoTab } from '@components/tabs/interactive-two-tab/interactive-two-tab';
import { InteractiveThreeTab } from '@components/tabs/interactive-three-tab/interactive-three-tab';
import { GenerateDependencyGraph } from '@wailsjs/go/app/App';
import { app } from '@wailsjs/go/models';
import { Button } from 'primeng/button';
import { ToggleButton } from 'primeng/togglebutton';
import { DisplayOptions } from '@/app/models/display-options';
import { Slider } from 'primeng/slider';
import { ScrollPanel } from 'primeng/scrollpanel';
import { LanguageService } from '@services/language-service';
import { Select } from 'primeng/select';
import Graph = app.Graph;
import GraphGenerationOptions = app.GraphGenerationOptions;
import Layout = app.Layout;

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
    SelectButton,
    FormsModule,
    InteractiveTwoTab,
    InteractiveThreeTab,
    Button,
    ToggleButton,
    Slider,
    ScrollPanel,
    Select,
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
      label: $localize`2D Interactive`,
      value: '2Di',
      icon: 'pi pi-stop',
    },
    {
      label: $localize`3D Interactive`,
      value: '3Di',
      icon: 'pi pi-box',
    }
  ];
  protected graphData?: Graph;

  protected displayOptions: DisplayOptions = {
    showIcons: true,
  };
  protected displayForm?: FormGroup<Form<DisplayOptions>>;
  protected languageOptions = [
    {
      label: 'English',
      value: 'en',
      countryCode: 'gb',
    },
    {
      label: 'Español',
      value: 'es',
      countryCode: 'es',
    },
  ];
  protected currentLanguage: string = 'en';

  constructor(
    private readonly fb: FormBuilder,
    private readonly messageService: MessageService,
    private readonly langService: LanguageService,
  ) {
    this.currentLanguage = this.langService.getCurrentLanguage();
    for (let l in Layout) {
      this.layoutOptions.push(l)
    }
    this.formGroup = this.fb.group<Form<GraphGenerationOptions>>({
      path: new FormControl<string>('', [Validators.required]),
      layout: new FormControl<Layout>(Layout.fdp, [Validators.required]),
    })
    this.displayForm = this.fb.group<Form<DisplayOptions>>({
      showIcons: new FormControl<boolean>(true, [Validators.required]),
      alphaDecay: new FormControl<number>(0.022, [Validators.required]),
      velocityDecay: new FormControl<number>(0.4, [Validators.required]),
    });
    this.displayOptions = this.displayForm.value as DisplayOptions;
    this.displayForm.valueChanges.subscribe(value => {
      this.displayOptions = value as DisplayOptions;
    });
  }


  protected async onGenerate() {
    const graphOptions = this.formGroup?.value
    if (!graphOptions) {
      return;
    }
    this.messageService.add({severity: 'info', summary: $localize`Generating graph`, detail: $localize`Generating graph...`});
    try {
      this.graphData = await GenerateDependencyGraph(graphOptions as GraphGenerationOptions);
      this.messageService.add({severity: 'success', summary: $localize`Graph generated`, detail: $localize`Graph generated successfully.`});
    } catch (error) {
      this.messageService.add({severity: 'error', summary: $localize`Something went wrong.`, detail: `Error: ${error}`});
      console.error("Error generating graph:", error);
    }
  }

  protected setLanguage(lang: string) {
    this.langService.setLanguage(lang);
  }
}
