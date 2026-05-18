/// <reference types="vite/client" />
/// <reference types="vite-plugin-svg-icons/client" />

interface ImportMetaEnv {
  readonly VITE_APP_ENV: 'development' | 'staging' | 'production';
  readonly VITE_API_BASE_URL: string;
  readonly VITE_APP_TITLE: string;
  readonly VITE_MOCK_ENABLED: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
