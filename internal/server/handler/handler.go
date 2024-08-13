package handler

import (
	"github.com/armiariyan/assessment-tsel/internal/infrastructure/container"
)

type Handler struct {
	healthCheckHandler *healthCheckHandler
	customersHandler   *customersHandler
	itemsHandler       *itemsHandler
	invoicesHandler    *invoicesHandler
}

func SetupHandler(container *container.Container) *Handler {
	return &Handler{
		healthCheckHandler: NewHealthCheckHandler().SetHealthCheckService(container.HealthCheckService).Validate(),
		customersHandler:   NewCustomersHandler().SetCustomersService(container.CustomersService).Validate(),
		itemsHandler:       NewItemsHandler().SetItemsService(container.ItemsService).Validate(),
		invoicesHandler:    NewInvoicesHandler().SetInvoicesService(container.InvoicesService).Validate(),
	}
}

func (h *Handler) Validate() *Handler {
	if h.healthCheckHandler == nil {
		panic("healthCheckHandler is nil")
	}
	return h
}
