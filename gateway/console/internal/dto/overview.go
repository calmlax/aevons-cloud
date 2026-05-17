package dto

import "gateway-console/internal/model"

type OverviewResponse struct {
	Message string         `json:"message"`
	Data    model.Overview `json:"data"`
}
