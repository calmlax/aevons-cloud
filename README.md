# Aevons Cloud

`aevons-cloud` 是基于 `aevons-framework` 组织的一套微服务项目，当前包含认证、系统管理、日志、代码生成、定时任务、前端、统一网关和 gRPC 契约。

公共框架见：

- [aevons-framework/README.md](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-framework/README.md:1)

**项目结构**

```text
aevons-cloud/
├── auth-service/
├── sys-service/
├── log-service/
├── gen-service/
├── job-service/
├── gateway-service/
├── frontend/
├── frontend-demo/
├── internal-grpc/
└── sql/
```

**模块说明**

- [auth-service](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/auth-service:1)
  - 认证中心
  - 登录、刷新、登出、邮箱验证码、注册、密码重置、OAuth2 授权码、Passkey、用户路由、用户资料
- [sys-service](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/sys-service:1)
  - 系统管理中心
  - 用户、角色、菜单、部门、岗位、字典、配置、语言、多语言资源、OAuth Client 管理
- [log-service](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/log-service:1)
  - 日志中心
  - 登录日志、操作日志
  - 同时对外提供 gRPC 写入服务
- [gen-service](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gen-service:1)
  - 代码生成服务
  - 生成表、字段设计、预览、下载、数据库导入
- [job-service](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/job-service:1)
  - 定时任务服务
  - 任务管理、执行日志、任务调度
- [gateway-service](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway-service:1)
  - 统一业务网关
  - 路由匹配、客户端识别、OAuth Client 资源校验、Bearer Token 校验、Consul 服务发现、反向代理、Swagger 聚合
- [frontend](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/frontend:1)
  - 主前端管理台
  - Vite + Vue
- [frontend-demo](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/frontend-demo:1)
  - 示例前端
- [internal-grpc](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/internal-grpc:1)
  - 内部 gRPC 契约
  - 当前主要是 `log_grpc`
- [sql](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/sql:1)
  - SQL 脚本与初始化相关内容

**统一后端结构**

业务服务基本采用这套分层：

```text
cmd/
configs/
api/
internal/
├── dto/
├── handler/
├── model/
├── repository/
├── router/
└── service/
```

特殊点：

- `log-service/internal/grpcs`
  - gRPC 服务端实现
- `sys-service/internal/middleware`
  - 系统管理侧操作日志等局部中间件
- `gateway-service/internal`
  - 网关侧独立拆分：
    - `auth`
    - `clientauth`
    - `discovery`
    - `proxy`
    - `swagger`
    - `gwcontext`

**启动方式**

当前后端服务已经基本统一到 `aevons-framework/core/bootstrap.go` 的启动模型。

普通 HTTP 服务：

```go
app, err := core.Bootstrap()
engine, err := router.Setup(app)
err = core.RunGin(app, engine)
```

当前适用：

- `auth-service`
- `sys-service`
- `job-service`
- `gateway-service`

带额外运行时能力的服务：

- `log-service`
  - `Bootstrap + gRPC 注册 + RunGinAndGRPC`
- `gen-service`
  - 除 HTTP 服务外还保留代码生成命令入口

**端口约定**

当前服务默认端口：

- `auth-service`
  - HTTP `10701`
  - gRPC `10801`
- `sys-service`
  - HTTP `10702`
  - gRPC `10802`
- `log-service`
  - HTTP `10703`
  - gRPC `10803`
- `gen-service`
  - HTTP `10704`
- `job-service`
  - HTTP `10705`
- `gateway-service`
  - HTTP `11080`
- `frontend`
  - dev server `5173`

**网关**

现在项目里真正承接业务流量的是：

- [gateway-service](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway-service:1)

它承担：

- 统一入口
- 服务发现
- 客户端资源校验
- Token 校验
- Swagger 聚合
- 入口层 CORS / XSS

业务服务侧已去掉统一入口治理中间件：

- `CORS`
- `XSSMiddleware`

这些现在统一放到网关层处理。

Swagger 页面入口：

- `http://127.0.0.1:11080/swagger`

Swagger 源通过 `service_id + Consul + path` 解析，不再写死实例 IP。

**前端**

[frontend](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/frontend:1) 当前开发时：

- 默认端口 `5173`
- 接口走 `VITE_API_PROXY_TARGET`
- 推荐代理到 `gateway-service`

当前常见开发链路：

```text
browser -> frontend(5173) -> gateway-service(11080) -> backend services
```

Passkey 登录已经和密码登录对齐：

- 前端会带 `client_id/client_secret`
- `auth-service` 会按 `passkey` grant type 校验客户端

**登录日志与操作日志**

日志写入分成两类：

- 业务访问日志
  - 由各服务自己的 `xlog` 输出到本地日志目录
- 业务行为日志
  - 登录日志 / 操作日志统一写入 `log-service`
  - 通过 `internal-grpc/log_grpc` 走 gRPC

目前：

- `auth-service` 会通过 gRPC 写登录日志
- `sys-service`、`gen-service`、`job-service` 会通过 gRPC 写操作日志
- `log-service` 负责落库

**OAuth Client 资源控制**

`gateway-service` 现在不再使用静态 fallback client 配置，唯一规则来源是：

- `sys_oauth_client.resources`

语义：

- `ALL`
- 或按服务名称逗号分割，例如：
  - `auth-service,sys-service,job-service`

同时：

- `gateway-service` 读 Redis 优先
- miss 后 DB 回源
- 有防击穿 / 防雪崩逻辑
- `sys-service` 提供手动刷新缓存入口

对应刷新接口：

- `POST /api/sys/v1/oauth/client/refresh-cache`

**OpenAPI / Swagger**

每个后端服务都有自己的：

- `api/swagger.json`

各服务 handler 里的 `@Router` 现在都已统一补成带 `/api/v1` 的完整路径。

如果重新生成某个服务的 Swagger JSON，可以使用：

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -d cmd/server,internal --parseDependency --parseInternal -o ./api --outputTypes json
```

通常在服务目录执行，例如：

```bash
cd auth-service
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -d cmd/server,internal --parseDependency --parseInternal -o ./api --outputTypes json
```

**开发建议**

推荐启动顺序：

1. MySQL / Redis / Consul
2. `log-service`
3. `auth-service`
4. `sys-service`
5. `gen-service`
6. `job-service`
7. `gateway-service`
8. `frontend`

这样能减少：

- gRPC 日志客户端启动早于 `log-service`
- 网关服务发现早于业务服务注册
- 前端代理目标未就绪

**当前状态**

这套项目已经不是纯脚手架状态，当前真实落地能力包括：

- 统一 Bootstrap 启动模型
- Consul 注册发现
- gRPC 日志中心写入
- Passkey 登录
- OAuth Client 资源治理
- 统一业务网关
- Swagger 聚合
- 前后端联调链路

如果继续推进，最适合往下收口的是：

- 自动刷新 OAuth Client 资源缓存
- 文档生成与启动脚本进一步统一
- 更细的观测、指标和审计能力
