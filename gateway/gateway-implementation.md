# Aevons Gateway 企业级网关实现设计文档（APISIX 方案）

本文档用于承接 [gateway-architecture.md](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/gateway-architecture.md:1) 的企业级网关架构设计，重点说明在 `APISIX + Aevons Plugin System + Gateway Console` 路线下，系统应如何落地实现。

本文档不再沿用“自研 Go 网关内核”的实现思路，而是明确以 APISIX 作为企业级网关数据面内核。

## 1. 实现目标

企业级阶段的实现目标是：

1. 使用 APISIX 作为统一高性能网关入口
2. 基于 Consul 实现服务发现与健康实例转发
3. 基于 APISIX 插件机制承载 Aevons 企业治理能力
4. 通过独立 Gateway Console 管理路由、插件、消费者与策略
5. 建立“控制面 + 数据面 + 插件层”的长期演进基础

第一阶段不再重点实现：

- 自研 HTTP 代理内核
- 自研通用路由引擎
- 自研长连接管理
- 自研高性能插件执行框架

## 2. 总体实现结构

企业级实现建议拆为三块：

```text
gateway/
├── apisix/
│   ├── docker-compose.yaml
│   ├── conf/
│   └── plugins/
├── console/
│   ├── api/
│   ├── service/
│   ├── repository/
│   └── model/
├── docs/
│   ├── gateway-architecture.md
│   └── gateway-implementation.md
└── scripts/
```

说明：

- `apisix/`
  - APISIX 节点部署、配置文件、插件目录
- `console/`
  - Aevons Gateway Console
  - 负责路由配置、插件配置、消费者配置、策略下发
- `docs/`
  - 架构与实现文档
- `scripts/`
  - 初始化、发布、同步配置、导入路由脚本

## 3. 运行时拓扑

```text
Client
  -> APISIX Data Plane
      -> APISIX Route / Upstream / Plugin Chain
          -> auth-service
          -> system-service
          -> log-service

Gateway Console
  -> APISIX Admin API
  -> Consul
  -> MySQL / Redis
```

职责划分：

- `APISIX`
  - 流量入口
  - 路由匹配
  - 上游转发
  - 插件执行
  - 基础流量治理

- `Gateway Console`
  - 路由管理
  - 插件管理
  - 消费者管理
  - 客户端策略管理
  - 发布与审计

## 4. 目录与模块职责

### 4.1 APISIX 侧

建议目录：

```text
apisix/
├── docker-compose.yaml
├── conf/
│   ├── config.yaml
│   └── apisix.yaml
├── plugins/
│   ├── client-resource-auth.lua
│   ├── jwt-enterprise-auth.lua
│   ├── tenant-isolation.lua
│   ├── audit-log.lua
│   └── ai-risk-control.lua
└── upstream/
```

### 4.2 Console 侧

建议目录：

```text
console/
├── cmd/server/
│   └── main.go
├── configs/
│   └── config.yaml
├── handler/
├── service/
├── repository/
├── model/
├── dto/
└── internal/
    ├── apisixadmin/
    ├── route/
    ├── plugin/
    ├── consumer/
    ├── policy/
    └── publish/
```

## 5. APISIX 需要管理的核心对象

企业级落地时，核心不是“自己写转发代码”，而是管理 APISIX 的这些对象：

### 5.1 Route

负责：

- 匹配请求路径
- 匹配方法
- 匹配 Header / Host / 版本
- 绑定插件链
- 绑定上游服务

### 5.2 Upstream

负责：

- 目标服务定义
- 负载均衡策略
- 健康检查
- 重试与超时
- 服务发现集成

### 5.3 Consumer

负责：

- 客户端身份建模
- API Key / JWT / 签名信息
- 客户端维度插件绑定

### 5.4 Plugin Config

负责：

- 插件参数模板
- 同类路由统一复用
- 环境差异化配置

### 5.5 Global Rule

负责：

- 全局安全策略
- 全局限流
- 全局 Header 注入
- 全局审计与追踪

## 6. Route + Predicate + Filter 在 APISIX 中的映射

企业级设计中的：

- `Route`
- `Predicate`
- `Filter`

在 APISIX 中可以映射为：

| 抽象 | APISIX 对应物 |
|---|---|
| Route | Route |
| Predicate | uri / methods / hosts / vars |
| Filter | plugins |

### 6.1 Route 示例

```json
{
  "id": "system-route",
  "uri": "/api/v1/system/*",
  "methods": ["GET", "POST", "PUT", "DELETE"],
  "plugins": {
    "jwt-enterprise-auth": {},
    "client-resource-auth": {}
  },
  "upstream_id": "system-service-upstream"
}
```

### 6.2 Predicate 映射

例如：

- `Path=/api/v1/system/**`
  - `uri`
- `Method=GET`
  - `methods`
- `Header=X-Version:v2`
  - `vars`
- `Host=admin.aevons.com`
  - `hosts`

### 6.3 Filter 映射

例如：

- `RewritePath`
- `JWTAuth`
- `RateLimit`
- `Retry`
- `CircuitBreaker`
- `AddHeader`

都通过 APISIX `plugins` 配置落地。

## 7. Consul 服务发现实现

### 7.1 目标

所有下游服务统一注册到 Consul，APISIX 不写死实例地址。

### 7.2 落地方式

推荐两种模式：

1. APISIX 通过服务发现机制直接对接 Consul
2. Gateway Console 监听 Consul 实例变化，并同步到 APISIX Upstream

