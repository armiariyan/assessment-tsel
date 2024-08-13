package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/armiariyan/assessment-tsel/internal/domain/entities"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"

	"github.com/armiariyan/bepkg/database/mysql"
	"golang.org/x/sync/errgroup"
)

type ItemsRepository interface {
	FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, itemType, name string) (result []entities.Item, count int, err error)
	FindByIDs(ctx context.Context, ids []int) (result []entities.Item, err error)
}

type repositoryItems struct {
	db *mysql.SQLDB
}

func NewItemsRepository(db *mysql.SQLDB) *repositoryItems {
	if db == nil {
		panic("db is nil")
	}

	return &repositoryItems{
		db: db,
	}
}

func (r *repositoryItems) FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, itemType, name string) (result []entities.Item, count int, err error) {
	offset := (pagination.Page - 1) * pagination.Limit

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() (egErr error) {
		query := "SELECT * FROM items WHERE deleted_at IS NULL"
		args := []interface{}{}

		if itemType != "" {
			query += " AND type LIKE ?"
			args = append(args, "%"+itemType+"%")
		}

		if name != "" {
			query += " AND name LIKE ?"
			args = append(args, "%"+name+"%")
		}

		query += " ORDER BY id DESC LIMIT ? OFFSET ?"
		args = append(args, pagination.Limit, offset)

		if err = r.db.SelectContext(ctx, &result, query, args...); err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() (egErr error) {
		query := "SELECT COUNT(*) FROM items WHERE deleted_at IS NULL"
		args := []interface{}{}

		if itemType != "" {
			query += " AND type LIKE ?"
			args = append(args, "%"+itemType+"%")
		}

		if name != "" {
			query += " AND name LIKE ?"
			args = append(args, "%"+name+"%")
		}

		if err := r.db.GetContext(egCtx, &count, query, args...); err != nil {
			return err
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		return nil, 0, err
	}

	return result, count, nil
}

func (r *repositoryItems) FindByIDs(ctx context.Context, ids []int) (result []entities.Item, err error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))

	// * create placeholder
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	// * join placeholder
	query := fmt.Sprintf("SELECT * FROM items WHERE id IN (%s)", strings.Join(placeholders, ","))

	err = r.db.SelectContext(ctx, &result, query, args...)

	return
}
