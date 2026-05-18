<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../store/modules/auth';
import { useThemeStore } from '../store/modules/theme';
import { useAppLocale, type AppLocale } from '../locale';
import { IconTranslate } from '@arco-design/web-vue/es/icon';
import { useI18n } from 'vue-i18n';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const themeStore = useThemeStore();
const { appLocale, localeOptions, setAppLocale } = useAppLocale();
const { t } = useI18n();

// 强制浅色：进入 public 页面时锁定 light，离开时恢复用户主题
onMounted(() => {
  document.body.removeAttribute('arco-theme');
  document.documentElement.dataset.theme = 'light';
  document.documentElement.style.colorScheme = 'light';
});

onUnmounted(() => {
  // 恢复用户实际选择的主题
  themeStore.applyTheme();
});

const navItems = computed(() => [
  { path: '/', label: t('portal.home') },
  { path: '/products', label: t('portal.products') },
  { path: '/about', label: t('portal.about') },
  { path: '/help', label: t('portal.help') },
]);

const handleActionBtn = () => {
  const targetPath = authStore.isAuthenticated ? '/dashboard' : '/login';
  const routeUrl = router.resolve({ path: targetPath });
  window.open(routeUrl.href, '_blank');
};

const onLocaleChange = (val: string | number | Record<string, any> | undefined) => {
  if (typeof val === 'string') {
    setAppLocale(val as AppLocale);
  }
};
</script>

<template>
  <div class="public-layout">
    <!-- 毛玻璃特效顶栏 -->
    <header class="public-header">
      <div class="header-container">
        <!-- Logo -->
        <div class="brand" @click="router.push('/')">
          <div class="logo">A</div>
          <span class="brand-name">Aevons</span>
        </div>

        <!-- 中部导航 -->
        <nav class="public-nav">
          <a
            v-for="item in navItems"
            :key="item.path"
            :class="['nav-link', { 'is-active': route.path === item.path }]"
            @click.prevent="router.push(item.path)"
          >
            {{ item.label }}
          </a>
        </nav>

        <!-- 右侧动作按钮区 -->
        <div class="header-actions">
          <a-dropdown @select="onLocaleChange">
            <a-button shape="circle" type="text" class="translate-btn">
              <IconTranslate />
            </a-button>
            <template #content>
              <a-doption 
                v-for="item in localeOptions" 
                :key="item.value" 
                :value="item.value"
                :style="{ color: appLocale === item.value ? 'rgb(var(--primary-6))' : '' }"
              >
                {{ item.label }}
              </a-doption>
            </template>
          </a-dropdown>

          <a-button 
            type="primary" 
            shape="round" 
            class="action-btn"
            @click="handleActionBtn"
          >
            {{ authStore.isAuthenticated ? (appLocale === 'en-US' ? 'Console' : '控制台') : (appLocale === 'en-US' ? 'Login' : '登录') }}
          </a-button>
        </div>
      </div>
    </header>

    <!-- 动态路由承载区 -->
    <main class="public-main">
      <router-view v-slot="{ Component }">
        <transition name="fade-slide" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- 通用页脚 -->
    <footer class="public-footer">
      <div class="footer-content">
        <div class="footer-info">
          <h2>Aevons</h2>
          <p>{{ t('portal.footer.desc') }}</p>
        </div>
        <div class="footer-links">
          <div class="link-group">
            <h3>{{ t('portal.footer.support') }}</h3>
            <a href="#">{{ t('portal.footer.docs') }}</a>
            <a href="#">{{ t('portal.footer.api') }}</a>
            <a href="#">{{ t('portal.footer.changelog') }}</a>
          </div>
          <div class="link-group">
            <h3>{{ t('portal.footer.discover') }}</h3>
            <!-- <a @click.prevent="router.push('/about')">{{ t('portal.footer.aboutUs') }}</a> -->
            <a href="#">{{ t('portal.footer.contact') }}</a>
            <!-- <a href="#">{{ t('portal.footer.join') }}</a> -->
          </div>
        </div>
      </div>
      <div class="footer-bottom">
        2026 Aevons All rights reserved.  | {{ t('portal.footer.copyright') }}
        <br />
        <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer" class="icp-link">黔ICP备2026005938号-1</a>
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* 整个容器铺满并提供平滑背景色 */
.public-layout {
  min-height: 100vh;
  height: 100vh; /* 必须显式声明高度，撑开全局 body 的隐藏截断 */
  display: flex;
  flex-direction: column;
  background-color: var(--color-bg-1);
  color: var(--color-text-1);
  /* 使用现代英文字体提升品牌感，中后置使用系统默认字符防形变 */
  font-family: 'Inter', system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 吸顶导航栏 (Glassmorphism) */
.public-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--color-border-1);
  z-index: 1000;
  transition: all 0.3s ease;
}

