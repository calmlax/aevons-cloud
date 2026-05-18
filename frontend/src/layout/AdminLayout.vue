<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRoute, useRouter, type RouteLocationRaw } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { Message, Modal } from '@arco-design/web-vue';
import {
  IconCheck,
  IconDesktop,
  IconDown,
  IconExport,
  IconFullscreen,
  IconFullscreenExit,
  IconLock,
  IconMenuFold,
  IconMenuUnfold,
  IconMoonFill,
  IconNotification,
  IconSearch,
  IconSunFill,
  IconTranslate,
  IconUser,
} from '@arco-design/web-vue/es/icon/index';
import { storeToRefs } from 'pinia';

import SideMenuNode from '../components/SideMenuNode.js';
import { useAppLocale, type AppLocale } from '../locale';
import { useAuthStore } from '../store/modules/auth';
import { useMenuStore } from '../store/modules/menu';
import { useThemeStore, type ThemeMode } from '../store/modules/theme';
import { useLayoutStore, type LayoutMode } from '../store/modules/layout';
import { useNotificationStore } from '../store/modules/notification';
import { createLeafRouteMap, findMenuChain, type MenuNode } from '../utils/menu';

// ── 产品服务菜单数据 ──────────────────────────────────
interface ProductItem {
  name: string;
  desc: string;
  icon: string;
  tag?: string;
}
interface ProductCategory {
  category: string;
  items: ProductItem[];
}

const productCategories: ProductCategory[] = [
  {
    category: '计算',
    items: [
      { name: '云服务器 CVM', desc: '弹性、安全、稳定的云端计算服务', icon: '🖥️' },
      { name: '容器服务 TKE', desc: '高度可扩展的 Kubernetes 容器集群', icon: '📦', tag: '热门' },
      { name: '弹性伸缩 AS', desc: '根据业务负载自动调整计算资源', icon: '⚡' },
      { name: '批量计算', desc: '大规模并行批处理计算服务', icon: '🔢' },
    ],
  },
  {
    category: '存储',
    items: [
      { name: '对象存储 COS', desc: '安全稳定、海量、便捷的云端存储', icon: '🗄️', tag: '热门' },
      { name: '文件存储 CFS', desc: '可扩展的共享文件存储服务', icon: '📁' },
      { name: '归档存储', desc: '低成本的长期数据归档存储', icon: '🗃️' },
      { name: '云硬盘 CBS', desc: '高可用、高可靠的块存储设备', icon: '💾' },
    ],
  },
  {
    category: '数据库',
    items: [
      { name: '云数据库 MySQL', desc: '高性能、高可靠的 MySQL 云数据库', icon: '🐬', tag: '热门' },
      { name: '云数据库 Redis', desc: '兼容 Redis 协议的内存数据库', icon: '🔴' },
      { name: 'TDSQL-C', desc: '新一代云原生关系型数据库', icon: '🌊', tag: 'New' },
      { name: '时序数据库 CTSDB', desc: '高效存储和查询时序数据', icon: '📈' },
    ],
  },
  {
    category: '网络',
    items: [
      { name: '私有网络 VPC', desc: '自定义逻辑隔离的网络空间', icon: '🌐' },
      { name: '负载均衡 CLB', desc: '将流量分发到多台后端服务器', icon: '⚖️', tag: '热门' },
      { name: 'CDN 内容分发', desc: '全球加速，就近分发内容', icon: '🚀' },
      { name: '专线接入 DC', desc: '高速稳定的专属网络连接', icon: '🔗' },
    ],
  },
  {
    category: '安全',
    items: [
      { name: 'Web 应用防火墙', desc: '防御 Web 攻击，保护网站安全', icon: '🛡️', tag: '热门' },
      { name: 'DDoS 防护', desc: '多维度 DDoS 攻击防护方案', icon: '🔒' },
      { name: '密钥管理 KMS', desc: '安全合规的密钥管理服务', icon: '🔑' },
      { name: '安全运营中心', desc: '统一安全事件管理与响应', icon: '🎯' },
    ],
  },
  {
    category: '人工智能',
    items: [
      { name: '自然语言处理', desc: '文本分析、情感识别、翻译', icon: '🤖', tag: 'AI' },
      { name: '图像识别', desc: '图片内容智能分析与识别', icon: '👁️', tag: 'AI' },
      { name: '语音识别 ASR', desc: '将语音转化为文字的云服务', icon: '🎙️' },
      { name: '大模型服务', desc: '企业级大语言模型 API 接入', icon: '✨', tag: 'New' },
    ],
  },
];
// ─────────────────────────────────────────────────────

