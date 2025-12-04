import { Injectable } from '@angular/core';
import { Environment } from '@wailsjs/runtime/runtime';

@Injectable({
  providedIn: 'root',
})
export class WailsRuntimeService {
  async getEnvironmentInfo() {
    return Environment();
  }

  // Add more runtime utilities as needed
}

