package repositories

import (
	"context"

	"github.com/armiariyan/assessment-tsel/internal/domain/entities"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
	"github.com/armiariyan/assessment-tsel/internal/pkg/utils"
	"gorm.io/gorm"

	"golang.org/x/sync/errgroup"
)

type ProductsRepository interface {
	FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, conds ...utils.DBCond) (result []entities.Product, count int64, err error)
}

type repositoryProducts struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *repositoryProducts {
	if db == nil {
		panic("db is nil")
	}

	return &repositoryProducts{
		db: db,
	}
}

func (r *repositoryProducts) FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, conds ...utils.DBCond) (result []entities.Product, count int64, err error) {
	limit := pagination.Limit
	offset := (pagination.Page - 1) * pagination.Limit
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (egErr error) {
		queryPayload := r.db.WithContext(egCtx).Limit(int(limit)).Offset(int(offset))
		return utils.CompileConds(queryPayload, conds...).Find(&result).Error
	})
	eg.Go(func() (egErr error) {
		countPayload := r.db.WithContext(egCtx).Model(&entities.Product{})
		return utils.CompileConds(countPayload, conds...).Count(&count).Error
	})
	err = eg.Wait()
	return
}