const route = useRoute();
const router = useRouter();
const { t } = useI18n();

const collapsed = ref(false);
const mobileMenuOpen = ref(false);
const isMobile = ref(false);
const isTablet = ref(false);
const accountMenuVisible = ref(false);
const localeMenuVisible = ref(false);
const searchMenuVisible = ref(false);
const themeMenuVisible = ref(false);
const productMenuVisible = ref(false);
const hiddenBreadcrumbsVisible = ref(false);

const getTagClass = (tag: string) => {
  const map: Record<string, string> = { New: 'product-item-tag--new', AI: 'product-item-tag--ai' };
  return ['product-item-tag', map[tag] ?? 'product-item-tag--hot'];
};
const isFullscreen = ref(false);
const headerElevated = ref(false);
const menuSearchKeyword = ref('');
const openKeys = ref<string[]>([]);

const themeStore = useThemeStore();
const { resolvedTheme, themeMode } = storeToRefs(themeStore);
const { appLocale, localeOptions, setAppLocale } = useAppLocale();
const authStore = useAuthStore();
const menuStore = useMenuStore();
const { currentUser } = storeToRefs(authStore);
const { logout } = authStore;
const { menuError, menuTree } = storeToRefs(menuStore);
const { initializeMenu } = menuStore;

// ── Layout mode ───────────────────────────────────────
const layoutStore = useLayoutStore();
const { layoutMode } = storeToRefs(layoutStore);

// ── Notification store ────────────────────────────────
const notifStore = useNotificationStore();
const { unreadCount: notifUnreadCount } = storeToRefs(notifStore);

// topnav 模式：当前激活的一级菜单 key
const activeTopKey = ref<string>('');
// topnav 模式：侧边栏显示的子菜单树（依赖 visibleMenuTree，在其后定义）
const topnavSubTree = computed<MenuNode[]>(() => {
  if (layoutMode.value !== 'topnav' || isMobile.value) return [];
  const node = visibleMenuTree.value.find((n) => n.key === activeTopKey.value);
  return node?.children ?? [];
});

const onTopNavItemClick = (key: string) => {
  activeTopKey.value = key;
  const node = visibleMenuTree.value.find((n) => n.key === key);
  if (!node?.children?.length) {
    const target = menuRouteMap.value.get(key);
    if (target) router.push(target);
  }
};

const layoutOptions: { label: string; value: LayoutMode }[] = [
  { label: '侧边栏模式', value: 'sidebar' },
  { label: '顶栏模式', value: 'topnav' },
];
const layoutMenuVisible = ref(false);

const onLayoutChange = (mode: LayoutMode) => {
  layoutStore.setLayoutMode(mode);
  layoutMenuVisible.value = false;
};

const visibleMenuTree = computed(() => {
  const filterVisible = (nodes: MenuNode[]): MenuNode[] =>
    nodes
      .filter((node) => node.meta?.hidden !== true)
      .map((node) => {
        if (!node.children?.length) {
          return node;
        }

        const visibleChildren = filterVisible(node.children);
        return {
          ...node,
          children: visibleChildren.length ? visibleChildren : undefined,
        };
      });

  return filterVisible(menuTree.value);
});

const menuRouteMap = computed(() => createLeafRouteMap(menuTree.value));
const themeOptions = computed(() => [
  { label: t('theme.light'), value: 'light' as ThemeMode },
  { label: t('theme.dark'), value: 'dark' as ThemeMode },
  { label: t('theme.system'), value: 'system' as ThemeMode },
]);
const themeTriggerIcon = computed(() => {
  if (themeMode.value === 'system') {
    return IconDesktop;
  }

  return resolvedTheme.value === 'dark' ? IconSunFill : IconMoonFill;
});
const fullscreenTriggerIcon = computed(() => (isFullscreen.value ? IconFullscreenExit : IconFullscreen));
const menuTheme = computed(() => (resolvedTheme.value === 'dark' ? 'dark' : 'light'));
const accountProfile = computed(() => {
  const user = currentUser.value;
  const displayName = user?.nickname || user?.username || 'A';
  const fallbackRole = user?.roles?.[0]?.role_name || user?.roles?.[0]?.role_key;

  return {
    name: displayName,
    roleKey: fallbackRole || 'Admin',
    initials: displayName.trim().charAt(0).toUpperCase() || 'A',
    avatar: user?.avatar,
  };
});
const accountActions = computed(() => [
  { key: 'profile' as const, labelKey: 'account.profile', icon: IconUser },
  { key: 'password' as const, labelKey: 'account.password', icon: IconLock },
  { key: 'logout' as const, labelKey: 'account.logout', icon: IconExport, danger: true },
]);