.header-container {
  max-width: 1200px;
  height: 100%;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
}

/* 品牌 Logo */
.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}
.logo {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, rgb(var(--primary-6)), rgb(var(--primary-4)));
  color: #ffffff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 16px;
  box-shadow: 0 2px 8px rgba(var(--primary-6), 0.25);
}
.brand-name {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.3px;
}

/* 居中导航 */
.public-nav {
  display: flex;
  gap: 36px;
}
.nav-link {
  font-size: 15px;
  font-weight: 500;
  color: var(--color-text-2);
  cursor: pointer;
  position: relative;
  text-decoration: none;
  transition: color 0.2s ease;
  padding: 8px 0;
}
.nav-link:hover {
  color: rgb(var(--primary-6));
}
.nav-link.is-active {
  color: rgb(var(--primary-6));
}
.nav-link.is-active::after {
  content: '';
  position: absolute;
  bottom: 0px;
  left: 50%;
  transform: translateX(-50%);
  width: 18px;
  height: 3px;
  background-color: rgb(var(--primary-6));
  border-radius: 2px;
}

/* 右侧按钮区 */
.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.translate-btn {
  font-size: 18px;
  color: var(--color-text-2);
  transition: all 0.2s ease;
}
.translate-btn:hover {
  color: rgb(var(--primary-6));
  background-color: var(--color-fill-2);
}

.action-btn {
  font-weight: 600;
  padding: 0 24px;
  transition: transform 0.2s ease, box-shadow 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.action-btn:hover {
  transform: translateY(-2px) scale(1.02);
  box-shadow: 0 6px 16px rgba(var(--primary-6), 0.35);
}

/* --- 主内容区 --- */
.public-main {
  flex: 1;
  padding-top: 64px; /* 为 Fixed Header 留出空间避免遮挡 */
  display: flex;
  flex-direction: column;
}

/* 路由全局转场动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(12px) scale(0.99);
}
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-12px) scale(0.99);
}

/* --- 极简商业风页脚 --- */
.public-footer {
  background-color: var(--color-bg-2);
  border-top: 1px solid var(--color-border-1);
  padding: 80px 24px 32px;
  margin-top: auto;
}
.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 48px;
  margin-bottom: 64px;
}
.footer-info h2 {
  font-size: 26px;
  font-weight: 800;
  margin: 0 0 16px 0;
  background: linear-gradient(to right, rgb(var(--primary-6)), rgb(var(--primary-4)));
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}
.footer-info p {
  color: var(--color-text-3);
  max-width: 320px;
  line-height: 1.7;
  font-size: 14px;
}
.footer-links {
  display: flex;
  gap: 80px;
}
.link-group h3 {
  font-size: 16px;
  color: var(--color-text-1);
  margin: 0 0 24px 0;
  font-weight: 600;
}
.link-group a {
  display: block;
  color: var(--color-text-3);
  margin-bottom: 14px;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  transition: color 0.2s ease, transform 0.2s ease;
}
.link-group a:hover {
  color: rgb(var(--primary-6));
  transform: translateX(3px);
}
.footer-bottom {
  max-width: 1200px;
  margin: 0 auto;
  text-align: center;
  padding-top: 32px;
  border-top: 1px solid var(--color-border-1);
  color: var(--color-text-4);
  font-size: 13px;
  letter-spacing: 0.5px;
}
.icp-link {
  color: var(--color-text-4);
  text-decoration: none;
  transition: color 0.2s ease;
}
.icp-link:hover {
  color: rgb(var(--primary-6));
}
</style>
