# Aevons Gateway 自研网关设计文档

## 1. 设计目标

`gateway-service` 作为 Aevons 微服务体系的统一入口，负责承接来自浏览器、管理后台、移动端、第三方客户端等流量，并在进入后端业务服务前完成统一治理。

网关的核心目标如下：

1. 统一接入入口，屏蔽后端服务拆分细节
2. 统一认证、鉴权、限流、审计和安全防护
3. 基于 Consul 服务发现实现动态路由与负载转发
4. 支持客户端级资源访问控制，避免不同客户端越权访问不属于自己的资源
5. 为后续灰度发布、熔断降级、流量观测和策略扩展预留能力

## 2. 设计原则

1. 网关只做入口治理，不承载复杂业务逻辑
2. 认证、客户端授权、资源控制在网关统一前置
3. 路由目标不写死地址，统一通过 Consul 服务发现
4. 访问控制规则要可配置、可审计、可扩展
5. 默认拒绝，显式放行
6. 客户端权限控制优先于服务内部业务权限控制

## 3. 网关职责边界

网关负责：

- TLS 终止与统一接入
- 请求路由与服务转发
- 身份令牌解析与基础用户上下文注入
- 客户端身份识别
- 客户端资源访问控制
- Swagger UI 聚合入口
- 请求级限流、黑白名单、基础风控
- 请求日志、追踪标识、审计埋点
- 错误格式统一

网关不负责：

- 复杂业务编排
- 领域内权限判断细节
- 业务数据聚合计算
- 替代下游服务的最终授权决策

说明：

网关负责“能不能进入这个资源入口”，下游服务负责“进入后能不能执行具体业务操作”。

## 4. 目标架构

建议结构如下：

```text
Client
  -> gateway-service
      -> 认证解析
      -> 客户端识别
      -> 资源访问控制
      -> 路由匹配
      -> Consul 服务发现
      -> 反向代理转发
          -> auth-service
          -> system-service
          -> log-service
          -> ...
```

网关内部建议拆分为以下模块：

- `router`
  负责匹配请求路径、方法和目标服务

- `discovery`
  负责通过 Consul 获取目标服务健康实例

- `proxy`
  负责请求转发、Header 透传、响应回写

- `auth`
  负责令牌解析、用户上下文提取

- `clientauth`
  负责客户端身份识别与客户端资源访问控制

- `middleware`
  负责请求日志、限流、追踪、异常恢复等通用治理

- `config`
  负责网关配置、路由配置、客户端资源规则配置

- `swagger`
  负责聚合各服务的 OpenAPI 文档，并统一提供 Swagger UI 入口

## 5. 请求处理流程

建议请求处理顺序如下：

1. 接收请求并生成 `request_id`
2. 识别客户端身份
3. 校验客户端是否允许访问当前资源
4. 校验用户令牌并提取用户上下文
5. 执行基础限流和安全策略
6. 匹配目标路由
7. 通过 Consul 获取目标服务实例
8. 选择一个健康实例转发请求
9. 回写响应并记录访问日志

说明：

客户端资源访问控制应在转发前完成，避免无效请求进入后端服务。

## 6. 服务路由设计

网关配置不建议细化到“每个接口写一条路由”，而应以“服务维度”定义接入规则。

也就是说：

- `system` 指向 `system-service`
- `auth` 指向 `auth-service`
- `log` 指向 `log-service`

网关根据服务级配置识别请求前缀，再转发到目标服务。

建议每个服务只配置以下核心信息：

- 服务配置 ID
- 服务名
- 匹配前缀
- 是否携带访问令牌
- 负载策略
- 排除认证/校验的路由
- 是否启用 Consul 注册发现

示例模型：

```yaml
services:
  - id: auth
    name: auth-service
    prefix: /api/v1/auth/**
    pass_access_token: true
    discovery: consul
    load_balance: round_robin
    exclude_auth_routes:
      - POST:/api/v1/auth/login
      - POST:/api/v1/auth/refresh
      - POST:/api/v1/auth/email/code

  - id: system
    name: system-service
    prefix: /api/v1/system/**
    pass_access_token: true
    discovery: consul
    load_balance: random
    exclude_auth_routes:
      - GET:/api/v1/system/public/**
```

说明：

1. `id` 用于网关内部标识和日志记录
2. `name` 用于服务发现时匹配服务名
3. `name` 同时作为 Consul 注册名和发现名
4. `prefix` 支持通配符
5. `pass_access_token` 控制是否向下游透传访问令牌
6. `exclude_auth_routes` 用于配置排除认证或排除访问控制的路由

## 7. 服务发现、Consul 注册与转发

网关不应写死下游地址，所有下游服务统一通过 Consul 服务发现。

建议规则：

