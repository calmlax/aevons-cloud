import type { App } from 'vue';

import { registerPermissionDirective } from './permission';

export const registerGlobalDirectives = (app: App) => {
  registerPermissionDirective(app);
};