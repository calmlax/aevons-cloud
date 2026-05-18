import { defineStore } from 'pinia';
import { ref } from 'vue';

/** 布局模式：sidebar = 传统左侧菜单，topnav = 一级菜单在顶栏 */
export type LayoutMode = 'sidebar' | 'topnav';

const STORAGE_KEY = 'aevo-layout-mode';

export const useLayoutStore = defineStore('layout', () => {
  const layoutMode = ref<LayoutMode>(
    (window.localStorage.getItem(STORAGE_KEY) as LayoutMode) ?? 'topnav',
  );

  const setLayoutMode = (mode: LayoutMode) => {
    layoutMode.value = mode;
    window.localStorage.setItem(STORAGE_KEY, mode);
  };

  return { layoutMode, setLayoutMode };
});
