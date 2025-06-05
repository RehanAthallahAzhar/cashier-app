package repositories

import (
	"context"
	"errors"
	"fmt"

	model "github.com/rehanazhar/cashier-app/models"
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

func (p *ProductRepository) FindProductByID(ctx context.Context, id uint) (model.Product, error) {
	var product model.Product
	result := p.db.WithContext(ctx).First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Product{}, model.ErrProductNotFound
		}
		return model.Product{}, fmt.Errorf("failed to read product by ID: %w", result.Error)
	}
	return product, nil
}

func (p *ProductRepository) AddProduct(ctx context.Context, product model.Product) error {
	if err := p.db.WithContext(ctx).Create(&product).Error; err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, id uint, product *model.Product) error {
	err := p.db.WithContext(ctx).Table("products").Where("id = ?", id).Updates(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	result := p.db.WithContext(ctx).Delete(&model.Product{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("product not found or already deleted")
	}
	return nil
}
