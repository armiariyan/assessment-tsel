package repositories

import (
	"context"
	"database/sql"

	"github.com/armiariyan/assessment-tsel/internal/domain/entities"

	"github.com/armiariyan/bepkg/database/mysql"
	"golang.org/x/sync/errgroup"
)

type InvoicesRepository interface {
	FindAllAndCountWithCondition(ctx context.Context, params entities.InvoiceListParams) (result []entities.Invoice, count int, err error)
	FindLastDataID(ctx context.Context) (lastID int, err error)
	FindByUniqueInvoiceID(ctx context.Context, invoiceID string) (result entities.Invoice, err error)
	InsertWithTx(ctx context.Context, entity entities.Invoice, tx *sql.Tx) (res sql.Result, err error)
	UpdateWithTx(ctx context.Context, entity entities.Invoice, tx *sql.Tx) (res sql.Result, err error)
	DeleteByUniqueInvoiceIDWithTx(ctx context.Context, invoiceID string, tx *sql.Tx) (res sql.Result, err error)
}

type repositoryInvoices struct {
	db *mysql.SQLDB
}

func NewInvoicesRepository(db *mysql.SQLDB) *repositoryInvoices {
	if db == nil {
		panic("db is nil")
	}

	return &repositoryInvoices{
		db: db,
	}
}

func (r *repositoryInvoices) FindAllAndCountWithCondition(ctx context.Context, params entities.InvoiceListParams) (result []entities.Invoice, count int, err error) {
	offset := (params.Page - 1) * params.Limit

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (egErr error) {
		query := `
		SELECT * FROM invoices
        JOIN invoice_summary ON invoice_summary.invoice_id = invoices.id
        JOIN customers ON customers.id = invoices.customer_id
		WHERE invoices.deleted_at IS NULL`
		args := []interface{}{}

		if params.Subject != "" {
			query += " AND subject LIKE ?"
			args = append(args, "%"+params.Subject+"%")
		}

		if params.InvoiceID != "" {
			query += " AND uq_invoice_id =  ?"
			args = append(args, params.InvoiceID)
		}

		if params.IssueDate != "" {
			query += " AND issue_date =  ?"
			args = append(args, params.IssueDate)
		}

		if params.DueDate != "" {
			query += " AND due_date =  ?"
			args = append(args, params.DueDate)
		}

		// * whom join invoice_summary
		if params.TotalItems != 0 {
			query += " AND invoice_summary.total_items = ?"
			args = append(args, params.TotalItems)
		}

		if params.Status != "" {
			query += " AND invoice_summary.is_paid = "
			if params.Status == "paid" {
				query += "true"
			} else {
				query += "false"
			}
		}

		// * whom join customers
		if params.Customer != "" {
			query += " AND customers.name LIKE ?"
			args = append(args, "%"+params.Customer+"%")
		}

		query += " ORDER BY invoices.id DESC LIMIT ? OFFSET ?"
		args = append(args, params.Limit, offset)

		err = r.db.SelectContext(ctx, &result, query, args...)
		return
	})

	eg.Go(func() (egErr error) {
		query := `
		SELECT COUNT(*) FROM invoices 
		JOIN invoice_summary ON invoice_summary.invoice_id = invoices.id
		JOIN customers ON customers.id = invoices.customer_id
		WHERE invoices.deleted_at IS NULL`
		args := []interface{}{}

		if params.Subject != "" {
			query += " AND subject LIKE ?"
			args = append(args, "%"+params.Subject+"%")
		}

		if params.InvoiceID != "" {
			query += " AND uq_invoice_id =  ?"
			args = append(args, params.InvoiceID)
		}

		if params.IssueDate != "" {
			query += " AND issue_date =  ?"
			args = append(args, params.IssueDate)
		}

		if params.DueDate != "" {
			query += " AND due_date =  ?"
			args = append(args, params.DueDate)
		}

		// * whom join invoice_summary
		if params.TotalItems != 0 {
			query += " AND invoice_summary.total_items = ?"
			args = append(args, params.TotalItems)
		}

		if params.Status != "" {
			query += " AND invoice_summary.is_paid = "
			if params.Status == "paid" {
				query += "true"
			} else {
				query += "false"
			}
		}

		// * whom join customers
		if params.Customer != "" {
			query += " AND customers.name LIKE ?"
			args = append(args, "%"+params.Customer+"%")
		}

		if err := r.db.GetContext(egCtx, &count, query, args...); err != nil && err != sql.ErrNoRows {
			return err
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		return nil, 0, err
	}

	return result, count, nil
}

func (r *repositoryInvoices) InsertWithTx(ctx context.Context, entity entities.Invoice, tx *sql.Tx) (res sql.Result, err error) {
	res, err = tx.ExecContext(ctx,
		`INSERT INTO invoices (
            uq_invoice_id, issue_date, due_date, subject, customer_id
            ) VALUES (?, ?, ?, ?, ?)`,
		entity.InvoiceID, entity.IssueDate, entity.DueDate, entity.Subject, entity.CustomerID,
	)

	return
}

func (r *repositoryInvoices) UpdateWithTx(ctx context.Context, entity entities.Invoice, tx *sql.Tx) (res sql.Result, err error) {
	query := `UPDATE invoices 
              SET 
                  uq_invoice_id = ?, 
                  issue_date = ?, 
                  due_date = ?, 
                  subject = ?, 
                  customer_id = ?
              WHERE id = ?`
	res, err = tx.ExecContext(ctx, query,
		entity.InvoiceID,
		entity.IssueDate,
		entity.DueDate,
		entity.Subject,
		entity.CustomerID,
		entity.ID)
	return
}

func (r *repositoryInvoices) FindLastDataID(ctx context.Context) (lastID int, err error) {
	err = r.db.GetContext(ctx, &lastID, "SELECT COUNT(*) FROM invoices")
	return
}

func (r *repositoryInvoices) FindByUniqueInvoiceID(ctx context.Context, invoiceID string) (result entities.Invoice, err error) {
	err = r.db.GetContext(ctx, &result, "SELECT id as pk_invoice_id, uq_invoice_id, issue_date, due_date, subject, customer_id FROM invoices WHERE uq_invoice_id = ?", invoiceID)

	return
}

func (r *repositoryInvoices) DeleteByUniqueInvoiceIDWithTx(ctx context.Context, invoiceID string, tx *sql.Tx) (res sql.Result, err error) {
	res, err = tx.ExecContext(ctx,
		`DELETE FROM invoices WHERE uq_invoice_id = ?`,
		invoiceID,
	)

	return
}