const selectedKey = computed(() => {
  if (route.meta.activeMenu) {
    return String(route.meta.activeMenu);
  }

  return route.meta.menuKey ? String(route.meta.menuKey) : '';
});
const exactMenuChain = computed(() => findMenuChain(menuTree.value, route.meta.menuKey ? String(route.meta.menuKey) : ''));
const routeBreadcrumbs = computed(() => {
  if (Array.isArray(route.meta.breadcrumbKeys)) {
    return route.meta.breadcrumbKeys.map((item) => ({ title: t(String(item)), path: '' }));
  }

  return exactMenuChain.value.map((item) => ({
    title: item.meta?.title || t(item.meta?.titleKey ?? 'layout.defaultTitle'),
    path: item.path,
  }));
});
const currentTitle = computed(() => {
  const fallbackTitleKey = exactMenuChain.value.at(-1)?.meta?.titleKey;

  if (typeof route.meta.titleKey === 'string') {
    return t(route.meta.titleKey);
  }

  return fallbackTitleKey ? t(fallbackTitleKey) : exactMenuChain.value.at(-1)?.meta?.title || t('layout.defaultTitle');
});
const breadcrumbs = computed(() => [{ title: t('layout.productName'), path: '/dashboard' }, ...routeBreadcrumbs.value]);
const displayedBreadcrumbs = computed(() => {
  if (breadcrumbs.value.length <= 3) {
    return breadcrumbs.value.map((item) => ({ ...item, ellipsis: false, hiddenItems: [] as typeof breadcrumbs.value }));
  }

  return [
    { ...breadcrumbs.value[0], ellipsis: false, hiddenItems: [] as typeof breadcrumbs.value },
    {
      title: '...',
      path: '',
      ellipsis: true,
      hiddenItems: breadcrumbs.value.slice(1, -1),
    },
    { ...breadcrumbs.value.at(-1)!, ellipsis: false, hiddenItems: [] as typeof breadcrumbs.value },
  ];
});
const searchableMenuItems = computed(() => {
  const entries: Array<{ key: string; label: string; trail: string; route: RouteLocationRaw }> = [];

  const walk = (nodes: MenuNode[], parents: string[] = []) => {
    for (const node of nodes) {
      const currentLabel = node.meta?.title || t(node.meta?.titleKey ?? 'layout.defaultTitle');
      const nextParents = [...parents, currentLabel];

      if (node.children?.length) {
        walk(node.children, nextParents);
        continue;
      }

      const routeTarget = menuRouteMap.value.get(node.key);
      if (routeTarget) {
        entries.push({
          key: node.key,
          label: currentLabel,
          trail: nextParents.join(' / '),
          route: routeTarget,
        });
      }
    }
  };

  walk(visibleMenuTree.value);
  return entries;
});
const filteredMenuItems = computed(() => {
  const keyword = menuSearchKeyword.value.trim().toLowerCase();

  if (!keyword) {
    return searchableMenuItems.value;
  }

  return searchableMenuItems.value.filter((item) => {
    return item.label.toLowerCase().includes(keyword) || item.trail.toLowerCase().includes(keyword);
  });
});

watch(
  exactMenuChain,
  (chain) => {
    // 确保当前路由的父节点也在展开列表里（不覆盖已有展开状态）
    const parentKeys = chain.slice(0, -1).map((item) => item.key);
    const merged = Array.from(new Set([...openKeys.value, ...parentKeys]));
    openKeys.value = merged;
  },
  { immediate: true }
);

// topnav 模式：根据当前路由自动激活一级菜单
watch(
  [exactMenuChain, layoutMode],
  ([chain]) => {
    if (layoutMode.value === 'topnav' && chain.length > 0) {
      activeTopKey.value = chain[0].key;
    }
  },
  { immediate: true },
);

// 收集菜单树中所有有子菜单的节点 key，用于默认全部展开
const collectParentKeys = (nodes: MenuNode[]): string[] => {
  const keys: string[] = [];
  for (const node of nodes) {
    if (node.children?.length) {
      keys.push(node.key);
      keys.push(...collectParentKeys(node.children));
    }
  }
  return keys;
};

// 菜单树加载后默认展开所有父节点
watch(
  visibleMenuTree,
  (tree) => {
    if (tree.length > 0) {
      const allParentKeys = collectParentKeys(tree);
      openKeys.value = Array.from(new Set([...openKeys.value, ...allParentKeys]));
    }
  },
  { immediate: true },
);

