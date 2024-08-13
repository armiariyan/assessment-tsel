package invoices

import (
	"time"

	"github.com/armiariyan/assessment-tsel/internal/pkg/constants"
)

// * Requests
type (
	GetListInvoicesRequest struct {
		constants.PaginationRequest
		InvoiceID  string `query:"invoiceId" validate:"omitempty"`
		IssueDate  string `query:"issueDate" validate:"omitempty"`
		Subject    string `query:"subject" validate:"omitempty"`
		Customer   string `query:"customer" validate:"omitempty"`
		DueDate    string `query:"dueDate" validate:"omitempty"`
		Status     string `query:"status" validate:"omitempty"`
		TotalItems int    `query:"totalItems" validate:"omitempty"`
	}

	CreateInvoiceRequest struct {
		IssueDate      string         `json:"issueDate" validate:"required"`
		DueDate        string         `json:"dueDate" validate:"required"`
		Subject        string         `json:"subject" validate:"required"`
		CustomerID     int            `json:"customerId" validate:"required,number"`
		InvoiceItems   []InvoiceItem  `json:"invoiceItems" validate:"required,dive"`
		InvoiceSummary InvoiceSummary `json:"invoiceSummary" validate:"required,dive"`
	}

	InvoiceItem struct {
		ItemID    int     `json:"itemId" validate:"required,number"`
		Quantity  float64 `json:"quantity" validate:"required,number"`
		UnitPrice float64 `json:"unitPrice" validate:"required,number"`
		Amount    float64 `json:"amount" validate:"required,number"`
	}

	InvoiceSummary struct {
		TotalItems int     `json:"totalItems" validate:"required,number"`
		SubTotal   float64 `json:"subTotal" validate:"required,number"`
		Tax        float64 `json:"tax" validate:"required,number"`
		GrandTotal float64 `json:"grandTotal" validate:"required,number"`
	}

	EditInvoiceRequest struct {
		InvoiceID      string         `json:"invoiceId" validate:"required"`
		IssueDate      string         `json:"issueDate" validate:"required"`
		DueDate        string         `json:"dueDate" validate:"required"`
		Subject        string         `json:"subject" validate:"required"`
		CustomerID     int            `json:"customerId" validate:"required,number"`
		InvoiceItems   []InvoiceItem  `json:"invoiceItems" validate:"required,dive"`
		InvoiceSummary InvoiceSummary `json:"invoiceSummary" validate:"required,dive"`
	}
)

// * Responses
type (
	GetListInvoicesResponse struct {
		InvoiceID    string    `json:"invoiceId"`
		Subject      string    `json:"subject"`
		CustomerName string    `json:"customerName"`
		Status       string    `json:"status"`
		TotalItems   int       `json:"total_items"`
		IssueDate    time.Time `json:"issueDate"`
		DueDate      time.Time `json:"dueDate"`
	}

	GetDetailInvoiceResponse struct {
		InvoiceID string    `json:"invoiceId"`
		Subject   string    `json:"subject"`
		IssueDate time.Time `json:"issueDate"`
		DueDate   time.Time `json:"dueDate"`
		Customer  Customer  `json:"customer"`
		Items     []Item    `json:"items"`
		Summary   Summary   `json:"summary"`
	}

	Customer struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		City     string `json:"city"`
		Country  string `json:"country"`
		Postcode string `json:"postcode"`
	}

	Item struct {
		Name      string  `json:"name"`
		Quantity  float64 `json:"quantity"`
		UnitPrice float64 `json:"unitPrice"`
		Amount    float64 `json:"amount"`
	}

	Summary struct {
		Status     string  `json:"status"`
		TotalItems int     `json:"totalItems"`
		SubTotal   float64 `json:"subTotal"`
		Tax        float64 `json:"tax"`
		GrandTotal float64 `json:"grandTotal"`
	}
)
