# Aevons Gateway

`gateway/` 用于承载 Aevons 企业级网关方案的落地实现。

当前网关方案采用：

```text
APISIX + Aevons Plugin System + Gateway Console
```

目标不是继续自研数据面代理内核，而是把以下三层稳定落地：

- `APISIX`：数据面网关，负责流量接入、路由、转发、插件执行
- `Aevons Plugins`：企业治理插件位点，承载鉴权、审计、租户、风控等能力
- `Gateway Console`：控制面，负责查看、聚合、发布和后续治理扩展

相关设计文档：

- [gateway-architecture.md](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/gateway-architecture.md:1)
- [gateway-implementation.md](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/gateway-implementation.md:1)

**目录说明**

```text
gateway/
├── apisix/
│   ├── conf/
│   │   ├── config.yaml
│   │   └── apisix.yaml
│   ├── plugins/
│   └── docker-compose.yaml
├── console/
│   ├── cmd/server/
│   ├── configs/
│   ├── handler/
│   ├── repository/
│   ├── service/
│   ├── internal/
│   └── ui/
├── scripts/
├── gateway-architecture.md
└── gateway-implementation.md
```

主要目录职责：

- `apisix/`
  - APISIX 容器部署文件
  - APISIX 启动配置
  - Aevons Lua 插件目录
- `apisix/conf/config.yaml`
  - APISIX 基础配置
  - 当前为 `traditional + yaml` 模式
- `apisix/conf/apisix.yaml`
  - 当前用于本地开发的静态路由示例
  - 现阶段仍包含业务 route/upstream 兜底配置
- `apisix/plugins/`
  - 企业插件实现位置
  - 当前已包含 `client-resource-auth`、`jwt-enterprise-auth`、`audit-log`、`tenant-isolation`、`ai-risk-control`
- `console/`
  - Gateway Console 控制面
  - 当前已支持目录查看、Swagger 聚合、APISIX 发布快照、发布入口
- `console/internal/apisixadmin/`
  - APISIX Admin API 客户端
- `console/internal/config/`
  - Console 自身扩展配置解析
- `console/internal/router/`
  - Console 路由装配
- `console/ui/swagger/`
  - 基于 `swagger-ui-dist` 的静态聚合页面
- `scripts/bootstrap.sh`
  - 本地启动提示脚本

**当前能力**

当前已经落地的能力：

- APISIX 容器部署骨架
- 三个服务的网关路由接入：
  - `auth-service`
  - `system-service`
  - `log-service`
- Gateway Console 启动与运行
- Swagger/OpenAPI 聚合页面
- 基于 Consul 的 upstream 发现与解析
- APISIX 发布快照与发布接口

当前仍处于第一阶段的内容：

- APISIX Lua 企业插件多为占位实现
- Console 仍以静态仓储数据为主，尚未演进为完整持久化管理面
- APISIX 当前仍使用本地 YAML 兜底路由
- APISIX 与 Console 的“控制面完全接管数据面”模式还在继续收口

**运行依赖**

启动网关前建议准备以下依赖：

- Docker / Docker Compose
- Go `1.25.0+`
- Consul
- 已注册到 Consul 的业务服务：
  - `auth-service`
  - `system-service`
  - `log-service`

如果要完整体验 Swagger 聚合和 upstream 发现，建议先启动：

1. Consul
2. `auth-service`
3. `system-service`
4. `log-service`
5. `gateway-console`
6. APISIX

**配置说明**

**1. APISIX**

文件：
[apisix/conf/config.yaml](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/apisix/conf/config.yaml:1)

当前关键配置：

- `deployment.role: traditional`
- `deployment.role_traditional.config_provider: yaml`
- 开启 `admin_key`
- 注册 Aevons 自定义插件列表

文件：
[apisix/conf/apisix.yaml](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/apisix/conf/apisix.yaml:1)

说明：

- 当前文件仍保留本地开发阶段的静态 route/upstream
- 长期目标是让业务 Route/Upstream 由 Gateway Console 通过 APISIX Admin API 下发
- 其中 `nodes` 不应被视为权威服务配置，最终应以 `service_name + consul discovery` 为准

**2. Gateway Console**

文件：
[console/configs/config.yaml](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/console/configs/config.yaml:1)

关键配置项：

- `server.port`
  - Console HTTP 端口，当前默认 `10900`