watch(
  () => route.fullPath,
  async () => {
    await nextTick();
    const mainElement = document.querySelector('.admin-main') as HTMLElement | null;
    if (mainElement) {
      mainElement.scrollTop = 0;
    }

    searchMenuVisible.value = false;
    hiddenBreadcrumbsVisible.value = false;
    menuSearchKeyword.value = '';
    headerElevated.value = false;

    if (isMobile.value) {
      mobileMenuOpen.value = false;
    }
  }
);

const syncViewport = () => {
  const mobile = window.matchMedia('(max-width: 767px)').matches;
  const tablet = window.matchMedia('(min-width: 768px) and (max-width: 1024px)').matches;

  isMobile.value = mobile;
  isTablet.value = tablet;

  if (mobile) {
    collapsed.value = false;
    return;
  }

  mobileMenuOpen.value = false;
  collapsed.value = tablet;
};

const onMenuClick = (key: string) => {
  const target = menuRouteMap.value.get(key);
  if (target) {
    router.push(target);
  }

  if (isMobile.value) {
    mobileMenuOpen.value = false;
  }
};

const onOpenKeysChange = (keys: string[]) => {
  openKeys.value = keys;
};

const onMainScroll = (event: Event) => {
  const target = event.target as HTMLElement | null;
  headerElevated.value = Boolean(target && target.scrollTop > 4);
};

let touchStartX = 0;

const onSiderTouchStart = (e: TouchEvent) => {
  touchStartX = e.touches[0].clientX;
};

const onSiderTouchEnd = (e: TouchEvent) => {
  const dx = e.changedTouches[0].clientX - touchStartX;
  if (dx < -50) {
    mobileMenuOpen.value = false;
  }
};

const onMenuSearchNavigate = (key: string) => {
  const target = menuRouteMap.value.get(key);
  if (!target) {
    return;
  }

  searchMenuVisible.value = false;
  menuSearchKeyword.value = '';
  router.push(target);
};

const onToggleSidebar = () => {
  if (isMobile.value) {
    mobileMenuOpen.value = !mobileMenuOpen.value;
    return;
  }

  collapsed.value = !collapsed.value;
};

const onBreadcrumbNavigate = (path?: string) => {
  if (!path) {
    return;
  }

  hiddenBreadcrumbsVisible.value = false;
  router.push(path);
};

const syncFullscreenState = () => {
  isFullscreen.value = Boolean(document.fullscreenElement);
};

const onFullscreenChange = () => {
  syncFullscreenState();
};

const onThemeChange = (value: string | number | Record<string, unknown> | undefined) => {
  if (value === 'light' || value === 'dark' || value === 'system') {
    themeStore.setThemeMode(value as ThemeMode);
    themeMenuVisible.value = false;
  }
};

const onLocaleChange = (value: string | number | Record<string, unknown> | undefined) => {
  if (value === 'zh-CN' || value === 'en-US') {
    setAppLocale(value as AppLocale);
    localeMenuVisible.value = false;
    window.location.reload();
  }
};

const onToggleFullscreen = async () => {
  if (!document.fullscreenEnabled) {
    Message.warning(t('layout.fullscreenUnsupported'));
    return;
  }

  if (document.fullscreenElement) {
    await document.exitFullscreen();
    return;
  }

  await document.documentElement.requestFullscreen();
};

const onAccountAction = (action: 'profile' | 'password' | 'logout') => {
  accountMenuVisible.value = false;

  if (action === 'profile') {
    router.push('/account/profile');
    return;
  }

  if (action === 'password') {
    router.push('/account/password');
    return;
  }

  Modal.confirm({
    title: t('account.logoutTitle'),
    content: t('account.logoutContent'),
    okText: t('account.logoutOk'),
    cancelText: t('account.logoutCancel'),
    onOk: async () => {
      await logout();
      Message.success(t('account.logoutSuccess'));
      router.push('/login');
    },
  });
};

