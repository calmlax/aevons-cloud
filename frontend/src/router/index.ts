import { createRouter, createWebHistory } from 'vue-router';

import AdminLayout from '../layout/AdminLayout.vue';
import AuthLayout from '../layout/AuthLayout.vue';
import { useAuthStore } from '../store/modules/auth';
import { useMenuStore } from '../store/modules/menu';
import PublicLayout from '../layout/PublicLayout.vue';
import type { MenuNode } from '../types/menu';

const viewModules = import.meta.glob('../views/**/*.vue');

const resolveViewModule = (component: string) => {
  const normalized = component.replace(/^views\//, '').replace(/^\/+/, '');
  const candidates = normalized.endsWith('.vue')
    ? [normalized]
    : [
        `${normalized}.vue`,
        `${normalized}/index.vue`,
        `${normalized.split('/').slice(0, -1).join('/')}/index.vue`,
      ];

  return candidates
    .filter((path) => path !== '.vue' && path !== '/index.vue')
    .map((path) => viewModules[`../views/${path}`])
    .find(Boolean);
};

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/auth',
      component: AuthLayout,
      meta: {
        guestOnly: true,
      },
      children: [
        {
          path: '',
          redirect: '/login',
        },
        {
          path: '/login',
          name: 'auth-login',
          component: () => import('../views/auth/LoginPage.vue'),
          meta: {
            guestOnly: true,
            titleKey: 'auth.loginTitle',
          },
        },
        {
          path: 'register',
          name: 'auth-register',
          component: () => import('../views/auth/RegisterPage.vue'),
          meta: {
            guestOnly: true,
            titleKey: 'auth.registerTitle',
          },
        },
        {
          path: 'forgot-password',
          name: 'auth-forgot-password',
          component: () => import('../views/auth/ForgotPasswordPage.vue'),
          meta: {
            guestOnly: true,
            titleKey: 'auth.resetTitle',
          },
        },
      ],
    },
    {
      path: '/result',
      component: AuthLayout,
      children: [
        {
          path: ':scene',
          name: 'auth-result',
          component: () => import('../views/auth/AuthResultPage.vue'),
          meta: {
            titleKey: 'result.pageTitle',
          },
        },
      ],
    },
    {
      path: '/403',
      name: 'permission-denied',
      component: () => import('../views/misc/PermissionDeniedPage.vue'),
      meta: {
        titleKey: 'permission.pageTitle',
      },
    },
    {
      path: '/oauth2/authorize',
      name: 'oauth2-authorize',
      component: () => import('../views/auth/AuthorizePage.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    // ---- 门户架构区 (Public Portal) ----
    {
      path: '/',
      component: PublicLayout,
      children: [
        {
          path: '',
          name: 'portal-home',
          component: () => import('../views/public/HomePage.vue'),
          meta: { titleKey: 'portal.home' },
        },
        {
          path: 'products',
          name: 'portal-products',
          component: () => import('../views/public/ProductsPage.vue'),
          meta: { titleKey: 'portal.products' },
        },
        {
          path: 'about',
          name: 'portal-about',
          component: () => import('../views/public/AboutPage.vue'),
          meta: { titleKey: 'portal.about' },
        },
        {
          path: 'help',
          name: 'portal-help',
          component: () => import('../views/public/HelpPage.vue'),
          meta: { titleKey: 'portal.help' },
        },
      ],
    },
    // ---- 后台管理架构区 (Admin Panel) ----
    {
      path: '/',
      name: 'admin',
      component: AdminLayout,
      meta: {
        requiresAuth: true,
      },
      children: [
        // {
        //   path: '',
        //   name: 'index',
        //   component: () => import('../views/dashboard/DashboardPage.vue'),
        //   meta: {
        //     titleKey: 'menu.dashboard',
        //     menuKey: 'dashboard',
        //   },
        // },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('../views/dashboard/DashboardPage.vue'),
          meta: {
            titleKey: 'menu.dashboard',
            menuKey: 'dashboard',
          },
        },
        {
          path: 'account/profile',
          name: 'account-profile',
          component: () => import('../views/account/AccountProfilePage.vue'),
          meta: {
            titleKey: 'account.profile',
            breadcrumbKeys: ['account.settings', 'account.profile'],
          },
        },
        {
          path: 'account/password',
          name: 'account-password',
          component: () => import('../views/account/AccountPasswordPage.vue'),
          meta: {
            titleKey: 'account.password',
            breadcrumbKeys: ['account.settings', 'account.password'],
          },
        },
        {
          path: 'notifications',
          name: 'notifications',
          component: () => import('../views/misc/NotificationsPage.vue'),
          meta: {
            titleKey: 'notifications.title',
          },
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      component: AuthLayout,
      children: [
        {
          path: '',
          name: 'not-found-fallback',
          component: () => import('../views/auth/AuthResultPage.vue'),
        },
      ],
    },
  ],
  scrollBehavior: () => ({ top: 0 }),
});

router.beforeEach(async (to) => {
  const authStore = useAuthStore();
  authStore.initializeAuth();

  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const guestOnly = to.matched.some((record) => record.meta.guestOnly);

  if (requiresAuth && !authStore.isAuthenticated) {
    return {
      path: '/login',
      query: { redirect: to.fullPath },
    };
  }

  if (guestOnly && authStore.isAuthenticated) {
    return authStore.getAuthHomePath();
  }

  // Dynamic menu routing injection
  if (authStore.isAuthenticated) {
    const menuStore = useMenuStore();
    if (!menuStore.menuReady) {
      try {
        const menus = await menuStore.initializeMenu();
        let added = false;

        const walk = (nodes: MenuNode[]) => {
          for (const node of nodes) {
            if (node.component && node.path && node.key !== 'dashboard') {
              if (!router.hasRoute(node.key)) {
                const importFn = resolveViewModule(node.component);
                if (importFn) {
                  router.addRoute('admin', {
                    path: node.path.replace(/^\//, ''),
                    name: node.key,
                    component: importFn,
                    meta: {
                      titleKey: node.meta?.titleKey,
                      menuKey: node.key,
                      permissions: node.meta?.permissions,
                      activeMenu: node.meta?.activeMenu,
                    },
                  });
                  added = true;
                } else {
                  console.warn(`[Router] Missing view component: ${node.component}`);
                }
              }
            }
            if (node.children?.length) {
              walk(node.children);
            }
          }
        };

        walk(menus);

        if (added) {
          // Re-evaluate current navigation using the raw fullPath
          // This ensures Vue Router ignores previous CatchAll match name and resolves against the new route table
          return { path: to.fullPath, replace: true, query: to.query, hash: to.hash };
        }
      } catch (error) {
        console.error('Failed to initialize dynamic menu routes', error);
      }
    }
  }

  const permissionList = to.matched.flatMap((record) => {
    const permissions = record.meta.permissions;
    if (typeof permissions === 'string') {
      return [permissions];
    }
    return Array.isArray(permissions) ? permissions.map(String) : [];
  });

  if (permissionList.length > 0 && !authStore.hasPermission(permissionList)) {
    return {
      path: '/403',
      query: { from: to.fullPath },
    };
  }

  return true;
});

export default router;
