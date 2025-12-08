import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class LanguageService {

  setLanguage(lang: string): void {
    window.location.href = window.location.href.replace(/wails:\/\/wails\/[a-z]{2}\//, `wails://wails/${lang}/`);
  }

  getCurrentLanguage(): string {
    return window.location.href.match(/wails:\/\/wails\/([a-z]{2})\//)?.[1] || 'en';
  }
}
