package entities

import "time"

type Customer struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Address   string     `db:"address" json:"address"`
	City      string     `db:"city" json:"city"`
	Postcode  string     `db:"postcode" json:"postcode"`
	Country   string     `db:"country" json:"country"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
