package handler

import (
	"net/http"

	"github.com/armiariyan/assessment-tsel/internal/pkg/utils"
	"github.com/armiariyan/assessment-tsel/internal/usecase/products"

	"github.com/labstack/echo/v4"
)

type productsHandler struct {
	productsService products.Service
}

func NewProductsHandler() *productsHandler {
	return &productsHandler{}
}

func (h *productsHandler) SetProductsService(service products.Service) *productsHandler {
	h.productsService = service
	return h
}

func (h *productsHandler) Validate() *productsHandler {
	if h.productsService == nil {
		panic("productsService is nil")
	}

	return h
}

func (h *productsHandler) GetListProducts(c echo.Context) (err error) {
	ctx := c.Request().Context()

	req := products.GetListProductsRequest{}
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	resp, err := h.productsService.GetListProducts(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, resp)

}
