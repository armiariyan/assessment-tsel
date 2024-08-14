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
	FindByIDOrError(ctx context.Context, id uint) (result entities.Product, err error)
	Create(ctx context.Context, entity *entities.Product) (err error)
	UpdateByID(ctx context.Context, id uint, entity *entities.Product) (result entities.Product, err error)
	DeleteByID(ctx context.Context, id uint) (err error)
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

func (r *repositoryProducts) FindByIDOrError(ctx context.Context, id uint) (result entities.Product, err error) {
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return
}

func (r *repositoryProducts) Create(ctx context.Context, entity *entities.Product) (err error) {
	err = r.db.WithContext(ctx).Create(&entity).Error
	return
}

func (r *repositoryProducts) UpdateByID(ctx context.Context, id uint, entity *entities.Product) (result entities.Product, err error) {
	tx := r.db.WithContext(ctx).Where("id = ?", id).Updates(&entity)
	err = tx.Error
	if err == nil && tx.RowsAffected < 1 {
		err = gorm.ErrRecordNotFound
		return
	}

	err = r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return
}

func (r *repositoryProducts) DeleteByID(ctx context.Context, id uint) (err error) {
	err = r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Product{}).Error
	return

}
