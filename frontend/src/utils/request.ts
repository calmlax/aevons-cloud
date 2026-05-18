import axios, { type AxiosRequestConfig, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';
import { Message } from '@arco-design/web-vue';
import { i18n ,getAppLocale} from '@/locale';
import { useAuthStore } from '@/store/modules/auth';
import { Session, Local } from './cache';
import { tansParams, blobValidate } from './common';
import { saveAs } from 'file-saver';

// 全局扩展 Axios 原生配置对象的类型定义，加入自定义的拦截器控制标志位
declare module 'axios' {
  export interface InternalAxiosRequestConfig {
    _retry?: boolean;            // 标记该请求是否已经历过 401 自动刷新无感重试，防止无限死循环
    _isRefreshRequest?: boolean; // 标记该请求是否其本身就是“去后端换取新 Token”的刷新请求
    isToken?: boolean;           // 是否需要携带 Token，false 表示忽略 Token 附加（如登录接口等）
    repeatSubmit?: boolean;      // 是否开启防重复提交控制（默认开启，false 关闭）
  }
}

// 基础网关地址：可以配置在 .env 系列文件中，如 VITE_API_BASE_URL=/api
const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

// 创建 axios 全局单例实例
const service = axios.create({
  baseURL: BASE_URL,
  timeout: 10000, // 设定请求超时时间为 10 秒
  headers: {
    'Content-Type': 'application/json;charset=utf-8', // 默认使用 JSON 数据交换格式
  },
  withCredentials: true, // 这是双 Token 必须的核心：允许跨域自动携带 HttpOnly Cookie（如 RefreshToken）给服务端
});

// ==========================================
// 并发 401 无感刷新锁机制用的全局状态变量
// ==========================================
let isRefreshing = false; // "正在刷新 Token" 的进程互斥锁
let refreshSubscribers: ((newToken: string) => void)[] = []; // 因为撞到 401 而排队休眠的请求列队（发誓等新 Token 来）

// Helper: 将请求压入等待队列
const subscribeTokenRefresh = (cb: (token: string) => void) => {
  refreshSubscribers.push(cb);
};

// Helper: 唤醒并执行队列中的所有请求
const onRefreshed = (token: string) => {
  refreshSubscribers.forEach((cb) => cb(token));
  refreshSubscribers = []; // 清空队列
};

// ==========================================
// [1] Axios 请求拦截器 (Request Interceptor)
// 发生在上送数据去服务器之前
// ==========================================
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {

    // --> 1. 动态附加身份令牌 (Access Token)
    const authStore = useAuthStore();
    const token = authStore.accessToken;
    // 检查此请求是否在白名单中配置了 isToken: false（不需要身份验证）
    const isToken = (config.headers as any).isToken === false;

    config.headers = config.headers || {};

    // 如果本地内存持有 Token，且该请求需要 Token
    if (token && !isToken) {
      config.headers.Authorization = `Bearer ${token}`; // 在 Header 中带上标准 OAuth2 Barer 证书
    }

    // --> 2. 国际化支持：向后端传递当前客户端语言环境
    config.headers['Accept-Language'] = getAppLocale();

    // --> 3. 处理 GET 请求的对象参数 (避免嵌套对象直接 stringify 变成 "[object Object]")
    if (config.method === 'get' && config.params) {
      // 通过自定义扁平化拍平工具，安全序列化参数拼接到 URL 尾巴上
      let url = config.url + '?' + tansParams(config.params);
      url = url.slice(0, -1);
      config.params = {};
      config.url = url;
    }

    // --> 4. 防抖与防重复提交机制 (防黑客重放/狂按点赞)
    const isRepeatSubmit = (config.headers as any).repeatSubmit === false || (config.headers as any).repeatSubmit === 'false';
    // 只有 Post 或 Put (增、改操作) 才需要被拦截
    if (!isRepeatSubmit && (config.method === 'post' || config.method === 'put')) {
      const requestObj = {
        url: config.url,
        data: typeof config.data === 'object' ? JSON.stringify(config.data) : config.data,
        time: new Date().getTime()
      }
      const sessionObj = Session.get<any>('sessionObj'); // 从本地 Session 获取上一次的历史提交信息

      // 如果内存为空，代表这是第一回提交
      if (sessionObj === undefined || sessionObj === null || sessionObj === '') {
        Session.set('sessionObj', requestObj)
      } else {
        const s_url = sessionObj.url;                // 上一次请求地址
        const s_data = sessionObj.data;              // 上一次请求载体(Payload)
        const s_time = sessionObj.time;              // 上一次请求发生的时间戳
        const interval = 1000;                       // 判定时间窗：同一请求在 1000 毫秒(1秒)内发出，即视为狂按

        // 地址相同、数据内容相同、时间差小于间隔，实锤异常！拦截！
        if (s_data === requestObj.data && requestObj.time - s_time < interval && s_url === requestObj.url) {
          const message = i18n.global.t('sys.api.repeatSubmit');
          console.warn(`[${s_url}]: ` + message);
          Message.warning(message);
          return Promise.reject(new Error(message));
        } else {
          Session.set('sessionObj', requestObj); // 更新正常的新鲜请求记录
        }
      }
    }

    return config;
  },
  (error) => {
    // 构建请求对象报错时的通用捕获回退
    console.error(error);
    return Promise.reject(error);
  }
);

