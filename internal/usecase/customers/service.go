package customers

import (
	"context"

	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
)

type CustomersService interface {
	GetListCustomers(ctx context.Context, req GetListCustomersRequest) (resp constants.DefaultResponse, err error)
}