onMounted(() => {
  syncFullscreenState();
  syncViewport();

  if (menuTree.value.length === 0) {
    void initializeMenu().catch(() => undefined);
  }

  // 初始化通知 mock 数据（仅在 store 为空时）
  if (notifStore.notifications.length === 0) {
    notifStore.setNotifications([
      { id: 1, level: 'error', title: t('notifications.items.shipmentAlert'), body: t('notifications.items.shipmentAlertBody'), time: '10:32', read: false },
      { id: 2, level: 'warning', title: t('notifications.items.roiWarning'), body: t('notifications.items.roiWarningBody'), time: '09:15', read: false },
      { id: 3, level: 'success', title: t('notifications.items.campaignDone'), body: t('notifications.items.campaignDoneBody'), time: '昨天 18:00', read: false },
      { id: 4, level: 'info', title: t('notifications.items.systemUpdate'), body: t('notifications.items.systemUpdateBody'), time: '昨天 12:00', read: true },
      { id: 5, level: 'success', title: t('notifications.items.exportReady'), body: t('notifications.items.exportReadyBody'), time: '2天前', read: true },
      { id: 6, level: 'info', title: t('notifications.items.newMember'), body: t('notifications.items.newMemberBody'), time: '3天前', read: true },
    ]);
  }

  document.addEventListener('fullscreenchange', onFullscreenChange);
  window.addEventListener('resize', syncViewport);
  document.addEventListener('click', (e) => {
    const target = e.target as HTMLElement;
    if (!target.closest('.top-nav-product-btn') && !target.closest('.product-panel')) {
      productMenuVisible.value = false;
    }
  });
});

onUnmounted(() => {
  document.removeEventListener('fullscreenchange', onFullscreenChange);
  window.removeEventListener('resize', syncViewport);
});
</script>