1. 服务配置中写逻辑服务名，例如 `system-service`
2. 请求转发前通过 Consul 查询健康实例列表
3. 默认支持 `random` 和 `round_robin`
4. 服务实例不可用时自动剔除
5. 同时支持 HTTP 服务端口和可选 gRPC 端口元信息
6. 网关自身也应支持注册到 Consul，便于统一治理与健康探测

转发时建议透传以下 Header：

- `X-Request-Id`
- `X-Forwarded-For`
- `X-Forwarded-Proto`
- `X-Client-Id`
- `X-User-Id`
- `X-Username`
- `Authorization`

其中 `Authorization` 是否透传，应由服务级配置 `pass_access_token` 控制。

## 8. 客户端资源访问控制设计

### 8.1 背景

微服务体系中，不同客户端通常拥有不同资源访问范围。

例如：

- 管理后台可以访问全部管理接口
- 门户站点只允许访问公开接口和部分用户接口
- 第三方集成客户端只能访问分配给它的开放 API

如果没有客户端资源访问控制，任意合法客户端都可能访问不属于自己的资源。

### 8.2 控制目标

客户端资源访问控制用于回答这个问题：

“当前客户端是否允许访问当前请求资源？”

这里的“客户端”可以是：

- OAuth Client
- 内部系统调用方
- 第三方应用
- 设备侧应用

### 8.3 资源定义

资源建议统一定义为：

`HTTP Method + Path Pattern`

例如：

- `GET:/api/v1/system/conf/**`
- `POST:/api/v1/auth/login`
- `DELETE:/api/v1/log/oper-log/**`

如果后续需要更细粒度，也可增加：

- 服务名
- 资源类型
- 资源标签

### 8.4 客户端资源规则模型

建议设计如下字段：

```yaml
clients:
  - client_id: admin-web
    enabled: true
    resources:
      - ALL

  - client_id: portal-web
    enabled: true
    resources:
      - POST:/api/v1/auth/login
      - POST:/api/v1/auth/refresh
      - GET:/api/v1/auth/user
      - GET:/api/v1/conf/public/**

  - client_id: third-party-a
    enabled: true
    resources:
      - GET:/api/v1/open/**
      - POST:/api/v1/open/token
```

### 8.5 `ALL` 语义定义

`ALL` 表示：

当前客户端允许访问网关侧定义的全部资源入口。

语义要求如下：

1. `ALL` 只表示客户端资源入口全开放
2. `ALL` 不绕过用户认证
3. `ALL` 不绕过下游服务自己的业务权限校验
4. `ALL` 不代表系统内部超级管理员
5. `ALL` 仅在客户端维度表示“允许通过网关访问全部资源”

也就是说：

- `ALL` 通过的是“客户端入口控制”
- 不是通过“业务权限控制”

### 8.6 匹配规则

建议按以下顺序匹配：

1. 客户端是否存在且启用
2. 资源列表中是否包含 `ALL`
3. 是否存在精确匹配 `METHOD:/path`
4. 是否存在路径模式匹配 `METHOD:/path/**`
5. 若均不匹配，则拒绝访问

路径模式应支持通配符，例如：

- `/api/v1/auth/**`
- `/api/v1/system/conf/**`
- `/api/v1/log/**`

建议第一阶段至少支持：

- 精确匹配
- 前缀通配 `/**`
- `ANY:/path/**` 形式的任意方法匹配

## 9. Swagger UI 聚合设计

可以在网关中统一集成 Swagger UI，并聚合所有服务的 `api/swagger.json` 文档。

这样做的好处是：

1. 前端和测试人员只需要访问网关一个入口
2. 所有服务文档统一展示，减少切换成本
3. 文档访问控制可以由网关统一收口
4. 可以根据环境开关决定是否对外暴露 Swagger UI

建议规则：

1. 各服务继续各自维护 `api/swagger.json`
2. 各服务通过自身 `/api/swagger.json` 对外暴露文档
3. 网关负责聚合这些 OpenAPI 地址
4. 统一通过一个 Swagger UI 页面展示多个服务文档

实现上建议直接使用现成 Swagger UI 库，而不是自己开发页面。

推荐方式：

- 使用 `swagger-ui-dist`
- 或其他兼容 OpenAPI 3 的现成 Swagger UI 静态资源

不建议：

- 手写前端页面解析 OpenAPI
- 在网关里自己实现 Swagger 渲染逻辑

### 9.1 访问开关

Swagger UI 必须支持配置开关控制。

建议开关语义：

- `swagger.enabled = true`
  允许暴露当前服务的 `/api/swagger.json`

- `swagger.ui_enabled = true`
  允许暴露网关 Swagger UI 页面

也就是说：

- 只开 `enabled`，可以访问 OpenAPI JSON
- 同时开 `enabled + ui_enabled`，才可以访问 Swagger UI

