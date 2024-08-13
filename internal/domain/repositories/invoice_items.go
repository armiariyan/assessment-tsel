package repositories

import (
	"context"
	"database/sql"

	"github.com/armiariyan/assessment-tsel/internal/domain/entities"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"

	"github.com/armiariyan/bepkg/database/mysql"
	"golang.org/x/sync/errgroup"
)

type InvoiceItemsRepository interface {
	FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest) (result []entities.InvoiceItems, count int, err error)
	FindByInvoiceID(ctx context.Context, id int) (result []entities.InvoiceItems, err error)
	InsertWithTx(ctx context.Context, entity entities.InvoiceItems, tx *sql.Tx) (res sql.Result, err error)
	UpdateWithTx(ctx context.Context, entity entities.InvoiceItems, tx *sql.Tx) (res sql.Result, err error)
	BulkInsertWithTx(ctx context.Context, entities []entities.InvoiceItems, tx *sql.Tx) (res sql.Result, err error)
	DeleteWithTx(ctx context.Context, id int, tx *sql.Tx) (res sql.Result, err error)
}

type repositoryInvoiceItems struct {
	db *mysql.SQLDB
}

func NewInvoiceItemsRepository(db *mysql.SQLDB) *repositoryInvoiceItems {
	if db == nil {
		panic("db is nil")
	}

	return &repositoryInvoiceItems{
		db: db,
	}
}

func (r *repositoryInvoiceItems) FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest) (result []entities.InvoiceItems, count int, err error) {
	offset := (pagination.Page - 1) * pagination.Limit

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (egErr error) {
		var query = `
		SELECT * FROM invoice_items ORDER BY id DESC LIMIT ? OFFSET ?`
		return r.db.SelectContext(egCtx, &result, query, pagination.Limit, offset)
	})

	eg.Go(func() (egErr error) {
		var query = `SELECT COUNT(*) FROM invoice_items`
		if err := r.db.GetContext(egCtx, &count, query); err != nil {
			return err
		}
		return
	})

	err = eg.Wait()

	return
}

func (r *repositoryInvoiceItems) InsertWithTx(ctx context.Context, entity entities.InvoiceItems, tx *sql.Tx) (res sql.Result, err error) {
	res, err = tx.ExecContext(ctx,
		`INSERT INTO invoice_items (
            invoice_id, item_id, quantity, unit_price, amount
            ) VALUES (?, ?, ?, ?, ?)`,
		entity.InvoiceID, entity.ItemID, entity.Quantity, entity.UnitPrice, entity.Amount,
	)
	return
}

func (r *repositoryInvoiceItems) UpdateWithTx(ctx context.Context, entity entities.InvoiceItems, tx *sql.Tx) (res sql.Result, err error) {
	query := `UPDATE invoice_items 
              SET 
                  quantity = ?, 
                  unit_price = ?, 
                  amount = ? 
              WHERE invoice_id = ? AND item_id = ?`
	res, err = tx.ExecContext(ctx, query,
		entity.Quantity,
		entity.UnitPrice,
		entity.Amount,
		entity.InvoiceID,
		entity.ItemID)
	return
}

func (r *repositoryInvoiceItems) BulkInsertWithTx(ctx context.Context, entities []entities.InvoiceItems, tx *sql.Tx) (res sql.Result, err error) {
	query := `INSERT INTO invoice_items (
                invoice_id, item_id, quantity, unit_price, amount
              ) VALUES `

	// * create the values part of the query
	values := ""
	args := []interface{}{}

	for i, entity := range entities {
		if i > 0 {
			values += ", "
		}
		values += "(?, ?, ?, ?, ?)"
		args = append(args, entity.InvoiceID, entity.ItemID, entity.Quantity, entity.UnitPrice, entity.Amount)
	}

	// * complete the query
	query += values

	// * execute the bulk insert query
	res, err = tx.ExecContext(ctx, query, args...)
	return
}

func (r *repositoryInvoiceItems) FindByInvoiceID(ctx context.Context, id int) (result []entities.InvoiceItems, err error) {
	err = r.db.SelectContext(ctx, &result, "SELECT * FROM invoice_items WHERE invoice_id = ?", id)
	return
}

func (r *repositoryInvoiceItems) DeleteWithTx(ctx context.Context, id int, tx *sql.Tx) (res sql.Result, err error) {
	query := `DELETE FROM invoice_items WHERE id = ?`
	res, err = tx.ExecContext(ctx, query, id)
	return
}