// ==========================================
// [2] Axios 响应拦截器 (Response Interceptor)
// 发生在收到了服务器返回的数据之后
// ==========================================
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, data, msg, message, args } = response.data || {};
    // 兼容取值 msg 或 message (视公司不同后端语言习惯)
    let effectiveMsg = msg || message || '';

    // 尝试进行多语言翻译
    if (effectiveMsg) {
      const translated = i18n.global.t(effectiveMsg, args || {});
      if (translated !== effectiveMsg) {
        effectiveMsg = translated;
      }
    }

    // 如果请求的是二进制文件流 (如 Excel 下载或图片下载)，直接不拆分外壳整体交还
    if (response.config.responseType === 'blob' || response.config.responseType === 'arraybuffer') {
      return response.data;
    }

    // 回退获取有效的响应层级状态码 (例如有的后端业务跑飞直接抛空 code，则回落看 Http 报文 200)
    const effectiveCode = code || response.status || 200;

    // --> 业务状态拦截分发器
    if (effectiveCode === 401) {
      // 第一时间：如果业务报文判定属于 401 会话失效
      // 我们抛出一个特殊的错误信息直接拒绝解析包体，把它丢给下方的 error catch (去享受无感刷新的队列排队服务)
      const err = new Error('401_ERR_TRIGGER');
      (err as any).response = { status: 401 };
      (err as any).config = response.config;
      return Promise.reject(err);
    } else if (effectiveCode === 500) {
      // 第二种情况：服务器内部宕机或抛出异常
      Message.error(effectiveMsg || i18n.global.t('sys.api.serverError'));
      return Promise.reject(new Error(effectiveMsg));
    } else if (effectiveCode === 601) {
      // 第三种情况：标准业务处理警示阻断 (如余额不足、状态不对等)
      Message.warning(effectiveMsg || i18n.global.t('sys.api.operationFailed'));
      return Promise.reject(new Error(effectiveMsg));
    } else if (effectiveCode === 200 || effectiveCode === 0 || effectiveCode === '200' || effectiveCode === '0') {
      // 第四种情况：完全成功！剥下外壳返回最里面的核心 data
      return data !== undefined ? data : response.data;
    } else {
      // 第五种情况：未在规范内的错误定义，强制托底抛错
      Message.error(effectiveMsg || i18n.global.t('sys.api.requestFailed'));
      return Promise.reject(new Error(effectiveMsg || 'error'));
    }
  },
  async (error) => {
    const config = error.config as InternalAxiosRequestConfig | undefined;
    const status = error.response?.status;

    if (!config) {
      Message.error(i18n.global.t('sys.api.networkError'));
      console.error('请求异常：无配置信息', error);
      return Promise.reject(error);
    }

    if (error.code === 'ECONNABORTED' || error.message.includes('timeout')) {
      Message.error(i18n.global.t('sys.api.timeout'));
      return Promise.reject(error);
    }

    const skipUrls = ['/v1/auth/login', '/v1/auth/refresh', '/v1/auth/logout', '/auth/login', '/auth/refresh', '/auth/logout'];
    const isAuthEndpoint = config.url && skipUrls.some((url) => config.url?.includes(url));
    // ====== 开启短效 401 无感刷新 Token 的守护动作 ======
    // 判断条件：撞墙状态为401，且该次请求它还没有被自动抢救过 (!config._retry) 以及并非核心接口
    if (status === 401 && !config._retry && !isAuthEndpoint) {

      // 如果早已经有另外一个并行的兄弟接口先撞了 401，且已经在去获取新 Token 的路上了...
      if (isRefreshing) {
        // 当前接口不准再向公司后端发起重复刷新，立刻将它包装进一个挂起的 Promise 契约丢入列队，并沉睡等待
        return new Promise((resolve) => {
          subscribeTokenRefresh((newToken) => {
            // 当 Token 拿到时我们这行代码会被唤醒！带上新 Token 发车。
            if (config.headers) {
              config.headers.Authorization = `Bearer ${newToken}`;
            }
            config._retry = true;
            resolve(service(config));
          });
        });
      }

      // 兄弟们都没去刷新，那好咱们当前这个请求就是首当其冲的那个。标记我们自己负责去刷新。
      config._retry = true;
      isRefreshing = true; // 举起互斥锁，拦截后续发车请求

      try {
        const clientId = import.meta.env.VITE_OAUTH_CLIENT_ID || '';
        const clientSecret = import.meta.env.VITE_OAUTH_CLIENT_SECRET || '';
        const baseURL = import.meta.env.VITE_API_BASE_URL || '/api';

        // 利用浏览器的双凭证特性，让 纯净的 Axios（避免拦截器死循环）带上 Basic Auth 标准头去换取新 AccessToken
        const refreshRes = await axios.post(
          `${baseURL}/v1/auth/refresh`,
          {},
          {
            headers: {
              Authorization: `Basic ${btoa(`${clientId}:${clientSecret}`)}`
            },
            withCredentials: true
          }
        );

        // 解析出新的 accessToken (兼容外壳嵌套或直接返回)
        const newAccessToken = refreshRes.data?.data?.access_token || refreshRes.data?.access_token;

        if (!newAccessToken) {
          throw new Error('Refresh failed, no token returned');
        }

        // 【内存同步】：拿到之后迅速灌入 Pinia，保证全网页实时最新
        const authStore = useAuthStore();
        authStore.setToken(newAccessToken as string);

        // 【队列释放】：释放并唤醒刚才排队在 `subscribeTokenRefresh` 沉睡的所有请求（让它们携带新 Token 重发）
        onRefreshed(newAccessToken as string);

        // 最后别忘了把我们作为引发这次雪崩的第一个倒霉蛋的原始请求，继续用新 Token 完成发送使命
        if (config.headers) {
          config.headers.Authorization = `Bearer ${newAccessToken}`;
        }
        return service(config);
      } catch (refreshError) {
        // 【终极灾难】：如果负责拿 Token 的请求都报错了说明长效 HttpOnly Cookie 也死了
        // 这时没有任何退路可言，全面强制退出并阻断一切队列。
        refreshSubscribers = [];
        const authStore = useAuthStore();
        authStore.clearAuth();
        Message.error(i18n.global.t('sys.api.sessionExpired'));
        return Promise.reject(refreshError);
      } finally {
        // 记得解锁
        isRefreshing = false;
      }
    }
    // ===========================================

    // == 其他常规 HTTP 协议级的报错进行文案自动翻译补齐 ==
    let message = error.response?.data?.message || error.message;
    const args = error.response?.data?.args || {};

    let isTranslated = false;
    if (error.response?.data?.message) {
      const translated = i18n.global.t(message, args);
      if (translated !== message) {
        message = translated;
        isTranslated = true;
      }
    }

    if (!isTranslated) {
      if (status === 401) {
        message = i18n.global.t('sys.api.sessionExpired');
      } else if (status === 403) {
        message = i18n.global.t('sys.api.permissionDenied');
      } else if (status === 404) {
        message = i18n.global.t('sys.api.notFound');
      } else if (status >= 500) {
        message = i18n.global.t('sys.api.serverInternalError');
      }
    }

    // 将最终修饰或翻译完毕的消息写回 error 对象，确保业务方 catch 的是处理后的文本
    error.message = message;

    // 弹窗输出报错提示，但是要排除自己创造的内置信号阻塞，防止连环吐司
    if (!config._isRefreshRequest && message !== '401_ERR_TRIGGER') {
      Message.error(message);
    }
    return Promise.reject(error);
  }
);

