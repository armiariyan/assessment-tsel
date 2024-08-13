package handler

import (
	"net/http"

	"github.com/armiariyan/assessment-tsel/internal/pkg/utils"
	"github.com/armiariyan/assessment-tsel/internal/usecase/items"

	"github.com/labstack/echo/v4"
)

type itemsHandler struct {
	itemsService items.ItemsService
}

func NewItemsHandler() *itemsHandler {
	return &itemsHandler{}
}

func (h *itemsHandler) SetItemsService(service items.ItemsService) *itemsHandler {
	h.itemsService = service
	return h
}

func (h *itemsHandler) Validate() *itemsHandler {
	if h.itemsService == nil {
		panic("itemsService is nil")
	}

	return h
}

func (h *itemsHandler) GetListItems(c echo.Context) (err error) {
	ctx := c.Request().Context()

	req := items.GetListItemsRequest{}
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	resp, err := h.itemsService.GetListItems(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, resp)

}
