# Aevons Gateway 实现设计草案

## 1. 文档目标

本文档用于承接 [gateway-architecture.md](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway-service/gateway-architecture.md:1) 的架构设计，进一步明确 `gateway-service` 的实现结构、配置模型与关键接口，便于后续直接进入编码阶段。

## 2. 实现目标

第一阶段建议优先实现以下能力：

1. HTTP 网关入口
2. 路由匹配
3. Consul 服务发现
4. 反向代理转发
5. 用户认证解析
6. 客户端资源访问控制
7. Swagger UI 聚合
8. 统一请求日志与错误响应

暂不优先实现：

- 熔断器
- 灰度发布
- 动态配置中心
- 多级缓存
- 管理后台可视化配置

## 3. 建议目录结构

```text
gateway-service/
├── cmd/server/
│   └── main.go
├── configs/
│   └── config.yaml
├── internal/
│   ├── config/
│   ├── model/
│   ├── router/
│   ├── discovery/
│   ├── proxy/
│   ├── clientauth/
│   ├── auth/
│   ├── swagger/
│   ├── middleware/
│   ├── response/
│   └── bootstrap/
└── gateway-architecture.md
```

## 4. 配置结构草案

建议网关配置拆成以下几部分：

```yaml
server:
  name: gateway-service
  host: 0.0.0.0
  port: 8080
  mode: release

consul:
  enabled: true
  address: http://127.0.0.1:8500
  register_host: 192.168.31.81

gateway:
  trusted_proxies:
    - 127.0.0.1
  timeout_seconds: 15
  max_body_bytes: 10485760

swagger:
  enabled: true
  ui_enabled: true
  docs:
    - name: auth-service
      url: /auth/apifox/openapi.json
    - name: system-service
      url: /system/apifox/openapi.json
    - name: log-service
      url: /log/apifox/openapi.json

services:
  - id: auth
    name: auth-service
    prefix: /api/v1/auth/**
    discovery: consul
    load_balance: round_robin
    pass_access_token: true
    exclude_auth_routes:
      - POST:/api/v1/auth/login
      - POST:/api/v1/auth/refresh
      - POST:/api/v1/auth/email/code

  - id: system
    name: system-service
    prefix: /api/v1/system/**
    discovery: consul
    load_balance: random
    pass_access_token: true
    exclude_auth_routes:
      - GET:/api/v1/system/public/**

clients:
  - client_id: admin-web
    enabled: true
    resources:
      - ALL

  - client_id: portal-web
    enabled: true
    resources:
      - POST:/api/v1/auth/login
      - GET:/api/v1/auth/user
```

## 5. 核心配置模型

建议定义以下结构：

```go
type GatewayConfig struct {
    TrustedProxies []string `yaml:"trusted_proxies"`
    TimeoutSeconds int      `yaml:"timeout_seconds"`
    MaxBodyBytes   int64    `yaml:"max_body_bytes"`
}

type SwaggerConfig struct {
    Enabled   bool               `yaml:"enabled"`
    UIEnabled bool               `yaml:"ui_enabled"`
    Docs      []SwaggerDocConfig `yaml:"docs"`
}

type SwaggerDocConfig struct {
    Name string `yaml:"name"`
    URL  string `yaml:"url"`
}

type ServiceConfig struct {
    ID                string   `yaml:"id"`
    Name              string   `yaml:"name"`
    Prefix            string   `yaml:"prefix"`
    Discovery         string   `yaml:"discovery"`
    LoadBalance       string   `yaml:"load_balance"`
    PassAccessToken   bool     `yaml:"pass_access_token"`
    ExcludeAuthRoutes []string `yaml:"exclude_auth_routes"`
}

type ClientRuleConfig struct {
    ClientID  string   `yaml:"client_id"`
    Enabled   bool     `yaml:"enabled"`
    Resources []string `yaml:"resources"`
}
```

## 6. 内存模型建议

配置加载后，不建议每次请求都直接遍历原始 YAML 结构，建议构建网关运行时模型：

