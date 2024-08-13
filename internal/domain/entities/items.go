package entities

import "time"

type Item struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Type      string     `db:"type" json:"type"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
