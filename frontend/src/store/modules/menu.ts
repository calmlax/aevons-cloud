import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import type { MenuNode } from '@/types/menu';
import { fetchMenuTree } from '@/api/menu';

export const useMenuStore = defineStore('menu', () => {
  const menuTree = ref<MenuNode[]>([]);
  const menuLoading = ref(false);
  const menuInitialized = ref(false);
  const menuError = ref<string | null>(null);

  let pendingRequest: Promise<MenuNode[]> | null = null;

  const loadMenu = async (force = false) => {
    if (!force && menuInitialized.value && menuTree.value.length > 0) {
      return menuTree.value;
    }

    if (!force && pendingRequest) {
      return pendingRequest;
    }

    menuLoading.value = true;
    menuError.value = null;

    pendingRequest = fetchMenuTree()
      .then((items) => {
        menuTree.value = items;
        menuInitialized.value = true;
        return items;
      })
      .catch((error: unknown) => {
        menuTree.value = [];
        menuInitialized.value = true;
        menuError.value = error instanceof Error ? error.message : 'menu.fetchFailed';
        throw error;
      })
      .finally(() => {
        menuLoading.value = false;
        pendingRequest = null;
      });

    return pendingRequest;
  };

  const initializeMenu = () => loadMenu(false);
  const refreshMenu = () => loadMenu(true);
  
  const resetMenu = () => {
    menuTree.value = [];
    menuLoading.value = false;
    menuInitialized.value = false;
    menuError.value = null;
  };

  const menuReady = computed(() => menuInitialized.value && !menuLoading.value);

  return {
    menuTree,
    menuLoading,
    menuError,
    menuReady,
    initializeMenu,
    refreshMenu,
    resetMenu,
  };
});
