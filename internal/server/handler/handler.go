package handler

import (
	"github.com/armiariyan/assessment-tsel/internal/infrastructure/container"
)

type Handler struct {
	healthCheckHandler *healthCheckHandler
	productsHandler    *productsHandler
}

func SetupHandler(container *container.Container) *Handler {
	return &Handler{
		healthCheckHandler: NewHealthCheckHandler().SetHealthCheckService(container.HealthCheckService).Validate(),
		productsHandler:    NewProductsHandler().SetProductsService(container.ProductService).Validate(),
	}
}

func (h *Handler) Validate() *Handler {
	if h.healthCheckHandler == nil {
		panic("healthCheckHandler is nil")
	}
	return h
}
