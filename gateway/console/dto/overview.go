package dto

import "aevons-cloud/gateway/console/internal/model"

type OverviewResponse struct {
	Message string         `json:"message"`
	Data    model.Overview `json:"data"`
}
