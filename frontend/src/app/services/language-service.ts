import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class LanguageService {
  private currentLanguage: string = 'en';

  setLanguage(lang: string): void {
    this.currentLanguage = lang;
    window.location.href = window.location.href.replace(/wails:\/\/wails\/[a-z]{2}\//, `wails://wails/${lang}/`);
  }

  getLanguage(): string {
    return this.currentLanguage;
  }
}
