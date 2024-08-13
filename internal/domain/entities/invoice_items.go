package entities

import "time"

type InvoiceItems struct {
	ID        int        `json:"id" db:"id"`
	InvoiceID int        `json:"invoice_id" db:"invoice_id"`
	ItemID    int        `json:"item_id" db:"item_id"`
	Quantity  float64    `json:"quantity" db:"quantity"`
	UnitPrice float64    `json:"unit_price" db:"unit_price"`
	Amount    float64    `json:"amount" db:"amount"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