let downloadLoadingInstance: any;

/**
 * 带有健壮拦截器包裹的大文件下载通用方法
 * 支持流转换识别逻辑 (如果请求被拦截抛出 json 而没有转为文件，能进行回弹识别)
 * 
 * @param url 文件接口路径
 * @param params 参数载体
 * @param filename 下载到系统硬盘时重命名的文件名
 * @param config 额外配置属性
 */
export function download(url: string, params: any, filename: string, config: any) {
  // Arco Design 全局 loading 挂载
  downloadLoadingInstance = Message.loading({
    content: "正在下载...",
    duration: 0
  });

  return service.post(url, params, {
    transformRequest: [(p) => tansParams(p)],
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    responseType: 'blob', // 明确要求浏览器将拿到的数据封存到 Blob 对象中
    ...config
  }).then(async (data: any) => {
    // 【边界防崩操作】：如果由于 Token 过期，按理说后端应该重定向 401 JSON 回来
    // 但是在这个下载请求里，不设防的 Axios 就会生硬地把这串 401 报错 JSON 保存成了个 2KB 的破损文件！
    // 这里的 blobValidate 会反过来拆解 Blob 检查它的头部究竟是啥类型。
    const isBlob = blobValidate(data);
    if (isBlob) {
      const blob = new Blob([data]);
      saveAs(blob, filename); // 妥投到系统下载栏中保存
    } else {
      // 实际上是个 JSON 封装下的报错结构？！赶紧拆包裹恢复本相
      const resText = await data.text();
      const rspObj = JSON.parse(resText);
      const errMsg = rspObj.msg || rspObj.message || '下载失败';
      Message.error(errMsg); // 优雅地通过控制台弹窗而不是生成坏文件！
    }
    // 妥善关闭加载窗闭包
    if (downloadLoadingInstance && downloadLoadingInstance.close) {
      downloadLoadingInstance.close();
    }
  }).catch((error) => {
    console.error(error);
    Message.error('下载文件失败');
    if (downloadLoadingInstance && downloadLoadingInstance.close) {
      downloadLoadingInstance.close();
    }
  });
}

