import { computed, Injectable, signal } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ThemeService {

  private isDarkTheme$: BehaviorSubject<boolean> = new BehaviorSubject(false);

  public toggleTheme(): void {
    const newTheme = !this.isDarkTheme$.getValue();
    this.isDarkTheme$.next(newTheme);
  }

  public getIsDarkTheme() {
    return this.isDarkTheme$.asObservable();
  }

}
