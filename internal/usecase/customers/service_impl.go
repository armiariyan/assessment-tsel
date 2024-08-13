package customers

import (
	"math"

	"github.com/armiariyan/assessment-tsel/internal/domain/repositories"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
	"github.com/armiariyan/assessment-tsel/internal/pkg/log"

	"context"

	"github.com/armiariyan/bepkg/database/mysql"
)

type service struct {
	db                  *mysql.SQLDB
	customersRepository repositories.CustomersRepository
}

func NewService() *service {
	return &service{}
}

func (s *service) SetDB(db *mysql.SQLDB) *service {
	s.db = db
	return s
}

func (s *service) SetCustomersRepository(repo repositories.CustomersRepository) *service {
	s.customersRepository = repo
	return s
}

func (s *service) Validate() CustomersService {
	if s.db == nil {
		panic("db is nil")
	}

	if s.customersRepository == nil {
		panic("customersRepository is nil")
	}

	return s
}

func (s *service) GetListCustomers(ctx context.Context, req GetListCustomersRequest) (resp constants.DefaultResponse, err error) {
	customers, count, err := s.customersRepository.FindAllAndCount(ctx, req.PaginationRequest)
	if err != nil {
		log.Error(ctx, "failed find all and count customers", err)
		return
	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data: constants.PaginationResponseData{
			Results: customers,
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