<template>
  <div class="admin-shell tcloud-shell">
    <header :class="['top-nav', { 'top-nav--elevated': headerElevated }]">
      <div class="top-nav-brand">
        <button  v-if="isMobile"
          type="button"
          class="top-nav-collapse"
          :aria-label="collapsed || mobileMenuOpen ? t('layout.openMenu') : t('layout.closeMenu')"
          @click="onToggleSidebar"
        >
          <component :is="collapsed || mobileMenuOpen ? IconMenuUnfold : IconMenuFold" />
        </button>

        <div class="top-nav-brand-copy">
          <strong>{{ t('layout.productName') }}</strong>
          <span v-if="!isMobile">{{ t('layout.console') }}</span>
        </div>
      </div>

      <nav v-if="!isMobile" class="top-nav-main">
        <!-- topnav 模式：一级菜单展示在顶栏 -->
        <template v-if="layoutMode === 'topnav' && !isMobile">
          <button
            v-for="node in visibleMenuTree"
            :key="node.key"
            type="button"
            :class="['top-nav-main-item', 'top-nav-menu-btn', { 'is-active': activeTopKey === node.key }]"
            @click="onTopNavItemClick(node.key)"
          >
            {{ node.meta?.title || t(node.meta?.titleKey ?? '') }}
          </button>
        </template>

        <!-- 产品服务入口 -->
        <!-- <button
          type="button"
          :class="['top-nav-main-item', 'top-nav-product-btn', { 'is-open': productMenuVisible }]"
          @click="productMenuVisible = !productMenuVisible"
        >
          产品服务
          <svg class="top-nav-product-arrow" width="12" height="12" viewBox="0 0 12 12" fill="none">
            <path d="M2 4l4 4 4-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button> -->
      </nav>

      <!-- 产品服务下拉面板（全宽） -->
      <transition name="product-panel">
        <div v-if="productMenuVisible && !isMobile" class="product-panel" @click.self="productMenuVisible = false">
          <div class="product-panel-inner">
            <div class="product-panel-header">
              <span class="product-panel-title">产品服务</span>
              <button type="button" class="product-panel-close" @click="productMenuVisible = false">✕</button>
            </div>
            <div class="product-panel-grid">
              <div v-for="cat in productCategories" :key="cat.category" class="product-cat">
                <div class="product-cat-label">{{ cat.category }}</div>
                <div class="product-cat-items">
                  <button
                    v-for="item in cat.items"
                    :key="item.name"
                    type="button"
                    class="product-item"
                    @click="productMenuVisible = false"
                  >
                    <span class="product-item-icon">{{ item.icon }}</span>
                    <span class="product-item-body">
                      <span class="product-item-name">
                        {{ item.name }}
                        <span v-if="item.tag" :class="getTagClass(item.tag)">{{ item.tag }}</span>
                      </span>
                      <span class="product-item-desc">{{ item.desc }}</span>
                    </span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </transition>

      <div v-if="!isMobile" class="top-nav-tools">
        <!-- 布局模式切换 -->
        <a-popover
          v-model:popup-visible="layoutMenuVisible"
          trigger="click"
          position="bl"
          popup-container="body"
          content-class="toolbar-popover"
        >
          <button type="button" class="header-tool-button" :aria-label="'布局模式'">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <rect x="1" y="1" width="5" height="14" rx="1" stroke="currentColor" stroke-width="1.4"/>
              <rect x="8" y="1" width="7" height="6" rx="1" stroke="currentColor" stroke-width="1.4"/>
              <rect x="8" y="9" width="7" height="6" rx="1" stroke="currentColor" stroke-width="1.4"/>
            </svg>
          </button>
          <template #content>
            <div class="toolbar-menu">
              <button
                v-for="item in layoutOptions"
                :key="item.value"
                type="button"
                :class="['toolbar-menu-button', { 'is-active': layoutMode === item.value }]"
                @click="onLayoutChange(item.value)"
              >
                <span>{{ item.label }}</span>
                <IconCheck v-if="layoutMode === item.value" class="toolbar-menu-check" />
              </button>
            </div>
          </template>
        </a-popover>
        <a-popover
          v-model:popup-visible="searchMenuVisible"
          trigger="click"
          position="bl"
          popup-container="body"
          content-class="toolbar-popover menu-search-popover"
        >
          <button type="button" class="header-tool-button" :aria-label="t('layout.menuSearch')">
            <IconSearch />
          </button>

          <template #content>
            <div class="menu-search-panel">
              <a-input-search
                v-model="menuSearchKeyword"
                class="menu-search-input"
                allow-clear
                :placeholder="t('layout.menuSearchPlaceholder')"
              />

              <div v-if="filteredMenuItems.length" class="menu-search-results">
                <button
                  v-for="item in filteredMenuItems"
                  :key="item.key"
                  type="button"
                  class="menu-search-item"
                  @click="onMenuSearchNavigate(item.key)"
                >
                  <div class="menu-search-item-copy">
                    <strong>{{ item.label }}</strong>
                    <span>{{ item.trail }}</span>
                  </div>
                </button>
              </div>

              <div v-else class="menu-search-empty">
                {{ t('layout.menuSearchEmpty') }}
              </div>
            </div>
          </template>
        </a-popover>

        <a-popover
          v-model:popup-visible="localeMenuVisible"
          trigger="click"
          position="bl"
          popup-container="body"
          content-class="toolbar-popover"
        >
          <button type="button" class="header-tool-button" :aria-label="t('locale.label')">
            <IconTranslate />
          </button>

          <template #content>
            <div class="toolbar-menu">
              <button
                v-for="item in localeOptions"
                :key="item.value"
                type="button"
                :class="['toolbar-menu-button', { 'is-active': appLocale === item.value }]"
                @click="onLocaleChange(item.value)"
              >
                <span>{{ item.label }}</span>
                <IconCheck v-if="appLocale === item.value" class="toolbar-menu-check" />
              </button>
            </div>
          </template>
        </a-popover>

        <a-popover
          v-model:popup-visible="themeMenuVisible"
          trigger="click"
          position="bl"
          popup-container="body"
          content-class="toolbar-popover"
        >
          <button type="button" class="header-tool-button" :aria-label="t('theme.label')">
            <component :is="themeTriggerIcon" />
          </button>

          <template #content>
            <div class="toolbar-menu">
              <button
                v-for="item in themeOptions"
                :key="item.value"
                type="button"
                :class="['toolbar-menu-button', { 'is-active': themeMode === item.value }]"
                @click="onThemeChange(item.value)"
              >
                <span>{{ item.label }}</span>
                <IconCheck v-if="themeMode === item.value" class="toolbar-menu-check" />
              </button>
            </div>
          </template>
        </a-popover>

        <button
          type="button"
          class="header-tool-button"
          :aria-label="isFullscreen ? t('layout.exitFullscreen') : t('layout.fullscreen')"
          @click="onToggleFullscreen"
        >
          <component :is="fullscreenTriggerIcon" />
        </button>

        <button type="button" class="header-tool-button" :aria-label="t('notifications.title')" @click="router.push('/notifications')">
          <a-badge :count="notifUnreadCount" :max-count="99" dot>
            <IconNotification />
          </a-badge>
        </button>
      </div>

      <a-popover
        v-model:popup-visible="accountMenuVisible"
        trigger="click"
        position="br"
        popup-container="body"
        content-class="account-popover"
      >
        <button type="button" class="account-trigger">
          <a-avatar :size="32" class="header-avatar">
            <img v-if="accountProfile.avatar" :src="accountProfile.avatar" alt="avatar" />
            <template v-else>{{ accountProfile.initials }}</template>
          </a-avatar>

          <div v-if="!isMobile" class="account-copy">
            <strong>{{ accountProfile.name }}</strong>
            <span>{{ accountProfile.roleKey }}</span>
          </div>

          <IconDown v-if="!isMobile" class="account-trigger-icon" />
        </button>

        <template #content>
          <div class="account-menu">
            <div class="account-menu-header">
              <a-avatar :size="42" class="header-avatar header-avatar--panel">
                <img v-if="accountProfile.avatar" :src="accountProfile.avatar" alt="avatar" />
                <template v-else>{{ accountProfile.initials }}</template>
              </a-avatar>
              <div>
                <strong>{{ accountProfile.name }}</strong>
                <span>{{ accountProfile.roleKey }}</span>
              </div>
            </div>

            <button
              v-if="isMobile"
              type="button"
              class="account-menu-button"
              :aria-label="t('locale.label')"
              @click="localeMenuVisible = !localeMenuVisible"
            >
              <IconTranslate class="account-menu-button-icon" />
              {{ t('locale.label') }}
            </button>

            <div v-if="isMobile && localeMenuVisible" class="toolbar-menu toolbar-menu--inline">
              <button
                v-for="item in localeOptions"
                :key="item.value"
                type="button"
                :class="['toolbar-menu-button', { 'is-active': appLocale === item.value }]"
                @click="onLocaleChange(item.value)"
              >
                <span>{{ item.label }}</span>
                <IconCheck v-if="appLocale === item.value" class="toolbar-menu-check" />
              </button>
            </div>

            <button
              v-if="isMobile"
              type="button"
              class="account-menu-button"
              :aria-label="t('theme.label')"
              @click="themeMenuVisible = !themeMenuVisible"
            >
              <component :is="themeTriggerIcon" class="account-menu-button-icon" />
              {{ t('theme.label') }}
            </button>

            <div v-if="isMobile && themeMenuVisible" class="toolbar-menu toolbar-menu--inline">
              <button
                v-for="item in themeOptions"
                :key="item.value"
                type="button"
                :class="['toolbar-menu-button', { 'is-active': themeMode === item.value }]"
                @click="onThemeChange(item.value)"
              >
                <span>{{ item.label }}</span>
                <IconCheck v-if="themeMode === item.value" class="toolbar-menu-check" />
              </button>
            </div>

            <button
              v-for="item in accountActions"
              :key="item.key"
              type="button"
              :class="['account-menu-button', { 'account-menu-button--danger': item.danger }]"
              @click="onAccountAction(item.key)"
            >
              <component :is="item.icon" class="account-menu-button-icon" />
              {{ t(item.labelKey) }}
            </button>
          </div>
        </template>
      </a-popover>
    </header>

    <div class="admin-body">
      <div
        v-if="isMobile"
        :class="['admin-sider-overlay', { 'admin-sider-overlay--visible': mobileMenuOpen }]"
        @click="mobileMenuOpen = false"
      />

      <!-- sidebar 模式：完整菜单树 / topnav 模式：当前一级菜单的子菜单 -->
      <a-layout-sider
        v-if="layoutMode === 'sidebar' || isMobile || (layoutMode === 'topnav' && topnavSubTree.length > 0)"
        :collapsed="isMobile ? false : collapsed"
        :width="220"
        :collapsed-width="64"
        :theme="menuTheme"
        :class="['admin-sider', { 'admin-sider--open': isMobile && mobileMenuOpen }]"
        hide-trigger
        @touchstart.passive="onSiderTouchStart"
        @touchend.passive="onSiderTouchEnd"
      >
        <a-menu
          :selected-keys="selectedKey ? [selectedKey] : []"
          :open-keys="openKeys"
          :collapsed="isMobile ? false : collapsed"
          :theme="menuTheme"
          class="nav-menu"
          @update:open-keys="onOpenKeysChange"
          @menu-item-click="onMenuClick"
        >
          <SideMenuNode
            v-for="item in (layoutMode === 'topnav' && !isMobile ? topnavSubTree : visibleMenuTree)"
            :key="item.key"
            :item="item"
          />
        </a-menu>

        <div v-if="menuError && !menuTree.length" class="nav-menu-state">
          {{ t('layout.menuLoadFailed') }}
        </div>

        <!-- 收起/展开按钮：固定在侧边栏底部右侧 -->
        <div v-if="!isMobile" class="sider-collapse-btn-wrap">
          <button
            type="button"
            class="sider-collapse-btn"
            :aria-label="collapsed ? t('layout.openMenu') : t('layout.closeMenu')"
            @click="onToggleSidebar"
          >
            <component :is="collapsed ? IconMenuUnfold : IconMenuFold" />
          </button>
        </div>
      </a-layout-sider>

      <main class="admin-main" @scroll="onMainScroll">
        <div class="content-topbar">
          <div class="console-breadcrumbs">
            <template v-for="(crumb, index) in breadcrumbs" :key="`${crumb.title}-${index}`">
              <span v-if="index > 0" class="console-breadcrumb-separator">/</span>
              <button
                v-if="crumb.path && index < breadcrumbs.length - 1"
                type="button"
                class="console-breadcrumb-button"
                @click="onBreadcrumbNavigate(crumb.path)"
              >
                {{ crumb.title }}
              </button>
              <span v-else class="console-breadcrumb-current">{{ crumb.title }}</span>
            </template>
          </div>
        </div>

        <div class="page-header">
          <h1>{{ currentTitle }}</h1>
        </div>

        <section class="admin-page-body">
          <router-view v-slot="{ Component }">
            <transition name="page-fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </section>
      </main>
    </div>
  </div>
