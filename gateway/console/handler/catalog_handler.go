package handler

import (
	"encoding/json"
	"net/http"

	"aevons-cloud/gateway/console/dto"
	"aevons-cloud/gateway/console/service"
)

type CatalogHandler struct {
	catalog *service.CatalogService
	publish *service.PublishService
}

func NewCatalogHandler(catalog *service.CatalogService, publish *service.PublishService) *CatalogHandler {
	return &CatalogHandler{
		catalog: catalog,
		publish: publish,
	}
}

func (h *CatalogHandler) Overview(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, dto.OverviewResponse{
		Message: "gateway console overview",
		Data:    h.catalog.Overview(),
	})
}

func (h *CatalogHandler) Routes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.catalog.Routes()})
}

func (h *CatalogHandler) Upstreams(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.catalog.Upstreams()})
}

func (h *CatalogHandler) Consumers(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.catalog.Consumers()})
}

func (h *CatalogHandler) Plugins(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.catalog.Plugins()})
}

func (h *CatalogHandler) Policies(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.catalog.Policies()})
}

func (h *CatalogHandler) PublishPlan(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "gateway publish plan",
		"data":    h.publish.Plan(),
	})
}

func (h *CatalogHandler) PublishSnapshot(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "apisix publish snapshot",
		"data":    h.publish.Snapshot(),
	})
}

func (h *CatalogHandler) PublishToAPISIX(w http.ResponseWriter, r *http.Request) {
	if err := h.publish.Publish(r.Context()); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"message": "publish to apisix failed",
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "publish to apisix succeeded",
	})
}

func (h *CatalogHandler) ControlPlaneHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"name":   "gateway-console",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
