export const appEnv = import.meta.env.VITE_APP_ENV as 'development' | 'staging' | 'production';
export const apiBaseUrl = import.meta.env.VITE_API_BASE_URL as string;
export const appTitle = import.meta.env.VITE_APP_TITLE as string;
export const mockEnabled = import.meta.env.VITE_MOCK_ENABLED === 'true';

export const isDev = appEnv === 'development';
export const isStaging = appEnv === 'staging';
export const isProd = appEnv === 'production';
