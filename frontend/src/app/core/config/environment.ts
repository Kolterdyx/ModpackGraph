export interface Environment {
  production: boolean;
  featureFlags: {
    persistState: boolean;
    enableExperimentalFeatures: boolean;
    debugMode: boolean;
  };
  wails: {
    timeout: number;
    retryAttempts: number;
  };
  cache: {
    maxRecentPaths: number;
  };
}

export const environment: Environment = {
  production: false,
  featureFlags: {
    persistState: true,
    enableExperimentalFeatures: true,
    debugMode: true,
  },
  wails: {
    timeout: 30000,
    retryAttempts: 3,
  },
  cache: {
    maxRecentPaths: 10,
  },
};

