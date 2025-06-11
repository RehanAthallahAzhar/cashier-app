package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID          string  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID   string  `gorm:"type:uuid;not null" json:"product_id"`
	UserID      string  `gorm:"type:uuid;not null" json:"user_id"`
	Quantity    float64 `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
	Description string  `json:"description"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

type CartResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	UserID      string  `json:"user_id"`
	ProductName string  `json:"name"`
	Quantity    float64 `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
	Description string  `json:"description"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty"`
}

type CartRequest struct {
	Quantity    int    `json:"quantity" validate:"required,min=1"`
	Description string `json:"description"`
}
