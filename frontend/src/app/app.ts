import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { Toast } from 'primeng/toast';
import { Button } from 'primeng/button';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    Toast,
    Button
  ],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App {

  protected isDarkTheme: boolean = true;

  constructor() {
    this.setTheme();
  }

  protected toggleTheme(): void {
    this.isDarkTheme = !this.isDarkTheme;
    this.setTheme();
  }

  private setTheme(): void {
    if (this.isDarkTheme) {
      document.body.classList.add('dark');
    } else {
      document.body.classList.remove('dark');
    }
  }
}