对于 Aevons 当前阶段，更建议：

### 第一阶段

Gateway Console 负责把 Consul 中的健康实例同步到 APISIX Upstream

原因：

- 更容易和当前 `aevons-framework/core/consul` 保持一致
- 可以在 Console 中统一做服务审计和发布控制
- 更方便后续接入自定义负载策略与治理能力

### 7.3 Upstream 示例

```json
{
  "id": "system-service-upstream",
  "type": "roundrobin",
  "nodes": {
    "192.168.31.81:10702": 1,
    "192.168.31.82:10702": 1
  }
}
```

## 8. Aevons Plugin System 实现建议

### 8.1 插件实现位置

Aevons 企业插件建议放在：

```text
apisix/plugins/
```

以 APISIX Lua 插件形式实现。

### 8.2 插件注册方式

通过 APISIX 配置注册：

- 插件名
- 执行阶段
- 配置 schema
- 优先级

### 8.3 插件分层

建议按以下类型划分：

- 安全插件
- 治理插件
- 审计插件
- 多租户插件
- 流量插件

## 9. Aevons 核心企业插件实现清单

### 9.1 client-resource-auth

职责：

- 识别客户端身份
- 判断客户端状态
- 校验客户端资源权限
- 实现 `ALL` 语义
- 实现多客户端资源隔离

输入：

- JWT 中 `client_id`
- API Key
- mTLS 身份

输出：

- 是否放行
- 客户端上下文 Header 注入

### 9.2 jwt-enterprise-auth

职责：

- 校验 JWT
- 提取用户信息
- 注入 `X-User-Id` / `X-Username`
- 构建角色与权限上下文

### 9.3 tenant-isolation

职责：

- 识别 TenantID
- 注入租户上下文
- 控制租户访问边界

### 9.4 audit-log

职责：

- 审计请求
- 记录敏感操作
- 风险事件上报

### 9.5 ai-risk-control

职责：

- 异常流量分析
- 风险识别
- 自动封禁
- AI WAF 扩展

## 10. Gateway Console 实现职责

Gateway Console 不是转发数据面，而是控制面。

建议职责：

### 10.1 Route 管理

- 新增/修改/删除路由
- 路由发布
- 路由灰度配置

### 10.2 Upstream 管理

- 服务注册信息同步
- Upstream 节点维护
- 权重与健康状态可视化

### 10.3 Consumer 管理

- 客户端注册
- 客户端密钥管理
- 客户端状态启停

### 10.4 Plugin 管理

- 插件启停
- 插件参数配置
- 插件模板配置

### 10.5 Policy 管理

- 客户端资源策略
- 多租户策略
- 后续 Casbin / OPA 策略

### 10.6 发布与审计

- 配置变更审计
- 操作人记录
- 发布记录
- 回滚能力

## 11. 客户端资源控制的落地方式

在 APISIX 方案下，客户端资源控制不再放在自研 Gin Middleware 中，而是落在：

```text
client-resource-auth 插件
```

### 11.1 第一阶段策略来源

可由 Gateway Console 从：

- MySQL
- Redis 缓存
- 配置中心

加载后下发到插件配置。

### 11.2 `ALL` 语义约束

仍保持当前架构稿中的边界：

1. `ALL` 仅代表业务接口入口全放行
2. 不包含系统保留资源
3. 不绕过用户认证
4. 不绕过业务数据权限
5. 不代表超级管理员

### 11.3 企业级演进

后续建议接入：

- `Casbin`
- `OPA`

承载复杂 `Policy`。

## 12. Swagger / OpenAPI 聚合实现

### 12.1 服务侧

各微服务继续维护：

```text
/api/swagger.json
```

### 12.2 网关侧

APISIX 只负责代理：

- `/auth/api/swagger.json`
- `/system/api/swagger.json`
- `/log/api/swagger.json`

### 12.3 UI 聚合方式

建议单独部署：

```text
swagger-ui
```

或由 Gateway Console 提供 Swagger 文档聚合页面。

不建议：

- 让 APISIX 自己解析 OpenAPI
- 在数据面网关里硬塞文档渲染逻辑

## 13. 安全体系实现

### 13.1 客户端身份

禁止：

```text
X-Client-Id 裸传
```

可信来源：

- JWT
- API Key
- mTLS

### 13.2 用户身份

通过：

```text
Authorization: Bearer xxx
```

传递，并由 `jwt-enterprise-auth` 插件解析。

### 13.3 安全能力落地

APISIX + Aevons 插件实现：

- JWT Auth
- API Key
- WAF
- IP 黑白名单
- Anti Replay
- 风险控制

## 14. 可观测性实现

### 14.1 指标

Prometheus：

- QPS
- 请求延迟
- 状态码分布
- 限流次数
- 熔断次数
- 插件拒绝次数

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

## 15. 实施阶段建议

### Phase 1

- APISIX 基础部署
- Consul 服务发现接入
- 基础 Route / Upstream 管理
- Swagger 聚合代理

### Phase 2

- `client-resource-auth`
- `jwt-enterprise-auth`
- 审计与日志插件

### Phase 3

- 多租户插件
- 插件模板
- Console 路由与插件配置界面

### Phase 4

- Policy Engine
- Casbin / OPA 集成

### Phase 5

- AI 风控
- 智能限流
- 风险识别

### Phase 6

- 完整 Control Plane
- 多环境发布
- 插件市场化与治理平台化

## 16. 为什么不完全自研网关

因为企业级网关最难的不是请求转发，而是：

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
