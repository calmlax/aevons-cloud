import type { Component } from 'vue';
import type { LocationQueryRaw, RouteLocationRaw } from 'vue-router';
import {
  IconApps,
  IconDashboard,
  IconNotification,
  IconOrderedList,
  IconSearch,
  IconUser,
  IconUserGroup,
} from '@arco-design/web-vue/es/icon/index';
import type { MenuNode, MenuQuery, MenuQueryValue, MenuRoute, MenuSourceNode } from '../types/menu';

const menuIconRegistry: Record<string, Component> = {
  components: IconApps,
  componentsAvatar: IconUser,
  componentsButton: IconApps,
  dashboard: IconDashboard,
  users: IconUserGroup,
  usersList: IconUser,
  usersSegments: IconNotification,
  usersSegmentsRecall: IconSearch,
  orders: IconOrderedList,
  ordersOverview: IconOrderedList,
  ordersReturns: IconNotification,
};

const resolveMenuIcon = (iconKey?: string) => {
  if (!iconKey) {
    return undefined;
  }

  return menuIconRegistry[iconKey];
};

export const hydrateMenuTree = (items: readonly MenuSourceNode[]): MenuNode[] => {
  return items.map((item) => {
    let parsedQuery = item.query;
    if (typeof parsedQuery === 'string') {
      try {
        parsedQuery = JSON.parse(parsedQuery);
      } catch (e) {
        console.warn('Failed to parse menu query JSON string', item.query);
      }
    }

    return {
      ...item,
      query: parsedQuery as any,
      meta: {
        ...item.meta,
        iconComponent: resolveMenuIcon(item.meta?.icon),
      },
      children: item.children?.length ? hydrateMenuTree(item.children) : undefined,
    };
  });
};

export const findMenuChain = (items: readonly MenuNode[], targetKey: string): MenuNode[] => {
  for (const item of items) {
    if (item.key === targetKey) {
      return [item];
    }

    if (item.children?.length) {
      const childChain = findMenuChain(item.children, targetKey);
      if (childChain.length) {
        return [item, ...childChain];
      }
    }
  }

  return [];
};

/** Find menu chain by current route path, using activePath as fallback */
export const findMenuChainByPath = (items: readonly MenuNode[], currentPath: string): MenuNode[] => {
  for (const item of items) {
    // Match by activePath first, then by path
    const matchPath = item.activePath ?? item.path;
    if (matchPath && (matchPath === currentPath || currentPath.startsWith(matchPath + '/'))) {
      return [item];
    }

    if (item.children?.length) {
      const childChain = findMenuChainByPath(item.children, currentPath);
      if (childChain.length) {
        return [item, ...childChain];
      }
    }
  }

  return [];
};

const serializeQueryValue = (value: MenuQueryValue): string | null | undefined => {
  if (value === null || value === undefined) {
    return value;
  }

  if (typeof value === 'string') {
    return value;
  }

  if (typeof value === 'number' || typeof value === 'boolean') {
    return String(value);
  }

  return JSON.stringify(value);
};

const serializeQuery = (query: MenuQuery): LocationQueryRaw => {
  return Object.fromEntries(
    Object.entries(query).map(([key, value]) => {
      if (Array.isArray(value)) {
        const serializedArray = value.every(
          (item) =>
            item === null ||
            item === undefined ||
            typeof item === 'string' ||
            typeof item === 'number' ||
            typeof item === 'boolean'
        )
          ? value.map((item) => serializeQueryValue(item as MenuQueryValue))
          : JSON.stringify(value);

        return [key, serializedArray];
      }

      return [key, serializeQueryValue(value)];
    })
  );
};

const normalizeMenuRoute = (route: MenuRoute): RouteLocationRaw | null => {
  if (!route.path && !route.name) {
    return null;
  }

  return {
    ...(route.path ? { path: route.path } : {}),
    ...(route.name ? { name: route.name } : {}),
    ...(route.hash ? { hash: route.hash } : {}),
    ...(route.query ? { query: serializeQuery(route.query) } : {}),
  };
};

const createRouteLocation = (node: MenuNode): RouteLocationRaw | null => {
  if (node.route) {
    return normalizeMenuRoute(node.route);
  }

  if (!node.path) {
    return null;
  }

  if (!node.query) {
    return node.path;
  }

  return {
    path: node.path,
    query: serializeQuery(node.query),
  };
};

export const createLeafRouteMap = (items: readonly MenuNode[]) => {
  const entries = new Map<string, RouteLocationRaw>();

  const walk = (nodes: readonly MenuNode[]) => {
    for (const node of nodes) {
      const routeLocation = createRouteLocation(node);
      if (routeLocation) {
        entries.set(node.key, routeLocation);
      }

      if (node.children?.length) {
        walk(node.children);
      }
    }
  };

  walk(items);
  return entries;
};

export type { MenuNode, MenuSourceNode } from '../types/menu';