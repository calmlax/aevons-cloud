package handler

import (
	"net/http"

	"aevons-cloud/gateway/console/dto"
	"aevons-cloud/gateway/console/service"

	frameworkresp "github.com/calmlax/aevons-framework/response"
	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	catalog     *service.CatalogService
	publish     *service.PublishService
	serviceName string
}

func NewCatalogHandler(catalog *service.CatalogService, publish *service.PublishService, serviceName string) *CatalogHandler {
	return &CatalogHandler{
		catalog:     catalog,
		publish:     publish,
		serviceName: serviceName,
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

func (h *CatalogHandler) PublishPlan(c *gin.Context) {
	frameworkresp.Success(c, gin.H{
		"message": "gateway publish plan",
		"data":    h.publish.Plan(),
	})
}

func (h *CatalogHandler) PublishSnapshot(c *gin.Context) {
	frameworkresp.Success(c, gin.H{
		"message": "apisix publish snapshot",
		"data":    h.publish.Snapshot(),
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
