package products

import (
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
	"gorm.io/datatypes"
)

// * Requests
type (
	GetListProductsRequest struct {
		constants.PaginationRequest
	}

	CreateProductRequest struct {
		Name        string         `json:"name" validate:"required,min=3,max=255"`
		Description string         `json:"description,omitempty" validate:"max=1000"`
		Price       float64        `json:"price" validate:"required,gt=0"`
		Variety     datatypes.JSON `json:"variety" validate:"omitempty"`
		Rating      *float64       `json:"rating,omitempty" validate:"omitempty,gte=0,lte=5"`
		Stock       float64        `json:"stock" validate:"required,gte=0"`
	}

	UpdateProductRequest struct {
		ID          uint           `json:"id" validate:"required,number"`
		Name        string         `json:"name" validate:"min=3,max=255"`
		Description string         `json:"description,omitempty" validate:"max=1000"`
		Price       float64        `json:"price" validate:"gt=0"`
		Variety     datatypes.JSON `json:"variety" validate:"omitempty"`
		Rating      *float64       `json:"rating,omitempty" validate:"omitempty,gte=0,lte=5"`
		Stock       float64        `json:"stock" validate:"gte=0"`
	}
)

// * Responses
type ()