</template>

<style scoped>
/* ── 侧边栏收起按钮 ── */
.sider-collapse-btn-wrap {
  position: absolute;
  bottom: 12px;
  right: 12px;
  z-index: 10;
}

.sider-collapse-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--sider-border, rgba(0,0,0,0.06));
  border-radius: 6px;
  background: var(--content-bg, #fff);
  color: var(--text-subtle, #4e5969);
  cursor: pointer;
  transition: all 0.15s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}

.sider-collapse-btn:hover {
  background: var(--color-fill-2, rgba(0,0,0,0.04));
  color: var(--text-main, #1d2129);
  border-color: var(--panel-hover-border, rgba(0,82,217,0.22));
}

/* ── 产品服务按钮 ── */
.top-nav-product-btn,
.top-nav-menu-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 14px;
  font-family: inherit;
}

.top-nav-menu-btn.is-active {
  color: var(--sider-item-selected-text, #0052d9);
  font-weight: 600;
}

.top-nav-product-arrow {
  transition: transform 0.2s ease;
  color: currentColor;
}

.top-nav-product-btn.is-open .top-nav-product-arrow {
  transform: rotate(180deg);
}

/* ── 产品面板 ── */
.product-panel {
  position: fixed;
  top: 56px;
  left: 0;
  right: 0;
  z-index: 99;
  background: var(--content-bg, #fff);
  border-bottom: 1px solid var(--topbar-border, rgba(0,0,0,0.06));
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.product-panel-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px 24px 24px;
}

.product-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.product-panel-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-main, #1d2129);
}

.product-panel-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--text-muted, #86909c);
  font-size: 13px;
  cursor: pointer;
  transition: background 0.15s;
}

.product-panel-close:hover {
  background: var(--color-fill-2, rgba(0,0,0,0.04));
  color: var(--text-main, #1d2129);
}

/* 6列网格，每列一个分类 */
.product-panel-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 0 16px;
}

.product-cat-label {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--text-muted, #86909c);
  padding: 0 8px 8px;
  border-bottom: 1px solid var(--divider, rgba(0,0,0,0.05));
  margin-bottom: 4px;
}

.product-cat-items {
  display: flex;
  flex-direction: column;
}

.product-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  width: 100%;
  padding: 7px 8px;
  border: none;
  border-radius: 6px;
  background: transparent;
  text-align: left;
  cursor: pointer;
  transition: background 0.15s;
}

.product-item:hover {
  background: var(--color-fill-2, rgba(0,0,0,0.04));
}

.product-item-icon {
  font-size: 16px;
  line-height: 1;
  flex-shrink: 0;
  margin-top: 1px;
}

.product-item-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.product-item-name {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-main, #1d2129);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.product-item-desc {
  font-size: 11px;
  color: var(--text-muted, #86909c);
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* tag 样式 */
.product-item-tag {
  display: inline-block;
  padding: 0 4px;
  border-radius: 3px;
  font-size: 10px;
  font-weight: 600;
  line-height: 16px;
  flex-shrink: 0;
}

.product-item-tag--hot {
  background: rgba(245, 63, 63, 0.1);
  color: #f53f3f;
}

.product-item-tag--new {
  background: rgba(0, 180, 42, 0.1);
  color: #00b42a;
}

.product-item-tag--ai {
  background: rgba(114, 46, 209, 0.1);
  color: #722ed1;
}

/* 进入/离开动画 */
.product-panel-enter-active,
.product-panel-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.product-panel-enter-from,
.product-panel-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
