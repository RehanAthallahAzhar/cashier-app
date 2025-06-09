package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        string         `gorm:"type:uuid" json:"id" validate:"required"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Name      string         `gorm:"type:varchar(100);unique" json:"name" validate:"required"`
	Price     float64        `json:"price" validate:"required"`
	Stock     int            `json:"stock" validate:"required"`
	Discount  float64        `json:"discount"`
	Type      string         `json:"type" validate:"required"`
	SellerID  string         `gorm:"type:uuid" json:"seller_id" validate:"required"`
}
