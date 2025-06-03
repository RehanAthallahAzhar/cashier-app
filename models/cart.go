package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	ProductID  uint    `json:"product_id"`
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
	ProductID uint `json:"product_id" validate:"required,min=1"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

type DeleteCartRequest struct {
	Id        uint `json:"id" validate:"required,min=1"`
	ProductID uint `json:"product_id" validate:"required,min=1"`
}