- `consul.address`
  - Consul 地址
- `console.apisix_admin_url`
  - APISIX Admin API 地址，默认 `http://127.0.0.1:9180`
- `console.apisix_admin_key`
  - APISIX Admin API Key
- `console.swagger_ui_url`
  - Swagger UI 外部配置展示字段

**部署说明**

**1. APISIX 部署**

文件：
[apisix/docker-compose.yaml](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/apisix/docker-compose.yaml:1)

当前容器映射端口：

- `9080`
  - APISIX 对外流量入口
- `9180`
  - APISIX Admin API

当前本地 Linux 环境采用 `host network` 启动 APISIX：

- APISIX 直接复用宿主机网络
- Consul discovery 使用 `http://127.0.0.1:8500`
- 业务服务和 Consul 都按宿主机地址暴露即可被 APISIX 访问

启动命令：

```bash
cd /home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/apisix
docker compose up -d
```

注意：

- 当前环境必须具备 Docker Socket 权限，否则容器无法启动
- 现有 `docker-compose.yaml` 保留了 `version` 字段，Docker 会提示该字段已过时，但不影响解析

**2. Gateway Console 部署**

启动命令：

```bash
cd /home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/console
go run ./cmd/server
```

当前 Console 已按 `aevons-framework` 的现行方式完成适配：

- 显式加载框架配置
- 使用 `core.NewApp`
- 使用 Gin 路由装配
- 支持 Consul 注册与注销

**启动说明**

可参考脚本：
[scripts/bootstrap.sh](/home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/scripts/bootstrap.sh:1)

推荐本地启动顺序：

```bash
# 1. 启动 APISIX
cd /home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/apisix
docker compose up -d

# 2. 启动 Gateway Console
cd /home/yhj/Desktop/data/mydata/aevons-cloud-dev/aevons-cloud/gateway/console
go run ./cmd/server
```

**访问入口**

Gateway Console 默认地址：

- 健康检查：
  - `http://127.0.0.1:10900/health`
- 网关概览：
  - `http://127.0.0.1:10900/api/v1/gateway/overview`
- Upstream 查看：
  - `http://127.0.0.1:10900/api/v1/gateway/upstreams`
- Swagger 聚合页面：
  - `http://127.0.0.1:10900/swagger/`
- Swagger 源列表：
  - `http://127.0.0.1:10900/api/v1/gateway/swagger/sources`

APISIX 默认地址：

- 流量入口：
  - `http://127.0.0.1:9080`
- Admin API：
  - `http://127.0.0.1:9180`

**发布说明**

当前 Gateway Console 已支持向 APISIX 下发快照：

- 发布计划：
  - `GET /api/v1/gateway/publish/plan`
- 发布快照：
  - `GET /api/v1/gateway/publish/snapshot`
- 执行发布：
  - `POST /api/v1/gateway/publish/run`

当前发布逻辑特性：

- route / consumer / policy 来自 Console 当前目录模型
- upstream 优先通过 Consul 发现健康实例
- `static_nodes_fallback` 仅作为兜底，不作为权威服务定义

**Swagger 使用说明**

当前 `gateway-console` 已集成 `swagger-ui-dist`，支持在一个页面中切换查看多个服务的 API 文档。

已接入的文档源：

- `auth-service`
- `system-service`
- `log-service`

聚合方式：

- Console 提供 Swagger 源清单接口
- Console 代理每个服务自己的 `/api/swagger.json`
- Swagger UI 页面通过下拉框切换不同服务文档

前提：

- 对应服务必须已启动
- 对应服务自己的 `api/swagger.json` 必须可访问

**已知限制**

- APISIX 当前尚未完全切到“仅 Console 下发”的控制模式
- `apisix/conf/apisix.yaml` 仍保留本地静态业务路由
- Aevons 企业插件目前多数仍为第一阶段占位实现
- 若 Consul 中没有健康实例，Console 发布会失败或退回兜底节点
- 若当前用户没有 Docker 权限，APISIX 容器无法启动

**后续建议**

建议按以下顺序继续演进：

1. 让 APISIX 从静态 YAML 路由逐步切换到完全由 Console 下发
2. 优先完成：
   - `jwt-enterprise-auth`
   - `client-resource-auth`
   - `audit-log`
3. 将 Console 从静态仓储升级为可持久化管理面
4. 补齐发布审计、回滚和可观测链路
