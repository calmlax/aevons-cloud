import enUS from './locales/en-US';
import zhCN from './locales/zh-CN';

export type AppLocale = 'zh-CN' | 'en-US';

const messages = {
  'zh-CN': zhCN,
  'en-US': enUS,
} as const;

export default messages;