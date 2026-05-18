# Aevons Gateway Service

`gateway-service` 是 Aevons 自研 HTTP 网关服务，负责统一接入、认证校验、客户端资源控制、Consul 服务发现转发，以及多服务 Swagger 聚合展示。

## 目录

```text
gateway-service/
├── cmd/server
├── configs
├── internal
│   ├── auth
│   ├── clientauth
│   ├── config
│   ├── discovery
│   ├── model
│   ├── proxy
│   ├── router
│   └── swagger
├── ui/swagger
├── gateway-architecture.md
├── gateway-implementation.md
└── README.md
```

## 已实现能力

- 服务前缀路由匹配
- Consul 服务发现
- 基于 `auth-service` 的 Bearer Token 校验
- 客户端资源白名单控制，支持 `ALL`、精确路径、`/**` 前缀和 `ANY:/path/**`
- 反向代理转发
- Swagger 源聚合与 Swagger UI 页面
- 网关请求上下文统一注入
- 健康检查、统一错误码和审计访问日志
- 默认已接入 `auth-service`、`sys-service`、`log-service`、`gen-service`、`job-service`

## 主要入口

- 服务地址：`http://127.0.0.1:11080`
- 健康检查：`GET /health`
- Swagger UI：`GET /swagger/`
- Swagger 源列表：`GET /api/v1/gateway/swagger/sources`
- Swagger JSON 代理：`GET /api/v1/gateway/swagger/:service/swagger.json`

## 配置

配置文件在 [configs/config.yaml](./configs/config.yaml)，主要分为：

- `server`
- `consul`
- `log`
- `cors`
- `xss`
- `gateway`
- `swagger`
- `services`
- `clients`

其中：

- `services` 定义服务前缀、发现方式、负载策略、公开接口白名单
- `clients` 定义客户端资源访问规则
- `swagger.allowed_ips` 控制 Swagger UI 和聚合接口访问来源
- `swagger.docs` 定义聚合展示的文档源

## 启动

```bash
cd aevons-cloud/gateway-service
go run ./cmd/server --config configs
```

如需指定环境配置：

```bash
APP_ENV=development go run ./cmd/server --config configs
```

## 校验

```bash
GOCACHE=/tmp/go-build-cache go test ./...
```
