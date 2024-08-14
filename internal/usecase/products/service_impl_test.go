package products

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/armiariyan/assessment-tsel/internal/domain/entities"
	mocksRepo "github.com/armiariyan/assessment-tsel/internal/domain/repositories/mocks"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
	"github.com/armiariyan/assessment-tsel/internal/pkg/log"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"go.uber.org/mock/gomock"
)

var (
	now                     = time.Now()
	rating         float64  = 4.5
	exampleRating  *float64 = &rating
	exampleVariety          = datatypes.JSON([]byte(`{"color": "red", "size": "M", "weight": 1.2}`))
)

func init() {
	log.New()
}

func TestValidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoProduct := mocksRepo.NewMockProductsRepository(ctrl)

	service := NewService().
		SetProductsRepository(mockRepoProduct)

	t.Run("panic when productRepository is nil", func(t *testing.T) {
		service.SetProductsRepository(nil)
		require.Panics(t, func() {
			service.Validate()
		}, "productRepository is nil")
	})

	service.SetProductsRepository(mockRepoProduct)

	t.Run("no panic when all are set", func(t *testing.T) {
		require.NotPanics(t, func() {
			service.Validate()
		}, "positive case")
	})
}

func TestProductService_GetListProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsRepo := mocksRepo.NewMockProductsRepository(ctrl)

	service := &service{
		productsRepository: mockProductsRepo,
	}

	tests := []struct {
		name              string
		req               GetListProductsRequest
		doMockProductRepo func(mock *mocksRepo.MockProductsRepository)
		wantRes           constants.DefaultResponse
		wantErr           error
	}{
		{
			name: "positive case",
			req: GetListProductsRequest{
				PaginationRequest: constants.PaginationRequest{
					Page:  1,
					Limit: 10,
				},
			},
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindAllAndCount(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]entities.Product{
						{
							ID:          1,
							Name:        "Test Name",
							Description: "Test Description",
							Price:       100000,
							Stock:       150,
							Rating:      exampleRating,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
						{
							ID:          2,
							Name:        "Test Name 2",
							Description: "Test Description 2",
							Price:       100000,
							Stock:       150,
							Rating:      exampleRating,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
					}, int64(2), nil).Times(1)
			},
			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_SUCCESS,
				Message: constants.MESSAGE_SUCCESS,
				Data: constants.PaginationResponseData{
					Results: []entities.Product{
						{
							ID:          1,
							Name:        "Test Name",
							Description: "Test Description",
							Price:       100000,
							Stock:       150,
							Rating:      exampleRating,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
						{
							ID:          2,
							Name:        "Test Name 2",
							Description: "Test Description 2",
							Price:       100000,
							Stock:       150,
							Rating:      exampleRating,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
					},

					PaginationData: constants.PaginationData{
						Page:        1,
						Limit:       10,
						TotalPages:  1,
						TotalItems:  2,
						HasNext:     false,
						HasPrevious: false,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "negative case - failed find all and count product",
			req: GetListProductsRequest{
				PaginationRequest: constants.PaginationRequest{
					Page:  1,
					Limit: 10,
				},
			},
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindAllAndCount(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]entities.Product{}, int64(0), errors.New("something went wrong. Please try again later (1)")).Times(1)
			},
			wantRes: constants.DefaultResponse{},
			wantErr: errors.New("something went wrong. Please try again later (1)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.doMockProductRepo(mockProductsRepo)

			resp, err := service.GetListProducts(context.TODO(), tt.req)
			require.Equal(t, tt.wantRes, resp)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestProductService_GetDetailProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsRepo := mocksRepo.NewMockProductsRepository(ctrl)

	service := &service{
		productsRepository: mockProductsRepo,
	}

	tests := []struct {
		name              string
		req               uint
		doMockProductRepo func(mock *mocksRepo.MockProductsRepository)
		wantRes           constants.DefaultResponse
		wantErr           error
	}{
		{
			name: "positive case",
			req:  uint(1),
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{
						ID:          1,
						Name:        "Test Name",
						Description: "Test Description",
						Price:       100000,
						Stock:       150,
						Rating:      exampleRating,
						CreatedAt:   now,
						UpdatedAt:   now,
					}, nil).Times(1)
			},
			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_SUCCESS,
				Message: constants.MESSAGE_SUCCESS,
				Data: entities.Product{
					ID:          1,
					Name:        "Test Name",
					Description: "Test Description",
					Price:       100000,
					Stock:       150,
					Rating:      exampleRating,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: nil,
		},
		{
			name: "negative case - data product not found",
			req:  uint(1),
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{}, gorm.ErrRecordNotFound).Times(1)
			},
			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_DATA_NOT_FOUND,
				Message: "data product not found",
				Data:    struct{}{},
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "negative case - failed find by id",
			req:  uint(1),
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{}, errors.New("something went wrong. Please try again later (1)")).Times(1)
			},
			wantRes: constants.DefaultResponse{},
			wantErr: errors.New("something went wrong. Please try again later (1)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.doMockProductRepo(mockProductsRepo)

			resp, err := service.GetDetailProduct(context.TODO(), tt.req)
			require.Equal(t, tt.wantRes, resp)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestProductService_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsRepo := mocksRepo.NewMockProductsRepository(ctrl)

	service := &service{
		productsRepository: mockProductsRepo,
	}

	tests := []struct {
		name              string
		req               CreateProductRequest
		doMockProductRepo func(mock *mocksRepo.MockProductsRepository)
		wantRes           constants.DefaultResponse
		wantErr           error
	}{
		{
			name: "positive case",
			req: CreateProductRequest{
				Name:        "test name product",
				Description: "test description product",
				Price:       100000,
				Variety:     exampleVariety,
				Rating:      exampleRating,
				Stock:       150,
			},
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_SUCCESS,
				Message: constants.MESSAGE_SUCCESS,
			},
			wantErr: nil,
		},
		{
			name: "negative case - failed create product",
			req: CreateProductRequest{
				Name:        "test name product",
				Description: "test description product",
				Price:       100000,
				Variety:     exampleVariety,
				Rating:      exampleRating,
				Stock:       150,
			},
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("something went wrong. Please try again later (1)")).Times(1)
			},
			wantRes: constants.DefaultResponse{},
			wantErr: errors.New("something went wrong. Please try again later (1)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.doMockProductRepo(mockProductsRepo)

			resp, err := service.CreateProduct(context.TODO(), tt.req)
			require.Equal(t, tt.wantRes, resp)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsRepo := mocksRepo.NewMockProductsRepository(ctrl)

	service := &service{
		productsRepository: mockProductsRepo,
	}

	tests := []struct {
		name              string
		req               UpdateProductRequest
		doMockProductRepo func(mock *mocksRepo.MockProductsRepository)
		wantRes           constants.DefaultResponse
		wantErr           error
	}{
		{
			name: "positive case",
			req: UpdateProductRequest{
				ID:          1,
				Name:        "UPDATED test name product",
				Description: "test description product",
				Price:       100000,
				Variety:     exampleVariety,
				Rating:      exampleRating,
				Stock:       150,
			},
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{
						ID:          1,
						Name:        "UPDATED test name product",
						Description: "test description product",
						Price:       100000,
						Variety:     exampleVariety,
						Rating:      exampleRating,
						Stock:       150,
					}, nil).Times(1)
				mock.EXPECT().UpdateByID(gomock.Any(), uint(1), gomock.Any()).Return(
					entities.Product{
						ID:          1,
						Name:        "UPDATED test name product",
						Description: "test description product",
						Price:       100000,
						Variety:     exampleVariety,
						Rating:      exampleRating,
						Stock:       150,
					}, nil).Times(1)
			},
			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_SUCCESS,
				Message: constants.MESSAGE_SUCCESS_UPDATE,
				Data: entities.Product{
					ID:          1,
					Name:        "UPDATED test name product",
					Description: "test description product",
					Price:       100000,
					Variety:     exampleVariety,
					Rating:      exampleRating,
					Stock:       150,
				},
			},
			wantErr: nil,
		},
		{
			name: "positive case - product rating empty (first rating)",
			req: UpdateProductRequest{
				ID:          1,
				Name:        "UPDATED test name product",
				Description: "test description product",
				Price:       100000,
				Variety:     exampleVariety,
				Rating:      exampleRating,
				Stock:       150,
			},
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{
						ID:          1,
						Name:        "UPDATED test name product",
						Description: "test description product",
						Price:       100000,
						Variety:     exampleVariety,
						Rating:      nil,
						Stock:       150,
					}, nil).Times(1)
				mock.EXPECT().UpdateByID(gomock.Any(), uint(1), gomock.Any()).Return(
					entities.Product{
						ID:          1,
						Name:        "UPDATED test name product",
						Description: "test description product",
						Price:       100000,
						Variety:     exampleVariety,
						Rating:      exampleRating,
						Stock:       150,
					}, nil).Times(1)
			},
			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_SUCCESS,
				Message: constants.MESSAGE_SUCCESS_UPDATE,
				Data: entities.Product{
					ID:          1,
					Name:        "UPDATED test name product",
					Description: "test description product",
					Price:       100000,
					Variety:     exampleVariety,
					Rating:      exampleRating,
					Stock:       150,
				},
			},
			wantErr: nil,
		},
		{
			name: "negative case - data product not found",
			req: UpdateProductRequest{
				ID:          1,
				Name:        "UPDATED test name product",
				Description: "test description product",
				Price:       100000,
				Variety:     exampleVariety,
				Rating:      exampleRating,
				Stock:       150,
			},
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{}, gorm.ErrRecordNotFound).Times(1)
			},
			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_DATA_NOT_FOUND,
				Message: "data product not found",
				Data:    struct{}{},
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "negative case - failed find by id",
			req: UpdateProductRequest{
				ID:          1,
				Name:        "UPDATED test name product",
				Description: "test description product",
				Price:       100000,
				Variety:     exampleVariety,
				Rating:      exampleRating,
				Stock:       150,
			},
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{}, errors.New("something went wrong. Please try again later (1)")).Times(1)
			},
			wantRes: constants.DefaultResponse{},
			wantErr: errors.New("something went wrong. Please try again later (1)"),
		},
		{
			name: "negative case - failed update by id",
			req: UpdateProductRequest{
				ID:          1,
				Name:        "UPDATED test name product",
				Description: "test description product",
				Price:       100000,
				Variety:     exampleVariety,
				Rating:      exampleRating,
				Stock:       150,
			},
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{
						ID:          1,
						Name:        "UPDATED test name product",
						Description: "test description product",
						Price:       100000,
						Variety:     exampleVariety,
						Rating:      exampleRating,
						Stock:       150,
					}, nil).Times(1)
				mock.EXPECT().UpdateByID(gomock.Any(), uint(1), gomock.Any()).Return(
					entities.Product{}, errors.New("something went wrong. Please try again later (2)")).Times(1)
			},
			wantRes: constants.DefaultResponse{},
			wantErr: errors.New("something went wrong. Please try again later (2)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.doMockProductRepo(mockProductsRepo)

			resp, err := service.UpdateProduct(context.TODO(), tt.req)
			require.Equal(t, tt.wantRes, resp)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsRepo := mocksRepo.NewMockProductsRepository(ctrl)

	service := &service{
		productsRepository: mockProductsRepo,
	}

	tests := []struct {
		name              string
		req               uint
		doMockProductRepo func(mock *mocksRepo.MockProductsRepository)
		wantRes           constants.DefaultResponse
		wantErr           error
	}{
		{
			name: "positive case",
			req:  uint(1),
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{
						ID:          1,
						Name:        "UPDATED test name product",
						Description: "test description product",
						Price:       100000,
						Variety:     exampleVariety,
						Rating:      exampleRating,
						Stock:       150,
					}, nil).Times(1)
				mock.EXPECT().DeleteByID(gomock.Any(), uint(1)).Return(nil).Times(1)
			},
			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_SUCCESS,
				Message: constants.MESSAGE_SUCCESS_DELETE,
			},
			wantErr: nil,
		},
		{
			name: "negative case - data product not found",
			req:  uint(1),
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{}, gorm.ErrRecordNotFound).Times(1)
			},
			wantRes: constants.DefaultResponse{
				Status:  constants.STATUS_DATA_NOT_FOUND,
				Message: "data product not found",
				Data:    struct{}{},
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "negative case - failed find by id",
			req:  uint(1),
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{}, errors.New("something went wrong. Please try again later (1)")).Times(1)
			},
			wantRes: constants.DefaultResponse{},
			wantErr: errors.New("something went wrong. Please try again later (1)"),
		},
		{
			name: "negative case - failed delete by id",
			req:  uint(1),
			doMockProductRepo: func(mock *mocksRepo.MockProductsRepository) {
				mock.EXPECT().FindByIDOrError(gomock.Any(), uint(1)).Return(
					entities.Product{
						ID:          1,
						Name:        "UPDATED test name product",
						Description: "test description product",
						Price:       100000,
						Variety:     exampleVariety,
						Rating:      exampleRating,
						Stock:       150,
					}, nil).Times(1)
				mock.EXPECT().DeleteByID(gomock.Any(), uint(1)).Return(errors.New("something went wrong. Please try again later (2)")).Times(1)
			},
			wantRes: constants.DefaultResponse{},
			wantErr: errors.New("something went wrong. Please try again later (2)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.doMockProductRepo(mockProductsRepo)

			resp, err := service.DeleteProduct(context.TODO(), tt.req)
			require.Equal(t, tt.wantRes, resp)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
