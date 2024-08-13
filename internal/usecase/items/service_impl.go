package items

import (
	"fmt"
	"math"

	"github.com/armiariyan/assessment-tsel/internal/domain/repositories"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"

	"context"

	"github.com/armiariyan/bepkg/database/mysql"
)

type service struct {
	db              *mysql.SQLDB
	itemsRepository repositories.ItemsRepository
}

func NewService() *service {
	return &service{}
}

func (s *service) SetDB(db *mysql.SQLDB) *service {
	s.db = db
	return s
}

func (s *service) SetItemsRepository(repo repositories.ItemsRepository) *service {
	s.itemsRepository = repo
	return s
}

func (s *service) Validate() ItemsService {
	if s.db == nil {
		panic("db is nil")
	}

	if s.itemsRepository == nil {
		panic("itemsRepository is nil")
	}

	return s
}

func (s *service) GetListItems(ctx context.Context, req GetListItemsRequest) (resp constants.DefaultResponse, err error) {
	items, count, err := s.itemsRepository.FindAllAndCount(ctx, req.PaginationRequest, req.Type, req.Name)
	if err != nil {
		fmt.Printf("failed find all and count items, err : %s", err.Error())
		err = fmt.Errorf("something went wrong[0]")
		return
	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data: constants.PaginationResponseData{
			Results: items,
			PaginationData: constants.PaginationData{
				Page:        req.Page,
				Limit:       req.Limit,
				TotalPages:  uint(math.Ceil(float64(count) / float64(req.Limit))),
				TotalItems:  uint(count),
				HasNext:     req.Page < uint(math.Ceil(float64(count)/float64(req.Limit))),
				HasPrevious: req.Page > 1,
			},
		},
		Errors: make([]string, 0),
	}

	return
}
