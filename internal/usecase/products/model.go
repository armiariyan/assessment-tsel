package products

import (
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
)

// * Requests
type (
	GetListProductsRequest struct {
		constants.PaginationRequest
		Type string `query:"type" validate:"omitempty"`
		Name string `query:"name" validate:"omitempty"`
	}
)

// * Responses
type ()
