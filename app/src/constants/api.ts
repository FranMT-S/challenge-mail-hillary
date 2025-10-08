

// windows.__ENV__ is used to get the environment variables in runtime in prod
export const API_CONFIG = {
  BASE_URL: window.__ENV__?.VITE_API_URL || import.meta.env.VITE_API_URL || 'http://localhost:8080',
  ENDPOINTS: {
    MAILS: '/api/mails/search',
  },
};

export const DEFAULT_PAGINATION = {
  PAGE: 1,
  LIMIT: 10,
};
