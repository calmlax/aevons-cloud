<script setup lang="ts">
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  IconCheck,
  IconDesktop,
  IconMoonFill,
  IconSunFill,
  IconTranslate,
} from '@arco-design/web-vue/es/icon/index';
import { storeToRefs } from 'pinia';

import { useAppLocale, type AppLocale } from '../locale';
import { useThemeStore, type ThemeMode } from '../store/modules/theme';

const { t } = useI18n();
const localeMenuVisible = ref(false);
const themeMenuVisible = ref(false);
const { appLocale, localeOptions, setAppLocale } = useAppLocale();
const themeStore = useThemeStore();
const { resolvedTheme, themeMode } = storeToRefs(themeStore);

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

const onLocaleChange = (value: string | number | Record<string, unknown> | undefined) => {
  if (value === 'zh-CN' || value === 'en-US') {
    setAppLocale(value as AppLocale);
    localeMenuVisible.value = false;
  }
};

const onThemeChange = (value: string | number | Record<string, unknown> | undefined) => {
  if (value === 'light' || value === 'dark' || value === 'system') {
    themeStore.setThemeMode(value as ThemeMode);
    themeMenuVisible.value = false;
  }
};
</script>

<template>
  <div class="auth-shell tcloud-auth-shell">
    <div class="auth-toolbar">
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
    </div>

    <section class="auth-showcase">
      <div class="auth-showcase-brand">
        <strong>{{ t('layout.productName') }}</strong>
      </div>

      <p>{{ t('auth.layoutDescription') }}</p>
    </section>

    <section class="auth-stage">
      <div class="auth-stage-inner">
        <router-view />
      </div>
    </section>
  </div>
</template>
