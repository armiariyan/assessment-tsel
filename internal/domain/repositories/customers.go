package repositories

import (
	"context"

	"github.com/armiariyan/assessment-tsel/internal/domain/entities"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"

	"github.com/armiariyan/bepkg/database/mysql"
	"golang.org/x/sync/errgroup"
)

type CustomersRepository interface {
	FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest) (result []entities.Customer, count int, err error)
	FindByID(ctx context.Context, id int) (result entities.Customer, err error)
}

type repositoryCustomers struct {
	db *mysql.SQLDB
}

func NewCustomersRepository(db *mysql.SQLDB) *repositoryCustomers {
	if db == nil {
		panic("db is nil")
	}

	return &repositoryCustomers{
		db: db,
	}
}

func (r *repositoryCustomers) FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest) (result []entities.Customer, count int, err error) {
	offset := (pagination.Page - 1) * pagination.Limit

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (egErr error) {
		var query = `
		SELECT * FROM customers ORDER BY id DESC LIMIT ? OFFSET ?`
		return r.db.SelectContext(egCtx, &result, query, pagination.Limit, offset)
	})

	eg.Go(func() (egErr error) {
		var query = `SELECT COUNT(*) FROM customers`
		if err := r.db.GetContext(egCtx, &count, query); err != nil {
			return err
		}
		return
	})

	err = eg.Wait()

	return
}

func (r *repositoryCustomers) FindByID(ctx context.Context, id int) (result entities.Customer, err error) {
	err = r.db.GetContext(ctx, &result, "SELECT * FROM customers WHERE id = ?", id)
	return
}
