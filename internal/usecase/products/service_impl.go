package products

import (
	"fmt"
	"math"

	"github.com/armiariyan/assessment-tsel/internal/domain/entities"
	"github.com/armiariyan/assessment-tsel/internal/domain/repositories"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
	"github.com/armiariyan/assessment-tsel/internal/pkg/log"
	"gorm.io/gorm"

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
	products, count, err := s.productsRepository.FindAllAndCount(ctx, req.PaginationRequest)
	if err != nil {
		log.Error(ctx, "failed to find list products", err)
		err = fmt.Errorf("something went wrong. Please try again later (1)")
		return
	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data: constants.PaginationResponseData{
			Results: products,
			PaginationData: constants.PaginationData{
				Page:        req.Page,
				Limit:       req.Limit,
				TotalPages:  uint(math.Ceil(float64(count) / float64(req.Limit))),
				TotalItems:  uint(count),
				HasNext:     req.Page < uint(math.Ceil(float64(count)/float64(req.Limit))),
				HasPrevious: req.Page > 1,
			},
		},
	}

	return
}

func (s *service) GetDetailProduct(ctx context.Context, id uint) (resp constants.DefaultResponse, err error) {
	product, err := s.productsRepository.FindByIDOrError(ctx, id)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to find product with id %d during get detail product", id), err)
		if err == gorm.ErrRecordNotFound {
			resp = constants.ErrorResponse(constants.STATUS_DATA_NOT_FOUND, "data product not found")
			return
		}
		err = fmt.Errorf("something went wrong. Please try again later (1)")
		return
	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    product,
	}

	return
}

func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (resp constants.DefaultResponse, err error) {
	// * mapping product request to entity
	payload := entities.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Variety:     req.Variety,
	}

	if req.Rating != nil {
		rating := math.Round(*req.Rating)
		payload.Rating = &rating
	}

	err = s.productsRepository.Create(ctx, &payload)
	if err != nil {
		log.Error(ctx, "failed to create product", payload, err)
		err = fmt.Errorf("something went wrong. Please try again later (1)")
		return
	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
	}

	return
}

func (s *service) UpdateProduct(ctx context.Context, req UpdateProductRequest) (resp constants.DefaultResponse, err error) {
	// * get the product by given id
	product, err := s.productsRepository.FindByIDOrError(ctx, req.ID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to find product with id %d during update product", req.ID), err)
		if err == gorm.ErrRecordNotFound {
			resp = constants.ErrorResponse(constants.STATUS_DATA_NOT_FOUND, "data product not found")
			return
		}
		err = fmt.Errorf("something went wrong. Please try again later (1)")
		return
	}

	// * mapping product request to entity
	payload := entities.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Variety:     req.Variety,
	}

	if req.Rating != nil {
		log.Info(ctx, fmt.Sprintf("[UPDATE] count new rating for product %s with id %d", product.Name, product.ID))
		log.Info(ctx, fmt.Sprintf("rating before %v", &product.Rating))

		var ratingValue float64
		if product.Rating == nil {
			ratingValue = *req.Rating
		} else {
			ratingValue = (*product.Rating + *req.Rating) / 2
		}

		roundedRating := math.Round(ratingValue*10) / 10
		payload.Rating = &roundedRating

		log.Info(ctx, fmt.Sprintf("rating after %v", &payload.Rating))
	}

	updatedProduct, err := s.productsRepository.UpdateByID(ctx, product.ID, &payload)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to update product with id %d", product.ID), err)
		err = fmt.Errorf("something went wrong. Please try again later (2)")
		return
	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS_UPDATE,
		Data:    updatedProduct,
	}

	return
}

func (s *service) DeleteProduct(ctx context.Context, id uint) (resp constants.DefaultResponse, err error) {
	// * get the product by given id
	_, err = s.productsRepository.FindByIDOrError(ctx, id)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to find product with id %d during delete product", id), err)
		if err == gorm.ErrRecordNotFound {
			resp = constants.ErrorResponse(constants.STATUS_DATA_NOT_FOUND, "data product not found")
			return
		}
		err = fmt.Errorf("something went wrong. Please try again later (1)")
		return
	}

	err = s.productsRepository.DeleteByID(ctx, id)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to delete product by id %d", id), err)
		err = fmt.Errorf("something went wrong. Please try again later (2)")
		return
	}

	resp = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS_DELETE,
	}

	return
}
