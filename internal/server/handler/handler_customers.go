package handler

import (
	"net/http"

	"github.com/armiariyan/assessment-tsel/internal/pkg/utils"
	"github.com/armiariyan/assessment-tsel/internal/usecase/customers"

	"github.com/labstack/echo/v4"
)

type customersHandler struct {
	customersService customers.CustomersService
}

func NewCustomersHandler() *customersHandler {
	return &customersHandler{}
}

func (h *customersHandler) SetCustomersService(service customers.CustomersService) *customersHandler {
	h.customersService = service
	return h
}

func (h *customersHandler) Validate() *customersHandler {
	if h.customersService == nil {
		panic("customersService is nil")
	}

	return h
}

func (h *customersHandler) GetListCustomers(c echo.Context) (err error) {
	ctx := c.Request().Context()

	req := customers.GetListCustomersRequest{}
	if err = utils.Validate(c, &req); err != nil {
		return
	}
	resp, err := h.customersService.GetListCustomers(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, resp)

}
