package items

import (
	"context"

	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
)

type ItemsService interface {
	GetListItems(ctx context.Context, req GetListItemsRequest) (resp constants.DefaultResponse, err error)
}
