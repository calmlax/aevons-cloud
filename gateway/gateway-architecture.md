# Aevons Gateway 企业级网关架构设计文档（APISIX 方案）

## 1. 文档目标

本文档定义 Aevons 微服务体系在企业级阶段的统一网关架构方案。

目标：

- 不重复自研成熟网关底层能力
- 基于成熟网关快速企业化
- 保留 Aevons 自身核心治理能力
- 构建可演进的企业级 API Gateway Platform

最终采用：

```text
Apache APISIX + Aevons Gateway Plugin System + Aevons Gateway Console
```

模式。

## 2. 架构选型结论

经过对以下方案评估：

- 自研 Go Gateway
- Spring Cloud Gateway
- Envoy
- Kong
- Traefik
- Apache APISIX

最终选择：

### 企业级主网关：Apache APISIX

原因：

| 能力 | APISIX |
|---|---|
| 高性能 | 极强（基于 Nginx / OpenResty） |
| 动态配置 | 原生支持 |
| 热更新 | 原生支持 |
| 插件体系 | 完整 |
| 灰度发布 | 支持 |
| 熔断限流 | 支持 |
| OpenTelemetry | 支持 |
| Prometheus | 支持 |
| 多节点 | 支持 |
| Admin API | 完整 |
| Kubernetes | 支持 |
| Lua 插件扩展 | 强 |
| 企业成熟度 | 高 |

## 3. 总体架构

```text
                    ┌────────────────────┐
                    │   Browser / App    │
                    └─────────┬──────────┘
                              │
                    ┌─────────▼──────────┐
                    │       APISIX       │
                    │  Enterprise Gateway│
                    └─────────┬──────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
 ┌───────▼────────┐  ┌───────▼────────┐  ┌───────▼────────┐
 │ auth-service   │  │ system-service │  │ log-service    │
 └────────────────┘  └────────────────┘  └────────────────┘
```

## 4. 架构分层

整体架构分为：

| 层 | 职责 |
|---|---|
| Access Layer | APISIX 统一流量入口 |
| Governance Layer | 鉴权、限流、审计、灰度、WAF |
| Plugin Layer | Aevons 自定义企业插件 |
| Service Layer | Gin 微服务集群 |
| Discovery Layer | Consul 服务发现 |
| Control Layer | Aevons Gateway Console |

## 5. 技术栈

| 模块 | 技术 |
|---|---|
| 微服务 | Go + Gin |
| 服务发现 | Consul |
| 网关 | APISIX |
| 配置中心 | Consul KV / Nacos（后期） |
| 数据库 | MySQL |
| 缓存 | Redis |
| 链路追踪 | OpenTelemetry |
| 指标监控 | Prometheus + Grafana |
| 日志 | Loki / ELK |
| Swagger | swagger-ui + OpenAPI |
| 权限策略 | Casbin（后期） |

## 6. APISIX 在体系中的职责

APISIX 负责：

### 6.1 流量入口

统一承接：

- HTTP
- HTTPS
- WebSocket
- HTTP/2

入口流量。

### 6.2 动态路由

负责：

- Path 路由
- Header 路由
- Host 路由
- 灰度路由
- 权重路由

例如：

```text
/api/v1/auth/**
    -> auth-service

/api/v1/system/**
    -> system-service
```

### 6.3 服务发现

通过 Consul 动态发现服务实例，不再配置固定 IP。

### 6.4 流量治理

APISIX 原生提供：

- 限流
- 熔断
- 重试
- 超时
- 健康检查
- 负载均衡
- 灰度发布
- 流量镜像

### 6.5 插件运行时

APISIX Plugin Runtime 作为：

```text
Aevons Enterprise Gateway Extension Layer
```

承载 Aevons 企业治理插件。

## 7. Aevons 自定义插件体系

Aevons 不重复造网关内核。

真正自研的是：

### 企业治理能力插件

也就是说：

- APISIX 负责高性能代理内核与通用流量治理
- Aevons 负责安全、权限、审计、租户、风控等企业能力

## 8. 插件架构设计

### 8.1 插件执行模型

```text
Request
   ↓
Global Plugin
   ↓
Route Plugin
   ↓
Consumer Plugin
   ↓
Proxy
   ↓
Response Plugin
```

### 8.2 插件分类

| 类型 | 说明 |
|---|---|
| Security Plugin | 安全插件 |
| Governance Plugin | 治理插件 |
| Transform Plugin | 请求改写插件 |
| Observe Plugin | 观测插件 |
| Traffic Plugin | 流量治理插件 |

## 9. Aevons 核心企业插件

### 9.1 client-resource-auth

核心插件：

```text
客户端资源访问控制
```

负责：

- 客户端识别
- 客户端资源权限
- `ALL` 权限控制
- 客户端状态校验
- 多客户端隔离

