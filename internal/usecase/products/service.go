package products

import (
	"context"

	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
)

type Service interface {
	GetListProducts(ctx context.Context, req GetListProductsRequest) (resp constants.DefaultResponse, err error)
}
