import { computed } from 'vue';
import { createI18n } from 'vue-i18n';
import arcoEnUS from '@arco-design/web-vue/es/locale/lang/en-us';
import arcoZhCN from '@arco-design/web-vue/es/locale/lang/zh-cn';
import langApi from '@/api/system/lang'

import messages, { type AppLocale } from './messages';

export type { AppLocale } from './messages';

const LOCALE_STORAGE_KEY = 'aevo-locale';
const defaultLocale: AppLocale = 'zh-CN';

const isAppLocale = (value: string | null): value is AppLocale => {
  return value === 'zh-CN' || value === 'en-US';
};

const readStoredLocale = (): AppLocale => {
  if (typeof window === 'undefined') {
    return defaultLocale;
  }

  const storedLocale = window.localStorage.getItem(LOCALE_STORAGE_KEY);
  return isAppLocale(storedLocale) ? storedLocale : defaultLocale;
};

export const i18n = createI18n({
  legacy: false,
  locale: defaultLocale,
  fallbackLocale: defaultLocale,
  messages,
});

const applyDocumentLanguage = (locale: AppLocale) => {
  if (typeof document === 'undefined') {
    return;
  }

  document.documentElement.lang = locale === 'zh-CN' ? 'zh-CN' : 'en';
};

export const initializeLocale = () => {
  const nextLocale = readStoredLocale();
  i18n.global.locale.value = nextLocale;
  applyDocumentLanguage(nextLocale);
};

export const setAppLocale = (locale: AppLocale) => {
  i18n.global.locale.value = locale;
  applyDocumentLanguage(locale);

  if (typeof window !== 'undefined') {
    window.localStorage.setItem(LOCALE_STORAGE_KEY, locale);
  }
};

export const getAppLocale = (): AppLocale => {
  return i18n.global.locale.value as AppLocale;
};

export const useAppLocale = () => {
    const appLocale = i18n.global.locale;
  const arcoLocale = computed(() => appLocale.value === 'en-US' ? arcoEnUS : arcoZhCN);

  // 语言列表（响应式）
  const localeOptions = ref<any[]>([]);

  // 获取后端语言列表
  const fetchLangList = async () => {
    try {
      const langList:any = await langApi.availableList();
      if (langList) {
        localeOptions.value = langList.map((lang: any) => ({
          label: lang.langName,
          value: lang.langCode,
        }));
      } else {
        setDefaultOptions();
      }
    } catch (err) {
      console.warn('语言列表接口失败，使用默认语言', err);
      setDefaultOptions();
    }
  };

  // 默认语言
  const setDefaultOptions = () => {
    localeOptions.value = [
      { label: i18n.global.t('locale.zhCN'), value: 'zh-CN' },
      { label: i18n.global.t('locale.enUS'), value: 'en-US' },
    ];
  };

  // 组件使用时自动加载
  onMounted(() => {
    if (localeOptions.value.length === 0) fetchLangList();
  });


  return {
    appLocale,
    arcoLocale,
    localeOptions,
    setAppLocale,
    getAppLocale,
  };
};