这是 Aevons 网关的核心能力之一。

### 9.2 jwt-enterprise-auth

企业 JWT 鉴权插件。

负责：

- JWT 校验
- 用户上下文注入
- Token 解析
- 用户身份提取
- RBAC 上下文构建

### 9.3 tenant-isolation

多租户隔离插件。

负责：

- TenantID 注入
- 租户路由隔离
- 数据访问边界

### 9.4 audit-log

审计插件。

负责：

- 请求审计
- 风险操作记录
- 安全事件上报

### 9.5 ai-risk-control

AI 风控插件。

未来可扩展：

- 风险识别
- 异常流量分析
- 自动封禁
- AI WAF

## 10. Route + Predicate + Filter 模型

企业级阶段，不再停留在：

```text
仅 Path Prefix 模式
```

而是演进为：

```text
Route + Predicate + Filter
```

模型。

### 10.1 Route

```json
{
  "id": "system-route",
  "uri": "http://system-service",
  "predicates": [
    "Path=/api/v1/system/**"
  ],
  "filters": [
    "jwt-enterprise-auth",
    "client-resource-auth"
  ]
}
```

### 10.2 Predicate

支持：

| Predicate | 作用 |
|---|---|
| Path | 路径匹配 |
| Method | 方法匹配 |
| Header | Header 匹配 |
| Host | Host 匹配 |
| Weight | 权重匹配 |
| ClientID | 客户端匹配 |

### 10.3 Filter

支持：

| Filter | 作用 |
|---|---|
| RewritePath | 路径改写 |
| StripPrefix | 前缀移除 |
| JWTAuth | JWT 鉴权 |
| RateLimit | 限流 |
| Retry | 重试 |
| CircuitBreaker | 熔断 |
| AddHeader | Header 注入 |

## 11. Swagger / OpenAPI 方案

采用：

```text
各服务维护自己的 OpenAPI
+
统一 swagger-ui 聚合
```

模式。

### 11.1 服务侧

每个服务继续维护：

```text
/api/swagger.json
```

独立文档。

Go 侧推荐：

- `swaggo/swag`
- `gin-swagger`

### 11.2 网关侧

APISIX 只负责代理：

```text
/auth/api/swagger.json
/system/api/swagger.json
/log/api/swagger.json
```

不解析 OpenAPI 内容。

### 11.3 文档聚合

单独部署：

```text
swagger-ui
```

统一聚合：

```js
urls: [
  {
    name: "auth-service",
    url: "/auth/api/swagger.json"
  },
  {
    name: "system-service",
    url: "/system/api/swagger.json"
  }
]
```

## 12. 服务发现方案

采用：

### Consul

负责：

- 服务注册
- 服务发现
- 健康检查
- 第一阶段配置存储

## 13. 安全体系

### 13.1 客户端身份

禁止：

```text
X-Client-Id 裸传
```

可信客户端来源：

- JWT
- API Key
- mTLS

### 13.2 用户身份

用户身份通过：

```text
Authorization: Bearer xxx
```

传递。

### 13.3 安全能力

通过 APISIX + 自定义插件实现：

- JWT Auth
- API Key
- WAF
- IP 黑白名单
- Anti Replay
- 风险控制

## 14. 可观测性体系

### 14.1 指标

Prometheus 关注：

- QPS
- 延迟
- 状态码
- 限流次数
- 熔断次数

### 14.2 日志

统一：

- Request Log
- Audit Log
- Error Log

### 14.3 Trace

OpenTelemetry：

- `request_id`
- `trace_id`
- `span_id`

## 15. 企业级演进方向

后续继续演进：

| 阶段 | 能力 |
|---|---|
| Phase 1 | APISIX + 基础插件 |
| Phase 2 | 企业权限插件 |
| Phase 3 | 多租户体系 |
| Phase 4 | Policy Engine |
| Phase 5 | AI Gateway |
| Phase 6 | Control Plane |

## 16. 为什么不完全自研网关

因为企业级网关最难的不是：

```text
请求转发
```

而是：

- 动态配置
- 热更新
- 插件体系
- HTTP 内核
- 长连接
- 流量治理
- 可观测性
- 高并发
- 灰度
- 多节点同步

这些能力 APISIX 已经非常成熟。

Aevons 真正应该自研的是：

```text
企业治理能力
```

而不是：

```text
重复造代理内核
```

## 17. 最终结论

Aevons Gateway 最合理路线：

```text
APISIX
    +
Aevons Enterprise Plugin System
    +
Aevons Gateway Console
```

即：

- 使用 APISIX 作为企业级高性能网关内核
- 使用自定义插件承载企业治理能力
- 使用 Gateway Console 实现统一治理平台

最终形成：

```text
企业级 API Gateway Platform
```

而不是单纯：

```text
请求转发器
```
