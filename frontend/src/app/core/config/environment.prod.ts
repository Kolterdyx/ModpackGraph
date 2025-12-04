import { Environment } from './environment';

export const environment: Environment = {
  production: true,
  featureFlags: {
    persistState: true,
    enableExperimentalFeatures: false,
    debugMode: false,
  },
  wails: {
    timeout: 60000,
    retryAttempts: 1,
  },
  cache: {
    maxRecentPaths: 20,
  },
};

