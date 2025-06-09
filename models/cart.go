package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	ProductID  string  `json:"product_id"`
	UserID     string  `json:"user_id"`
	Quantity   float64 `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

type JoinCart struct {
	Id         uint    `json:"id"`
	ProductId  uint    `json:"product_id"`
	Name       string  `json:"name"`
	Quantity   float64 `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

type CartRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}
