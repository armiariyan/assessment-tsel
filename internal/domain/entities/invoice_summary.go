package entities

import "time"

type InvoiceSummary struct {
	ID         int        `json:"id" db:"id"`
	InvoiceID  int        `json:"invoice_id" db:"invoice_id"`
	TotalItems int        `json:"total_items" db:"total_items"`
	Subtotal   float64    `json:"subtotal" db:"subtotal"`
	Tax        float64    `json:"tax" db:"tax"`
	GrandTotal float64    `json:"grand_total" db:"grand_total"`
	IsPaid     bool       `json:"is_paid" db:"is_paid"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
