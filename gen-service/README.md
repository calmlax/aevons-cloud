# System service

## OpenAPI

生成 `api/swagger.json`：

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -d cmd/server,hander,dto,model,router,service,repository --parseDependency --parseInternal -o ./api --outputTypes json
```
