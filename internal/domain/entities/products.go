package entities

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"column:id" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	Price       float64        `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Stock       float64        `gorm:"column:stock;type:int;not null" json:"stock"`
	Rating      *float64       `gorm:"column:rating;type:decimal(2,1);" json:"rating"`
	Variety     datatypes.JSON `gorm:"column:variety;type:jsonb;default:'[]'" json:"variety"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp with time zone;default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp with time zone;default:current_timestamp;autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;type:timestamp with time zone" json:"-"`
}
