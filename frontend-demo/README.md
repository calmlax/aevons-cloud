# Aevons Frontend Demo

`aevons-cloud/frontend-demo` 是一个最小可运行的第三方应用示例，用来演示如何通过 Aevons 授权中心完成 OAuth2 授权码模式登录。

它的职责很单一：

- 打开授权中心登录页
- 接收授权回调里的 `code`
- 使用 `client_id + client_secret` 换取 Token
- 在页面上展示登录结果

## 目录说明

- [src/views/HomeView.vue](aevons-cloud/frontend-demo/src/views/HomeView.vue:1)
  第三方应用首页，负责发起授权和接收回调消息
- [src/views/CallbackView.vue](aevons-cloud/frontend-demo/src/views/CallbackView.vue:1)
  OAuth2 回调页，负责接收 `code` 并通过 `postMessage` 回传给首页
- [src/api/auth.ts](aevons-cloud/frontend-demo/src/api/auth.ts:1)
  用授权码换取 Token
- [src/stores/auth.ts](aevons-cloud/frontend-demo/src/stores/auth.ts:1)
  保存 `client_id`、`client_secret` 和登录后的 `tokenPair`
- [vite.config.ts](aevons-cloud/frontend-demo/vite.config.ts:1)
  Vite 启动端口和 `/api` 代理配置
- [.env](aevons-cloud/frontend-demo/.env:1)
  本地开发环境变量

## 授权流程

当前示例走的是标准授权码模式：

1. 在首页点击“使用授权中心登录”
2. 浏览器打开授权中心页面：
   `VITE_AUTHORIZE_CENTER`
3. 授权中心完成登录和授权确认后，跳回：
   `VITE_REDIRECT_URI`
4. 回调页拿到 `code`
5. 回调页通过 `postMessage` 把 `code` 发回首页
6. 首页调用 `/api/auth/v1/login`
   - `grant_type=authorization_code`
   - `code=<授权码>`
   - `Authorization: Basic base64(client_id:client_secret)`
7. 后端返回 TokenPair，前端保存到本地状态

注意：

- 授权码是一次性的，消费后立即失效
- 这个 demo 会自动消费授权码，不适合拿同一个 `code` 再手工重复请求

## 环境变量

本地默认配置见 [.env](aevons-cloud/frontend-demo/.env:1)：

```env
VITE_OAUTH_CLIENT_ID=auth-demo
VITE_OAUTH_CLIENT_SECRET=123456
VITE_API_PROXY_TARGET=http://localhost:11080
VITE_AUTHORIZE_CENTER=http://localhost:5173/oauth2/authorize
VITE_REDIRECT_URI=http://localhost:5174/callback
```

字段说明：

- `VITE_OAUTH_CLIENT_ID`
  第三方应用客户端 ID
- `VITE_OAUTH_CLIENT_SECRET`
  第三方应用客户端密钥
- `VITE_API_PROXY_TARGET`
  本地开发时 `/api` 的代理目标，默认走网关入口 `11080`
- `VITE_AUTHORIZE_CENTER`
  授权中心页面地址，通常是主前端的 `/oauth2/authorize`
- `VITE_REDIRECT_URI`
  第三方应用回调地址，必须与 `sys_oauth_client.web_server_redirect_uri` 完全一致

## 启动

安装依赖：

```bash
cd aevons-cloud/frontend-demo
pnpm install
```

启动开发服务：

```bash
pnpm dev
```

默认访问地址：

- `http://localhost:5174`
- 局域网也可访问，例如：`http://192.168.0.102:5174`

当前 Vite 已配置：

- `host: 0.0.0.0`
- `port: 5174`

## 前置依赖

要跑通这套登录链，至少需要这些服务是可用的：

1. `frontend`
   - 提供授权中心页面
   - 默认：`http://localhost:5173`

2. `gateway-service` 或 APISIX 网关
   - 提供 `/api/auth/v1/*`
   - 默认：`http://localhost:11080`

3. `auth-service`
   - 提供 OAuth2 授权和换 token 能力

4. `sys_oauth_client` 数据
   - 必须存在 `auth-demo` 这个客户端
   - `authorized_grant_types` 必须包含 `authorization_code`
   - `web_server_redirect_uri` 必须和 `VITE_REDIRECT_URI` 完全一致

## 数据库配置要求

`auth-demo` 这条客户端至少需要满足：

- `client_id = auth-demo`
- `client_secret` 对应 `123456`
- `authorized_grant_types` 包含 `authorization_code`
- `resources` 至少允许访问 `auth-service`
- `web_server_redirect_uri = http://localhost:5174/callback`

如果你改成局域网访问，也要同步改成完全一致的回调地址，例如：

```text
http://192.168.0.102:5174/callback
```

## 常见问题

### 1. `oauth2.redirect_uri_mismatch`

原因：

- `VITE_REDIRECT_URI` 和数据库里的 `web_server_redirect_uri` 不一致

解决：

- 改前端 `.env`
- 或改数据库 `sys_oauth_client.web_server_redirect_uri`

两边必须完全一致，包括：

- 协议
- host
- 端口
- path

### 2. `oauth2.invalid_client`

原因：

- `client_id` 不存在
- `client_secret` 不匹配
- `authorized_grant_types` 没有 `authorization_code`

解决：

- 检查 `sys_oauth_client`

### 3. 页面能打开，但换 token 失败

先检查：

- `VITE_API_PROXY_TARGET` 是否正确指向当前网关入口
- 当前默认应为：
  `http://localhost:11080`

### 4. 只能 `localhost:5174` 访问，不能局域网 IP 访问

当前已支持局域网访问，因为：

- [vite.config.ts](aevons-cloud/frontend-demo/vite.config.ts:1)
  配置了：

```ts
server: {
  host: '0.0.0.0',
  port: 5174,
}
```

如果仍然无法访问，请检查：

- 本机实际 IP 是否正确
- 防火墙是否拦截 `5174`

## 与主前端的关系

- `aevons-cloud/frontend`
  是授权中心前端
- `aevons-cloud/frontend-demo`
  是第三方应用 demo

两者配合后，可以完整演示：

- 登录授权
- 授权确认
- 回调拿 code
- code 换 token

