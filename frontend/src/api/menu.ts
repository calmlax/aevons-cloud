import { hydrateMenuTree, type MenuNode, type MenuSourceNode } from '@/utils/menu';
import { request } from '@/utils/request';

const wait = (duration: number) => new Promise((resolve) => window.setTimeout(resolve, duration));

const commonMenuNodes: MenuSourceNode[] = [
  {
    key: 'dashboard',
    path: '/dashboard',
    component: 'dashboard/DashboardPage',
    meta: {
      titleKey: 'menu.dashboard',
      icon: 'dashboard',
    },
  },
];

// ===============================================
// The backend /auth/v1/routers already filters by role
// so we no longer need local `filterMenuByPermission`
// ===============================================

export const fetchMenuTree = async (): Promise<MenuNode[]> => {
  // Pull authenticated dynamic menu payload from API endpoint
  const dynamicMenuTree = await request<MenuSourceNode[]>({ url: '/auth/v1/routers', method: 'get' });
  
  // Directly hydrate since the backend guarantees permitted routes only
  return hydrateMenuTree([...commonMenuNodes, ...(dynamicMenuTree || [])]);
};