默认建议：

- 开发环境：可开启
- 测试环境：按需开启
- 生产环境：默认关闭

### 9.2 网关聚合方式

网关不应内嵌所有 OpenAPI 内容，而应聚合“文档地址列表”。

例如：

```yaml
swagger:
  enabled: true
  ui_enabled: true
  docs:
    - name: auth-service
      url: /auth/api/swagger.json
    - name: system-service
      url: /system/api/swagger.json
    - name: log-service
      url: /log/api/swagger.json
```

Swagger UI 通过这个列表动态展示多个服务文档。

### 9.3 安全要求

Swagger UI 虽然是文档能力，但同样属于入口能力，需要统一治理。

建议：

1. 必须受配置开关控制
2. 可以额外受 IP 白名单控制
3. 可以只允许开发、测试、内网环境访问
4. 相关访问请求应纳入网关日志审计

### 9.4 与客户端资源控制的关系

Swagger UI 属于网关管理能力，不建议混入普通业务客户端资源规则中。

更推荐的处理方式是：

- 通过 Swagger 专用开关控制是否开放
- 通过网关附加安全策略控制访问范围

而不是把 Swagger UI 当成普通业务 API 去做客户端资源授权。

### 8.7 默认策略

默认策略建议为：

- 未识别客户端：拒绝
- 客户端被禁用：拒绝
- 没有资源配置：拒绝
- 没有命中资源规则：拒绝

返回码建议：

- `401`：客户端身份无效
- `403`：客户端无权访问当前资源

## 9. 认证与客户端识别

网关需要区分两类身份：

1. 用户身份
2. 客户端身份

客户端身份可来自：

- `client_id`
- Basic Auth
- JWT 中的 `aud/client_id`
- 内部服务调用凭证
- API Key

建议优先统一成网关内部上下文：

```text
request context:
  client_id
  user_id
  username
  roles
  permissions
```

客户端资源访问控制主要依赖 `client_id`。

## 10. 配置来源建议

客户端资源规则不建议硬编码在代码中，建议支持以下来源：

1. 本地配置文件
2. 配置中心
3. 数据库

演进建议：

- 初期可用 YAML 配置
- 中期可迁到数据库或配置中心
- 网关本地缓存配置，并支持热更新

## 11. 数据结构建议

建议抽象以下核心结构：

```go
type RouteRule struct {
    ID                  string
    Method              string
    Path                string
    Service             string
    AuthRequired        bool
    ClientResourceCheck bool
}

type ClientResourceRule struct {
    ClientID  string
    Enabled   bool
    Resources []string
}
```

建议在内存中建立：

- 路由表
- 客户端资源规则表
- 路径匹配索引

## 12. 鉴权决策顺序

建议网关内的最终决策顺序如下：

1. 路由是否存在
2. 客户端是否有效
3. 客户端是否允许访问当前资源
4. 是否要求用户认证
5. 用户令牌是否有效
6. 转发到目标服务

这样可以避免：

- 未授权客户端先进入业务认证
- 无效流量进入下游服务

## 13. 安全与审计要求

网关应记录以下审计信息：

- 请求时间
- `request_id`
- 客户端 ID
- 用户 ID
- 请求方法与路径
- 匹配到的目标服务
- 是否命中 `ALL`
- 是否被资源规则拒绝
- 响应状态码
- 请求耗时

对 `ALL` 客户端应额外关注：

- 使用范围
- 配置变更审计
- 谁在何时赋予了 `ALL`

原则上：

`ALL` 应只授予受控客户端，不应随意发给普通第三方应用。

## 14. 故障与降级策略

当 Consul 或下游服务异常时，网关应具备基本降级能力：

- 找不到健康实例时返回统一错误
- 可对部分公开接口做静态降级
- 发现服务不可达时快速失败
- 对查询型接口可考虑短时缓存策略

客户端资源规则加载异常时，建议：

- 默认拒绝高风险接口
- 对明确允许匿名访问的接口可按白名单放行

## 15. 实施建议

建议按以下顺序实现：

1. 完成基础反向代理与路由模块
2. 接入 Consul 服务发现
3. 接入统一认证解析
4. 实现客户端资源规则与 `ALL` 匹配
5. 增加请求日志与审计
6. 增加限流、黑名单、熔断等治理能力

## 16. 结论

`gateway-service` 的核心价值不只是“转发请求”，而是：

- 统一入口
- 统一治理
- 统一客户端资源访问控制

其中，客户端资源访问控制是微服务入口层必须具备的能力之一。

`ALL` 的设计可以提升配置效率，但必须严格限定它的语义：

- 它代表客户端资源入口全部可访问
- 不代表绕过认证
- 不代表绕过业务权限
- 不代表超级管理员

这样才能在保持灵活性的同时，避免权限失控。
