package handler

import (
	"github.com/armiariyan/assessment-tsel/internal/infrastructure/container"

	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo, cnt *container.Container) {
	h := SetupHandler(cnt).Validate()

	e.GET("/", h.healthCheckHandler.HealthCheck)

	v1 := e.Group("/v1")
	{
		customers := v1.Group("/customers")
		{
			customers.GET("", h.customersHandler.GetListCustomers)
		}

		items := v1.Group("/items")
		{
			items.GET("", h.itemsHandler.GetListItems)
		}

		invoices := v1.Group("/invoices")
		{
			invoices.GET("", h.invoicesHandler.GetListInvoices)
			invoices.GET("/:invoiceID", h.invoicesHandler.GetDetailInvoice)
			invoices.POST("", h.invoicesHandler.CreateInvoice)
			invoices.PATCH("", h.invoicesHandler.EditInvoice)
			invoices.DELETE("/:invoiceID", h.invoicesHandler.DeleteDetailInvoice)
		}
	}
}
