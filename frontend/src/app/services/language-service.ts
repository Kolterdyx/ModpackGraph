import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class LanguageService {

  currentLanguage: string = window.location.href.match(/wails:\/\/wails\/([a-z]{2})\//)?.[1] || 'en';

  setLanguage(lang: string): void {
    window.location.href = window.location.href.replace(/wails:\/\/wails\/[a-z]{2}\//, `wails://wails/${lang}/`);
  }
}
