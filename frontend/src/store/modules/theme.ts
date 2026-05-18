import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

export type ThemeMode = 'light' | 'dark' | 'system';
type ResolvedTheme = 'light' | 'dark';

export const useThemeStore = defineStore('theme', () => {
  const THEME_STORAGE_KEY = 'aevo-theme-mode';
  const themeMode = ref<ThemeMode>('light');
  const systemPrefersDark = ref(false);

  let initialized = false;
  let mediaQuery: MediaQueryList | null = null;

  const resolvedTheme = computed<ResolvedTheme>(() => {
    if (themeMode.value === 'system') {
      return systemPrefersDark.value ? 'dark' : 'light';
    }

    return themeMode.value;
  });

  const isThemeMode = (value: string | null): value is ThemeMode => {
    return value === 'light' || value === 'dark' || value === 'system';
  };

  const syncArcoTheme = (nextTheme: ResolvedTheme) => {
    if (typeof document === 'undefined' || !document.body) {
      return;
    }

    if (nextTheme === 'dark') {
      document.body.setAttribute('arco-theme', 'dark');
      return;
    }

    document.body.removeAttribute('arco-theme');
  };

  const applyTheme = () => {
    if (typeof document === 'undefined') {
      return;
    }

    const nextTheme = resolvedTheme.value;

    document.documentElement.dataset.theme = nextTheme;
    document.documentElement.style.colorScheme = nextTheme;

    if (document.body) {
      syncArcoTheme(nextTheme);
      return;
    }

    window.requestAnimationFrame(() => syncArcoTheme(nextTheme));
  };

  const persistTheme = () => {
    if (typeof window === 'undefined') {
      return;
    }

    window.localStorage.setItem(THEME_STORAGE_KEY, themeMode.value);
  };

  const handleSystemThemeChange = (event: MediaQueryListEvent) => {
    systemPrefersDark.value = event.matches;

    if (themeMode.value === 'system') {
      applyTheme();
    }
  };

  const initializeTheme = () => {
    if (initialized || typeof window === 'undefined') {
      return;
    }

    mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    systemPrefersDark.value = mediaQuery.matches;

    const storedTheme = window.localStorage.getItem(THEME_STORAGE_KEY);
    if (isThemeMode(storedTheme)) {
      themeMode.value = storedTheme;
    }

    if (typeof mediaQuery.addEventListener === 'function') {
      mediaQuery.addEventListener('change', handleSystemThemeChange);
    } else {
      mediaQuery.addListener(handleSystemThemeChange);
    }

    applyTheme();
    initialized = true;
  };

  const setThemeMode = (nextTheme: ThemeMode) => {
    themeMode.value = nextTheme;
    persistTheme();
    applyTheme();
  };

  return {
    themeMode,
    resolvedTheme,
    initializeTheme,
    setThemeMode,
    applyTheme,
  };
});
