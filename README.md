# Aevons Cloud

`aevons-cloud` 是基于 `aevons-framework` 组织的一组微服务。

当前主要服务包括：

- `auth-service`
- `sys-service`
- `log-service`
- `gen-service`
- `job-service`
- `gateway/console`

**统一启动约定**

普通 HTTP 服务现在统一使用：

```go
app, err := core.Bootstrap()
engine, err := router.Setup(app)
err = core.RunGin(app, engine)
```

适用服务：

- `auth-service`
- `sys-service`
- `job-service`

带额外运行时能力的服务按“Bootstrap + 少量服务特有装配 + 统一运行时”模式处理：

- `gen-service`
  - 在 HTTP 服务之外还保留代码生成命令入口
- `log-service`
  - 只保留 gRPC 服务注册
  - 使用 `core.RunGinAndGRPC(...)` 接管 HTTP + gRPC 生命周期
- `gateway/console`
  - 使用 `core.BootstrapWithOptions(...)`
  - 不初始化 DB / Redis

**OpenAPI**

各服务生成 `api/swagger.json` 的命令：

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -d cmd/server,internal --parseDependency --parseInternal -o ./api --outputTypes json
```