```go
type ServiceRule struct {
    ID               string
    Name             string
    Prefix           string
    Discovery        string
    LoadBalance      string
    PassAccessToken  bool
    ExcludeAuthRules []string
}

type ClientRule struct {
    ClientID    string
    Enabled     bool
    AllowAll    bool
    ExactRules  map[string]struct{}
    PrefixRules []string
}
```

说明：

- `AllowAll` 用于快速判断是否命中 `ALL`
- `ExactRules` 用于 `GET:/api/v1/auth/user` 这类精确匹配
- `PrefixRules` 用于 `/api/v1/system/conf/**` 这类前缀匹配
- `ServiceRule` 用于表示一个完整服务的接入规则

## 7. 路由匹配器设计

建议定义接口：

```go
type ServiceMatcher interface {
    Match(path string) (*ServiceRule, bool)
}
```

第一阶段实现建议采用：

- 长前缀优先
- 前缀通配次之
- 不做复杂正则

建议匹配顺序：

1. 精确前缀优先
2. 长前缀优先于短前缀
3. 通配符前缀次之
4. 都不命中则返回无路由

例如：

- `/api/v1/auth/**` -> `auth-service`
- `/api/v1/system/**` -> `system-service`
- `/api/v1/log/**` -> `log-service`

## 8. 客户端资源匹配器设计

建议定义接口：

```go
type ClientResourceChecker interface {
    Allow(clientID, method, path string) bool
}
```

建议核心实现方法：

```go
func (c *Checker) Allow(clientID, method, path string) bool
```

匹配逻辑：

1. 查找客户端规则
2. 客户端不存在或禁用则拒绝
3. 如果 `AllowAll == true`，直接允许
4. 匹配 `METHOD:/path` 精确规则
5. 匹配 `METHOD:/path/**` 前缀规则
6. 都不匹配则拒绝

## 9. `ALL` 处理规则

`ALL` 建议在配置加载阶段就预处理为：

```go
AllowAll = true
```

而不是每次请求时遍历字符串列表。

这样可以减少判断开销，也能避免大小写或空白字符问题。

建议规范：

- 只认大写 `ALL`
- 配置加载时自动 `TrimSpace`
- 不允许与非法资源规则混用而静默忽略

## 10. 客户端身份识别设计

建议定义：

```go
type ClientIdentity struct {
    ClientID string
    Source   string
}
```

建议识别来源顺序：

1. `Authorization: Basic ...`
2. `X-Client-Id`
3. 用户 JWT 中的 `client_id`
4. 内部服务固定 Header

建议统一收口方法：

```go
type ClientResolver interface {
    Resolve(r *http.Request) (*ClientIdentity, error)
}
```

## 11. Swagger 聚合实现设计

建议在网关中增加一个独立的 `swagger` 模块，专门负责 Swagger UI 能力。

建议职责：

1. 根据配置决定是否暴露 Swagger UI
2. 返回 Swagger UI 配置列表
3. 提供统一的 Swagger UI 页面入口

建议接口：

```go
type SwaggerProvider interface {
    Enabled() bool
    UIEnabled() bool
    Docs() []SwaggerDocConfig
}
```

建议路由：

- `/swagger`
- `/swagger/index.html`
- `/swagger/swagger-config`

实现方式建议：

1. 使用现成 Swagger UI 静态资源
2. 通过 `/swagger/swagger-config` 返回文档列表
3. Swagger UI 前端页面加载后，根据配置动态展示多个服务文档

建议直接使用：

- `swagger-ui-dist`

不建议在网关中自己实现 OpenAPI 渲染页面。

## 12. 用户认证解析设计

网关用户认证建议只做轻量解析，不做复杂业务校验。

建议定义：

```go
type UserIdentity struct {
    UserID      int64
    Username    string
    Roles       []string
    Permissions []string
}
```

如果使用统一 token store，可在网关通过：

- Bearer Token
- 共享 Redis
- 或认证中心 introspection

获取用户信息。

## 13. 反向代理设计

建议定义接口：

```go
type Proxy interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request, instance string, rule *ServiceRule) error
}
```

建议基于标准库 `httputil.ReverseProxy` 实现，并统一处理：

