import type { Component } from 'vue';

type MenuQueryScalar = string | number | boolean | null | undefined;
type MenuQueryValue = MenuQueryScalar | Record<string, unknown> | MenuQueryValue[];

export type MenuQuery = Record<string, MenuQueryValue>;

export interface MenuRoute {
  path?: string;
  name?: string;
  query?: MenuQuery;
  hash?: string;
}

export interface MenuMeta {
  titleKey?: string;
  title?: string;
  permissions?: string[];
  icon?: string;
  hidden?: boolean;
  activeMenu?: string;
}

export interface MenuNodeMeta extends MenuMeta {
  iconComponent?: Component;
}

export interface MenuSourceNode {
  key: string;
  path?: string;
  component?: string;
  query?: MenuQuery | string;
  route?: MenuRoute;
  meta?: MenuMeta;
  activePath?: string;  // 激活菜单完整路由，用于路径匹配回退
  children?: MenuSourceNode[];
}

export interface MenuNode extends Omit<MenuSourceNode, 'children' | 'meta' | 'query'> {
  meta?: MenuNodeMeta;
  children?: MenuNode[];
  query?: MenuQuery;
}

export type { MenuQueryValue };