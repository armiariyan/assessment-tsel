package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/armiariyan/assessment-tsel/internal/pkg/log"
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

func (h *productsHandler) GetDetailProduct(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// * get and convert id to int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(ctx, "failed convert id from param into integer", err)
		c.Set("invalid-format", true)
		err = fmt.Errorf("invalid id")
		return
	}

	resp, _ := h.productsService.GetDetailProduct(ctx, uint(id))

	return c.JSON(http.StatusOK, resp)

}

func (h *productsHandler) CreateProduct(c echo.Context) (err error) {
	ctx := c.Request().Context()

	req := products.CreateProductRequest{}
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	resp, err := h.productsService.CreateProduct(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *productsHandler) UpdateProduct(c echo.Context) (err error) {
	ctx := c.Request().Context()

	req := products.UpdateProductRequest{}
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	resp, _ := h.productsService.UpdateProduct(ctx, req)

	return c.JSON(http.StatusOK, resp)
}

func (h *productsHandler) DeleteProduct(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// * get and convert id to int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(ctx, "failed convert id from param into integer", err)
		c.Set("invalid-format", true)
		err = fmt.Errorf("invalid id")
		return
	}

	resp, _ := h.productsService.DeleteProduct(ctx, uint(id))

	return c.JSON(http.StatusOK, resp)

}
