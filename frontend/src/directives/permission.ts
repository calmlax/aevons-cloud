import type { App, DirectiveBinding, ObjectDirective } from 'vue';

import { useAuthStore } from '../store/modules/auth';

type PermissionBinding = string | string[] | undefined;

const applyPermission = (el: HTMLElement, binding: DirectiveBinding<PermissionBinding>) => {
  const required = binding.value;
  const authStore = useAuthStore();
  if (!authStore.hasPermission(required)) {
    el.style.display = 'none';
  } else {
    el.style.display = '';
  }
};

const permissionDirective: ObjectDirective<HTMLElement, PermissionBinding> = {
  mounted(el, binding) {
    applyPermission(el, binding);
  },
  updated(el, binding) {
    applyPermission(el, binding);
  },
};

export const registerPermissionDirective = (app: App) => {
  app.directive('permission', permissionDirective);
};