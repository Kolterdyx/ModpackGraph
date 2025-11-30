import { ApplicationConfig, provideBrowserGlobalErrorListeners, provideZoneChangeDetection } from '@angular/core';

import { providePrimeNG } from 'primeng/config';
import DefaultTheme from '@primeuix/themes/aura';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { definePreset } from '@primeuix/themes';

const Theme = definePreset({
  options: {
    darkModeSelector: '.dark-theme',
  },
  components: {
    toast: {
      colorScheme: {
        dark: {
          info: {
            background: '{surface.900}',
          }
        }
      }
    },
  }
}, DefaultTheme);

export const appConfig: ApplicationConfig = {
  providers: [
    provideAnimationsAsync(),
    providePrimeNG({
      theme: {
        preset: Theme,
        options: {
          darkModeSelector: '.dark-theme',
        }
      }
    }),
    provideBrowserGlobalErrorListeners(),
    provideZoneChangeDetection({eventCoalescing: true}),
  ]
};
