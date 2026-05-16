package main

import (
	"context"
	"log"

	"aevons-cloud/gateway/console/internal/router"

	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/xjson"
	"github.com/gin-gonic/gin/binding"
)

func init() {
	binding.EnableDecoderUseNumber = true
	xjson.InitGin()
}

func main() {
	application, err := core.New("gateway-console", "./configs", core.WithRoutes(router.New()))
	if err != nil {
		log.Fatalf("bootstrap gateway-console: %v", err)
	}

	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("run gateway-console: %v", err)
	}

	core.WaitForShutdown(application)
}
