package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name     string  `gorm:"type:varchar(100);unique_index" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Stock    int     `json:"stock" validate:"required"`
	Discount float64 `json:"discount"`
	Type     string  `json:"type" validate:"required"`
}
