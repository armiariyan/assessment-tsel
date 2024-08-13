package invoices

import (
	"context"

	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
)

type InvoicesService interface {
	GetListInvoices(ctx context.Context, req GetListInvoicesRequest) (resp constants.DefaultResponse, err error)
	GetDetailInvoice(ctx context.Context, uniqueInvoiceID string) (resp constants.DefaultResponse, err error)
	CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (resp constants.DefaultResponse, err error)
	EditInvoice(ctx context.Context, req EditInvoiceRequest) (resp constants.DefaultResponse, err error)
	DeleteDetailInvoice(ctx context.Context, uniqueInvoiceID string) (resp constants.DefaultResponse, err error)
}
