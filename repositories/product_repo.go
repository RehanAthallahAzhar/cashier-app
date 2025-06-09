package repositories

import (
	"context"
	"errors"
	"fmt"

	model "github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db}
}

func (p *ProductRepository) FindAllProducts(ctx context.Context) ([]model.Product, error) {
	results := []model.Product{}
	err := p.db.WithContext(ctx).Table("products").Select("*").Where("deleted_at is null").Find(&results).Error
	if err != nil {
		return []model.Product{}, err
	}
	return results, nil
}

func (p *ProductRepository) FindProductByID(ctx context.Context, id string) (model.Product, error) {
	var product model.Product
	result := p.db.WithContext(ctx).First(&product, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Product{}, model.ErrProductNotFound
		}
		return model.Product{}, fmt.Errorf("failed to read product by ID: %w", result.Error)
	}
	return product, nil
}

func (p *ProductRepository) FindProductBySellerID(ctx context.Context, sellerId string) ([]model.Product, error) {
	var product []model.Product
	result := p.db.WithContext(ctx).First(&product, "seller_id = ?", sellerId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []model.Product{}, model.ErrProductNotFound
		}
		return []model.Product{}, fmt.Errorf("failed to read product by ID: %w", result.Error)
	}
	return product, nil
}

func (p *ProductRepository) AddProduct(ctx context.Context, product *model.Product) error {
	if err := p.db.WithContext(ctx).Create(&product).Error; err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, id string, product *model.Product, sellerId string) error {
	err := p.db.WithContext(ctx).
		Table("products").
		Where("id = ? AND seller_id = ?", id, sellerId).
		Updates(product).Error

	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, id string, sellerId string) error {
	result := p.db.WithContext(ctx).
		Table("products").
		Where("id = ? AND seller_id = ?", id, sellerId).
		Delete(&model.Product{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("product not found or already deleted")
	}
	return nil
}
