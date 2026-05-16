package router

import (
	"net/http"

	"aevons-cloud/gateway/console/handler"
	"aevons-cloud/gateway/console/internal/apisixadmin"
	"aevons-cloud/gateway/console/repository"
	"aevons-cloud/gateway/console/service"

	"github.com/calmlax/aevons-framework/web"
)

type Routes struct {
	handler *handler.CatalogHandler
}

func New() web.RouteRegistrar {
	repo := repository.NewStaticRepository()
	catalogService := service.NewCatalogService(repo)
	apisixClient := apisixadmin.New("http://127.0.0.1:9180", "replace-with-real-admin-key")
	publishService := service.NewPublishService(catalogService, apisixClient)

	return &Routes{
		handler: handler.NewCatalogHandler(catalogService, publishService),
	}
}

func (r *Routes) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/gateway/overview", r.handler.Overview)
	mux.HandleFunc("/api/v1/gateway/routes", r.handler.Routes)
	mux.HandleFunc("/api/v1/gateway/upstreams", r.handler.Upstreams)
	mux.HandleFunc("/api/v1/gateway/consumers", r.handler.Consumers)
	mux.HandleFunc("/api/v1/gateway/plugins", r.handler.Plugins)
	mux.HandleFunc("/api/v1/gateway/policies", r.handler.Policies)
	mux.HandleFunc("/api/v1/gateway/publish/plan", r.handler.PublishPlan)
	mux.HandleFunc("/api/v1/gateway/publish/snapshot", r.handler.PublishSnapshot)
	mux.HandleFunc("/api/v1/gateway/publish/run", r.handler.PublishToAPISIX)
	mux.HandleFunc("/api/v1/gateway/healthz/control-plane", r.handler.ControlPlaneHealth)
}
