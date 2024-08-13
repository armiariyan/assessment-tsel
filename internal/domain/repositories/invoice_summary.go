package repositories

import (
	"context"
	"database/sql"

	"github.com/armiariyan/assessment-tsel/internal/domain/entities"
	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"

	"github.com/armiariyan/bepkg/database/mysql"
	"golang.org/x/sync/errgroup"
)

type InvoiceSummaryRepository interface {
	FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest) (result []entities.InvoiceSummary, count int, err error)
	FindByInvoiceID(ctx context.Context, id int) (result entities.InvoiceSummary, err error)
	InsertWithTx(ctx context.Context, entity entities.InvoiceSummary, tx *sql.Tx) (res sql.Result, err error)
	UpdateWithTx(ctx context.Context, entity entities.InvoiceSummary, tx *sql.Tx) (res sql.Result, err error)
}

type repositoryInvoiceSummary struct {
	db *mysql.SQLDB
}

func NewInvoiceSummaryRepository(db *mysql.SQLDB) *repositoryInvoiceSummary {
	if db == nil {
		panic("db is nil")
	}

	return &repositoryInvoiceSummary{
		db: db,
	}
}

func (r *repositoryInvoiceSummary) FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest) (result []entities.InvoiceSummary, count int, err error) {
	offset := (pagination.Page - 1) * pagination.Limit

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (egErr error) {
		var query = `
		SELECT * FROM invoice_summary ORDER BY id DESC LIMIT ? OFFSET ?`
		return r.db.SelectContext(egCtx, &result, query, pagination.Limit, offset)
	})

	eg.Go(func() (egErr error) {
		var query = `SELECT COUNT(*) FROM invoice_summary`
		if err := r.db.GetContext(egCtx, &count, query); err != nil {
			return err
		}
		return
	})

	err = eg.Wait()

	return
}

func (r *repositoryInvoiceSummary) InsertWithTx(ctx context.Context, entity entities.InvoiceSummary, tx *sql.Tx) (res sql.Result, err error) {
	res, err = tx.ExecContext(ctx,
		`INSERT INTO invoice_summary (
            invoice_id, total_items, subtotal, tax, grand_total
            ) VALUES (?, ?, ?, ?, ?)`,
		entity.InvoiceID, entity.TotalItems, entity.Subtotal, entity.Tax, entity.GrandTotal,
	)

	return
}

func (r *repositoryInvoiceSummary) UpdateWithTx(ctx context.Context, entity entities.InvoiceSummary, tx *sql.Tx) (res sql.Result, err error) {
	query := `UPDATE invoice_summary 
              SET 
                  total_items = ?, 
                  subtotal = ?, 
                  tax = ?, 
                  grand_total = ?, 
                  is_paid = ? 
              WHERE invoice_id = ?`
	res, err = tx.ExecContext(ctx, query,
		entity.TotalItems,
		entity.Subtotal,
		entity.Tax,
		entity.GrandTotal,
		entity.IsPaid,
		entity.InvoiceID)
	return
}

func (r *repositoryInvoiceSummary) FindByInvoiceID(ctx context.Context, id int) (result entities.InvoiceSummary, err error) {
	err = r.db.GetContext(ctx, &result, "SELECT * FROM invoice_summary WHERE invoice_id = ?", id)
	return
}
