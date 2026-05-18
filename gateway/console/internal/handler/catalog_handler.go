package handler

import (
	"io"
	"net/http"
	"time"

	consoleconfig "gateway-console/internal/config"
	"gateway-console/internal/dto"
	"gateway-console/internal/model"
	"gateway-console/internal/service"

	frameworkresp "github.com/calmlax/aevons-framework/response"
	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	catalog     *service.CatalogService
	publish     *service.PublishService
	serviceName string
	consoleCfg  consoleconfig.Settings
	httpClient  *http.Client
}

func NewCatalogHandler(catalog *service.CatalogService, publish *service.PublishService, serviceName string, consoleCfg consoleconfig.Settings) *CatalogHandler {
	return &CatalogHandler{
		catalog:     catalog,
		publish:     publish,
		serviceName: serviceName,
		consoleCfg:  consoleCfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h *CatalogHandler) Overview(c *gin.Context) {
	frameworkresp.Success(c, dto.OverviewResponse{
		Message: "gateway console overview",
		Data:    h.catalog.Overview(),
	})
}

func (h *CatalogHandler) Routes(c *gin.Context) {
	frameworkresp.Success(c, h.catalog.Routes())
}

func (h *CatalogHandler) Upstreams(c *gin.Context) {
	frameworkresp.Success(c, h.catalog.Upstreams())
}

func (h *CatalogHandler) Consumers(c *gin.Context) {
	frameworkresp.Success(c, h.catalog.Consumers())
}

func (h *CatalogHandler) Plugins(c *gin.Context) {
	frameworkresp.Success(c, h.catalog.Plugins())
}

func (h *CatalogHandler) Policies(c *gin.Context) {
	frameworkresp.Success(c, h.catalog.Policies())
}

func (h *CatalogHandler) SwaggerSources(c *gin.Context) {
	sources := h.catalog.SwaggerSources()
	for i := range sources {
		sources[i].ProxyURL = "/api/v1/gateway/swagger/" + sources[i].Service + "/swagger.json"
	}

	frameworkresp.Success(c, gin.H{
		"ui_url":  h.consoleCfg.SwaggerUIURL,
		"sources": sources,
	})
}

func (h *CatalogHandler) ProxySwagger(c *gin.Context) {
	serviceName := c.Param("service")
	source, ok := h.lookupSwaggerSource(serviceName)
	if !ok {
		frameworkresp.Fail(c, http.StatusNotFound, http.StatusNotFound, "gateway.swagger_source_not_found", map[string]any{
			"service": serviceName,
		})
		return
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, source.TargetURL, nil)
	if err != nil {
		frameworkresp.FailServerError(c, "gateway.swagger_proxy_failed", map[string]any{"error": err.Error()})
		return
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		frameworkresp.Fail(c, http.StatusBadGateway, http.StatusBadGateway, "gateway.swagger_proxy_failed", map[string]any{
			"service": serviceName,
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		frameworkresp.FailServerError(c, "gateway.swagger_proxy_failed", map[string]any{"error": err.Error()})
		return
	}

	if resp.StatusCode >= 300 {
		c.Data(resp.StatusCode, "application/json; charset=utf-8", body)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json; charset=utf-8"
	}
	c.Data(resp.StatusCode, contentType, body)
}

func (h *CatalogHandler) PublishPlan(c *gin.Context) {
	frameworkresp.Success(c, gin.H{
		"message": "gateway publish plan",
		"data":    h.publish.Plan(),
	})
}

func (h *CatalogHandler) PublishSnapshot(c *gin.Context) {
	snapshot, err := h.publish.Snapshot()
	if err != nil {
		frameworkresp.Fail(c, http.StatusBadGateway, http.StatusBadGateway, "gateway.snapshot_failed", map[string]any{
			"error": err.Error(),
		})
		return
	}

	frameworkresp.Success(c, gin.H{
		"message": "apisix publish snapshot",
		"data":    snapshot,
	})
}

func (h *CatalogHandler) PublishToAPISIX(c *gin.Context) {
	if err := h.publish.Publish(c.Request.Context()); err != nil {
		frameworkresp.Fail(c, http.StatusBadGateway, http.StatusBadGateway, "gateway.publish_failed", map[string]any{
			"error": err.Error(),
		})
		return
	}

	frameworkresp.Success(c, gin.H{
		"message": "publish to apisix succeeded",
	})
}

func (h *CatalogHandler) ControlPlaneHealth(c *gin.Context) {
	frameworkresp.Success(c, gin.H{
		"status": "ok",
		"name":   h.serviceName,
	})
}

func (h *CatalogHandler) lookupSwaggerSource(serviceName string) (source model.SwaggerSource, ok bool) {
	for _, source := range h.catalog.SwaggerSources() {
		if source.Service == serviceName {
			return source, true
		}
	}
	return model.SwaggerSource{}, false
}
