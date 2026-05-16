# Aevons Gateway

本目录用于承载 Aevons 企业级网关方案的落地实现。

当前实现按以下分层组织：

- `apisix/`：APISIX 数据面部署、配置与企业插件
- `console/`：Gateway Console 控制面
- `scripts/`：初始化与运维脚本
- `gateway-architecture.md`：架构设计文档
- `gateway-implementation.md`：实现设计文档

当前阶段的目标不是继续自研网关数据面，而是先把：

- APISIX 网关内核
- Aevons 企业插件位点
- Gateway Console 控制面骨架

这三层立起来。
