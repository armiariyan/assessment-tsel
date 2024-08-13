package entities

import (
	"time"
)

type Invoice struct {
	ID         int        `json:"id" db:"pk_invoice_id"`
	InvoiceID  string     `json:"invoiceId" db:"uq_invoice_id"`
	IssueDate  time.Time  `json:"issueDate" db:"issue_date"`
	DueDate    time.Time  `json:"dueDate" db:"due_date"`
	Subject    string     `json:"subject" db:"subject"`
	CustomerID int        `json:"customerId" db:"customer_id"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`

	// * relations
	InvoiceSummary `json:"-"`
	InvoiceItems   `json:"-"`
	Customer       `json:"-"`
}

type InvoiceListParams struct {
	Limit      uint
	Page       uint
	InvoiceID  string
	Subject    string
	Customer   string
	IssueDate  string
	DueDate    string
	Status     string
	TotalItems int
}
