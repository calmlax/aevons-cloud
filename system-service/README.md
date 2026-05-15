# System service

## OpenAPI

生成 `api/swagger.json`：

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o ./api --outputTypes json
```
