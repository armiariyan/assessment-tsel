package handler

import (
	"net/http"

	"github.com/armiariyan/assessment-tsel/internal/pkg/utils"
	"github.com/armiariyan/assessment-tsel/internal/usecase/invoices"

	"github.com/labstack/echo/v4"
)

type invoicesHandler struct {
	invoicesService invoices.InvoicesService
}

func NewInvoicesHandler() *invoicesHandler {
	return &invoicesHandler{}
}

func (h *invoicesHandler) SetInvoicesService(service invoices.InvoicesService) *invoicesHandler {
	h.invoicesService = service
	return h
}

func (h *invoicesHandler) Validate() *invoicesHandler {
	if h.invoicesService == nil {
		panic("invoicesService is nil")
	}

	return h
}

func (h *invoicesHandler) GetListInvoices(c echo.Context) (err error) {
	ctx := c.Request().Context()

	req := invoices.GetListInvoicesRequest{}
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	resp, err := h.invoicesService.GetListInvoices(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *invoicesHandler) GetDetailInvoice(c echo.Context) (err error) {
	ctx := c.Request().Context()

	uniqueInvoiceID := c.Param("invoiceID")

	resp, err := h.invoicesService.GetDetailInvoice(ctx, uniqueInvoiceID)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *invoicesHandler) CreateInvoice(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// bind and validate request
	req := invoices.CreateInvoiceRequest{}
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	resp, err := h.invoicesService.CreateInvoice(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *invoicesHandler) EditInvoice(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// bind and validate request
	req := invoices.EditInvoiceRequest{}
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	resp, err := h.invoicesService.EditInvoice(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *invoicesHandler) DeleteDetailInvoice(c echo.Context) (err error) {
	ctx := c.Request().Context()

	uniqueInvoiceID := c.Param("invoiceID")

	resp, err := h.invoicesService.DeleteDetailInvoice(ctx, uniqueInvoiceID)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, resp)

}
