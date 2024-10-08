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
		products := v1.Group("/products")
		{
			products.GET("", h.productsHandler.GetListProducts)
			products.GET("/:id", h.productsHandler.GetDetailProduct)
			products.POST("", h.productsHandler.CreateProduct)
			products.PATCH("", h.productsHandler.UpdateProduct)
			products.DELETE("/:id", h.productsHandler.DeleteProduct)
		}
	}
}
