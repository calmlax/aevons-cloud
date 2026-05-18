# Aevons

一个基于 Vue 3、Vite、Vue Router 和 Arco Design Vue 的后台管理示例项目。

项目当前已经集成以下能力：

- 明暗主题切换，支持浅色、深色、跟随系统
- 中英文切换
- Mock 登录、注册、邮箱验证码、找回密码、修改密码
- 路由权限、菜单权限、按钮权限
- 菜单由 JSON 驱动，并通过接口风格的封装获取
- 组件示例页面，当前包含 Button 和 Avatar

## 技术栈

- Vue 3
- Vite
- TypeScript
- Vue Router 4
- Vue I18n 9
- Arco Design Vue

## 功能概览

### 1. 后台框架

- 左侧菜单支持多级嵌套
- 菜单支持图标，一级和子菜单都可配置 icon
- 顶部工具栏支持语言切换、主题切换、全屏、菜单搜索、账户操作
- 面包屑和页面标题随当前路由自动变化

### 2. 主题系统

- 主题模式：`light`、`dark`、`system`
- Arco 组件主题通过 `body[arco-theme="dark"]` 切换
- 自定义页面样式通过 `html[data-theme="light|dark"]` 控制
- 主题状态保存在浏览器 `localStorage`

### 3. 国际化

- 当前提供：简体中文、English
- 每种语言拆分为单独文件维护
- 菜单、认证页、结果页、组件示例页均接入 i18n

### 4. 认证与权限

- 登录
- 注册
- 邮箱验证码
- 找回密码
- 修改密码
- 403 权限不足页，并支持定时返回首页
- 路由权限控制
- 菜单权限过滤
- 按钮级权限指令 `v-permission`

### 5. Mock 数据能力

- 菜单数据位于 `src/mock/menu.json`
- 用户数据位于 `src/mock/auth-users.json`
- Mock Auth API 位于 `src/api/auth.ts`
- Mock Menu API 位于 `src/api/menu.ts`
- 所有认证状态、验证码、会话等通过 `localStorage` 模拟真实环境

### 6. 组件示例

当前已提供以下组件演示页面：

- Button
- Avatar

后续可以按同样模式继续扩展常用组件目录。

## 演示账号

项目内置了 3 个 mock 账号，可直接登录体验不同权限。

| 角色 | 姓名 | 邮箱 | 密码 |
| --- | --- | --- | --- |
| 超级管理员 | 安夏 | `admin@aevo.local` | `Admin@123` |
| 运营管理员 | 周岚 | `operator@aevo.local` | `Operator@123` |
| 只读审计员 | 顾言 | `viewer@aevo.local` | `Viewer@123` |

权限差异说明：

- 超级管理员：可访问全部示例功能
- 运营管理员：可访问仪表盘、用户中心、组件示例、账户设置
- 只读审计员：可访问仪表盘、订单中心、组件示例、个人中心

## 项目结构

根目录结构：

```text
.
├─ index.html
├─ package.json
├─ vite.config.ts
├─ tsconfig.json
├─ tsconfig.node.json
├─ src/
└─ dist/                  # 打包后生成
```

`src` 目录说明：

```text
src
├─ api/                   # 接口层封装，当前主要是 mock auth/menu API
├─ assets/
│  └─ style/              # 全局样式
├─ components/            # 通用组件，例如侧边菜单节点渲染
├─ config/                # 项目配置，例如 mock 延迟、会话过期时间
├─ directives/            # 全局指令，例如 v-permission
├─ filters/               # 预留目录
├─ hooks/                 # 组合式逻辑，如 useAuth/useMenu/useTheme
├─ layout/                # 页面布局，如 AdminLayout/AuthLayout
├─ locale/                # 国际化入口与语言包
├─ mock/                  # mock 数据源
├─ router/                # 路由与前置守卫
├─ store/                 # 预留目录
├─ types/                 # 全局类型定义
├─ utils/                 # 工具函数，如菜单图标映射、菜单路由处理
├─ views/                 # 页面视图
├─ App.vue
└─ main.ts
```

## 关键目录说明

### `src/api`

- `auth.ts`：Mock 认证接口
- `menu.ts`：菜单获取与权限过滤

### `src/hooks`

- `useAuth.ts`：认证状态、登录登出、权限判断
- `useMenu.ts`：菜单初始化、刷新、重置
- `useTheme.ts`：主题模式管理
- `useLocale` 相关逻辑位于 `src/locale`

### `src/layout`

- `AdminLayout.vue`：后台主布局
- `AuthLayout.vue`：登录 / 注册 / 找回密码布局

### `src/views`

主要页面包含：

- 仪表盘
- 用户列表
- 订单列表
- 多级菜单示例页
- 个人中心
- 修改密码
- 登录 / 注册 / 找回密码
- 结果页 / 403 页
- Button / Avatar 组件示例页

### `src/mock`

- `auth-users.json`：内置演示账号
- `menu.json`：菜单结构、权限、iconKey 配置

## 菜单配置说明

菜单采用 JSON 配置，示例来源：`src/mock/menu.json`。

每个菜单节点支持：

- `key`：菜单唯一标识
- `labelKey`：i18n 文案 key
- `path`：菜单跳转路径
- `permissions`：权限列表
- `iconKey`：图标 key
- `children`：子菜单

菜单图标映射在 `src/utils/menu.ts` 中维护。

如果需要给子菜单增加图标：

1. 在 `src/mock/menu.json` 给对应节点增加 `iconKey`
2. 在 `src/utils/menu.ts` 中补充对应图标映射

## 启动项目

### 1. 安装依赖

项目当前带有 `package-lock.json`，默认使用 npm：

```bash
npm install
```

### 2. 启动开发环境

```bash
npm run dev
```

启动后按 Vite 默认方式在本地打开开发地址。

### 3. 预览生产构建

```bash
npm run preview
```

## 打包项目

执行生产打包：

```bash
npm run build
```

打包命令实际执行：

```bash
vue-tsc --noEmit && vite build
```

说明：

- 先执行 TypeScript 类型检查
- 再执行 Vite 生产构建
- 输出目录为 `dist/`

## 开发说明

### 主题切换

- 主题初始化入口：`src/hooks/useTheme.ts`
- Arco 暗黑模式：通过设置 `body` 上的 `arco-theme="dark"`
- 自定义页面样式：通过设置 `html[data-theme]`

### 权限控制

- 路由守卫：`src/router/index.ts`
- 菜单过滤：`src/api/menu.ts`
- 按钮权限：`src/directives/permission.ts`

### 认证状态

- 用户、会话、验证码都保存在浏览器 `localStorage`
- 刷新页面后仍会保留当前登录态，直到 session 过期或主动退出

## 已验证命令

当前项目已验证通过的构建命令：

```bash
npm run build
```

## 后续扩展建议

可以继续沿着当前结构扩展：

- 增加更多常用组件示例页
- 接入真实后端接口替换 mock API
- 增加 Pinia 等状态管理
- 拆分更细的业务模块和权限模型