- 目标实例地址拼接
- Header 注入
- 转发超时
- 错误响应
- 响应状态透传
- 是否透传访问令牌

注入 Header 建议包括：

- `X-Request-Id`
- `X-Client-Id`
- `X-User-Id`
- `X-Username`

其中：

- `PassAccessToken == true` 时透传 `Authorization`
- `PassAccessToken == false` 时主动移除 `Authorization`

## 14. Consul 发现接口建议

建议复用 `aevons-framework/core/consul` 的发现能力，并在网关本地封装一层：

```go
type ServiceDiscovery interface {
    Discover(serviceName string) ([]Instance, error)
    Register(serviceName string, host string, port int) error
}
```

网关本地只关心：

- 服务名
- 健康实例列表
- 负载选择策略
- 选中的目标实例地址

不在网关里重复实现 Consul 协议细节，也不在配置中写死下游地址。

## 15. 负载选择策略

第一阶段建议只实现：

1. 随机
2. 轮询

建议定义接口：

```go
type LoadBalancer interface {
    Select(service string, instances []Instance) (*Instance, error)
}
```

后续可以扩展：

- 权重
- 一致性哈希
- 基于客户端会话粘性

## 16. 中间件执行顺序

建议中间件顺序如下：

1. `Recovery`
2. `RequestID`
3. `AccessLog`
4. `ClientResolve`
5. `ServiceMatch`
6. `ClientResourceCheck`
7. `ExcludeAuthCheck`
8. `UserAuth`
9. `RateLimit`
10. `ProxyDispatch`

说明：

- `ClientResourceCheck` 应在 `UserAuth` 前执行，减少无效鉴权开销
- `ServiceMatch` 应先命中目标服务，再决定是否需要排除认证
- `ExcludeAuthCheck` 用于判断当前请求是否命中服务级排除路由
- `ProxyDispatch` 作为最后一步

## 17. HTTP 上下文约定

建议在请求上下文中统一注入：

```go
type RequestContext struct {
    RequestID string
    Service   *ServiceRule
    Client    *ClientIdentity
    User      *UserIdentity
}
```

建议通过 context key 封装读取函数，不直接散落 `context.WithValue`。

## 18. 错误响应建议

统一返回结构建议为：

```json
{
  "code": "GATEWAY_FORBIDDEN",
  "message": "client has no access to this resource",
  "requestId": "xxx"
}
```

常见错误码建议：

- `GATEWAY_ROUTE_NOT_FOUND`
- `GATEWAY_CLIENT_INVALID`
- `GATEWAY_CLIENT_FORBIDDEN`
- `GATEWAY_TOKEN_INVALID`
- `GATEWAY_SERVICE_UNAVAILABLE`
- `GATEWAY_PROXY_ERROR`

## 19. 审计与观测建议

请求日志建议至少记录：

- 请求方法
- 请求路径
- 路由 ID
- 客户端 ID
- 用户 ID
- 目标服务
- 命中实例
- 响应状态码
- 耗时
- 是否命中 `ALL`

同时建议后续对接：

- Prometheus 指标
- OpenTelemetry Trace
- 审计日志服务

## 20. 第一阶段实现顺序

建议编码顺序如下：

1. 定义服务级配置结构
2. 实现服务匹配器
3. 实现排除认证路由匹配
4. 实现 Swagger UI 聚合模块
5. 实现客户端资源规则加载与匹配
6. 实现 Consul 注册与发现封装
7. 实现负载策略选择器
8. 实现反向代理
9. 实现用户认证解析
10. 按顺序装配中间件
11. 增加审计与统一错误响应

## 21. 结论

这版实现设计的核心思想是：

- 先把最关键、最稳定的骨架搭起来
- 采用服务维度配置，不把路由拆碎到每个接口
- 不提前引入过重的灰度、规则引擎、动态编排复杂度
- 优先把“客户端资源访问控制 + Consul 注册发现 + 服务级转发”做扎实

其中客户端资源访问控制建议作为网关第一批核心能力直接实现，`ALL` 语义也应在实现阶段就固定下来，避免后续解释漂移。
