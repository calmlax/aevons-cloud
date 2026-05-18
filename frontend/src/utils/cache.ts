export interface CacheItem {
  value: unknown;
  time: number;
  expire: number | null;
}

const DEFAULT_CACHE_TIME = 60 * 60 * 24 * 7 * 1000; // 7 days
const APP_PREFIX = 'aevo_';

const getPrefixedKey = (key: string) => {
  return key.startsWith(APP_PREFIX) ? key : `${APP_PREFIX}${key}`;
};

/**
 * 通用的 WebStorage 封装类
 */
class WebStorage {
  private storage: Storage;

  constructor(storage: Storage) {
    this.storage = storage;
  }

  set(key: string, value: unknown, expire: number | null = DEFAULT_CACHE_TIME) {
    const parsedKey = getPrefixedKey(key);
    
    const data: CacheItem = {
      value,
      time: Date.now(),
      expire: expire !== null ? Date.now() + expire : null,
    };
    
    this.storage.setItem(parsedKey, JSON.stringify(data));
  }

  get<T = unknown>(key: string, def?: T): T | null {
    const parsedKey = getPrefixedKey(key);
    const raw = this.storage.getItem(parsedKey);
    
    if (!raw) {
      return def ?? null;
    }
    
    try {
      const data: CacheItem = JSON.parse(raw);
      const { value, expire } = data;
      
      // 如果设置了过期时间，并且此刻已经过期，那么清理缓存并返回 null（或者默认值）
      if (expire !== null && expire < Date.now()) {
        this.remove(key);
        return def ?? null;
      }
      
      return value as T;
    } catch (error) {
      this.remove(key);
      return def ?? null;
    }
  }

  remove(key: string) {
    const parsedKey = getPrefixedKey(key);
    this.storage.removeItem(parsedKey);
  }

  clear() {
    Object.keys(this.storage).forEach((key) => {
      // 仅清理带有前缀的缓存，防止误删其它缓存
      if (key.startsWith(APP_PREFIX)) {
        this.storage.removeItem(key);
      }
    });
  }
}

// 模拟 SSR 下的 Storage 环境以防止报错
const isBrowser = typeof window !== 'undefined';
const dummyStorage = {
  getItem: () => null,
  setItem: () => {},
  removeItem: () => {},
  clear: () => {},
  length: 0,
  key: () => null,
} as unknown as Storage;

/** 暴露出基于 LocalStorage 的完整实例 */
export const Local = new WebStorage(isBrowser ? window.localStorage : dummyStorage);

/** 暴露出基于 SessionStorage 的完整实例 */
export const Session = new WebStorage(isBrowser ? window.sessionStorage : dummyStorage);

// =============== 保留原有的直接导出方法以兼容已有代码 (指向 LocalStorage) ===============
export const setCache = Local.set.bind(Local);
export const getCache = Local.get.bind(Local);
export const removeCache = Local.remove.bind(Local);
export const clearCache = Local.clear.bind(Local);

// =============== 提供基于 SessionStorage 的快速导出方法 ===============
export const setSession = Session.set.bind(Session);
export const getSession = Session.get.bind(Session);
export const removeSession = Session.remove.bind(Session);
export const clearSession = Session.clear.bind(Session);