// 基础 request
export const request = <T = any>(config: AxiosRequestConfig): Promise<T> => service.request(config);

// 导出你想要的 get/post/put/delete 方法 ✅
export const requestGet = <T = any>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> =>
  service.get(url, { params, ...config });

export const requestPost = <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> =>
  service.post(url, data, config);

export const requestPut = <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> =>
  service.put(url, data, config);

export const requestDelete = <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> =>
  service.delete(url, config);

export const requestPatch = <T = any>(
  url: string,
  data?: any,
  config?: AxiosRequestConfig
): Promise<T> => service.patch(url, data, config);

export const requestHead = <T = any>(
  url: string,
  config?: AxiosRequestConfig
): Promise<T> => service.head(url, config);

export const requestOptions = <T = any>(
  url: string,
  config?: AxiosRequestConfig
): Promise<T> => service.options(url, config);

// ==================== 上传文件（特别常用） ====================
export const requestUpload = <T = any>(
  url: string,
  formData: FormData,
  config?: AxiosRequestConfig
): Promise<T> =>
  service.post(url, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    ...config,
  });

// ==================== 下载文件（返回 blob） ====================
export const requestDownload = <T = Blob>(
  url: string,
  params?: any,
  config?: AxiosRequestConfig
): Promise<T> =>
  service.get(url, {
    params,
    responseType: 'blob',
    ...config,
  });

// 如果你想直接用 request.get 这种写法，也可以挂载上去
request.get = requestGet;
request.post = requestPost;
request.put = requestPut;
request.delete = requestDelete;
request.patch = requestPatch;
request.head = requestHead;
request.options = requestOptions;
request.upload = requestUpload;
request.download = requestDownload;

export default service;
