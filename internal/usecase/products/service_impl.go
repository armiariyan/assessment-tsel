package products

import (
	"github.com/armiariyan/assessment-tsel/internal/domain/repositories"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"

	"context"
)

type service struct {
	productsRepository repositories.ProductsRepository
}

func NewService() *service {
	return &service{}
}

func (s *service) SetProductsRepository(repo repositories.ProductsRepository) *service {
	s.productsRepository = repo
	return s
}

func (s *service) Validate() Service {
	if s.productsRepository == nil {
		panic("productsRepository is nil")
	}

	return s
}

func (s *service) GetListProducts(ctx context.Context, req GetListProductsRequest) (resp constants.DefaultResponse, err error) {
	// products, count, err := s.productsRepository.FindAllAndCount(ctx, req.PaginationRequest, req.Type, req.Name)
	// if err != nil {
	// 	fmt.Printf("failed find all and count products, err : %s", err.Error())
	// 	err = fmt.Errorf("something went wrong[0]")
	// 	return
	// }

	// resp = constants.DefaultResponse{
	// 	Status:  constants.STATUS_SUCCESS,
	// 	Message: constants.MESSAGE_SUCCESS,
	// 	Data: constants.PaginationResponseData{
	// 		Results: products,
	// 		PaginationData: constants.PaginationData{
	// 			Page:          req.Page,
	// 			Limit:         req.Limit,
	// 			TotalPages:    uint(math.Ceil(float64(count) / float64(req.Limit))),
	// 			TotalProducts: uint(count),
	// 			HasNext:       req.Page < uint(math.Ceil(float64(count)/float64(req.Limit))),
	// 			HasPrevious:   req.Page > 1,
	// 		},
	// 	},
	// 	Errors: make([]string, 0),
	// }

	return
}
