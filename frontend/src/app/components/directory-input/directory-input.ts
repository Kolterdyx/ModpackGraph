import { Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { Button } from "primeng/button";
import { AbstractControl, FormsModule, ReactiveFormsModule } from "@angular/forms";
import { InputGroup } from "primeng/inputgroup";
import { InputGroupAddon } from "primeng/inputgroupaddon";
import { InputText } from "primeng/inputtext";
import { OpenDirectoryDialog } from '@wailsjs/go/app/App';
import { Tooltip } from 'primeng/tooltip';

@Component({
  selector: 'app-directory-input',
  imports: [
    Button,
    FormsModule,
    InputGroup,
    InputGroupAddon,
    InputText,
    ReactiveFormsModule,
    Tooltip
  ],
  templateUrl: './directory-input.html',
  styleUrl: './directory-input.scss',
})
export class DirectoryInput implements OnInit {
  @Input() control: AbstractControl | null = null;
  @Input() placeholder: string = 'Select directory';
  @Input() title: string = 'Select Directory';

  @ViewChild('dirInput') dirInput!: ElementRef<HTMLInputElement>

  @Output() folderSelected: EventEmitter<string> = new EventEmitter<string>();

  protected value: string = '';

  ngOnInit() {
    if (this.control) {
      this.value = this.control.getRawValue();
    }
  }

  protected async onSelectFolder() {
    try {
      this.value = await OpenDirectoryDialog({
        title: this.title,
        defaultDirectory: this.control?.getRawValue(),
      });
      if (this.value) {
        this.control?.setValue(this.value);
        this.folderSelected.emit(this.value);
      }
    } catch (error) {
      console.error("Error selecting folder:", error);
    }
  }
}
