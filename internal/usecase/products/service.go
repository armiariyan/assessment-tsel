package products

import (
	"context"

	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
)

type Service interface {
	GetListProducts(ctx context.Context, req GetListProductsRequest) (resp constants.DefaultResponse, err error)
	GetDetailProduct(ctx context.Context, id uint) (resp constants.DefaultResponse, err error)
	CreateProduct(ctx context.Context, req CreateProductRequest) (resp constants.DefaultResponse, err error)
	UpdateProduct(ctx context.Context, req UpdateProductRequest) (resp constants.DefaultResponse, err error)
	DeleteProduct(ctx context.Context, id uint) (resp constants.DefaultResponse, err error)
}
