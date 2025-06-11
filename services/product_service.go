package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/repositories"
)

// list function
type ProductService interface {
	FindAllProducts(ctx context.Context) ([]models.Product, error)
	FindProductBySellerID(ctx context.Context, sellerID string) ([]models.Product, error)
	AddProduct(ctx context.Context, userID, username, role string, productData *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, productID string, productData *models.Product, sellerID string) (*models.Product, error)
	DeleteProduct(ctx context.Context, productID string, sellerID string) error
}

type productServiceImpl struct {
	productRepo repositories.ProductRepository
}

var allowedRoles = map[string]bool{
	"superadmin": true,
	"admin":      true,
	"seller":     true,
	// Tambahkan peran lain jika diperlukan
}

// Dependency Injection
func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productServiceImpl{productRepo: productRepo}
}

func (s *productServiceImpl) FindAllProducts(ctx context.Context) ([]models.Product, error) {
	res, err := s.productRepo.FindAllProducts(ctx)
	if err != nil {
		// spesific error
		if errors.Is(err, models.ErrProductNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("service: failed to retrieve all products: %w", err)
	}
	if len(res) == 0 {
		return nil, models.ErrProductNotFound
	}

	return res, nil
}

func (s *productServiceImpl) FindProductBySellerID(ctx context.Context, sellerID string) ([]models.Product, error) {
	res, err := s.productRepo.FindProductBySellerID(ctx, sellerID)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("service: failed to retrieve products by seller ID %s: %w", sellerID, err)
	}
	if len(res) == 0 {
		return nil, models.ErrProductNotFound
	}

	return res, nil
}

func (s *productServiceImpl) AddProduct(ctx context.Context, userID, username, role string, productData *models.Product) (*models.Product, error) {
	if !allowedRoles[role] { // Memeriksa apakah 'role' ada di map 'allowedRoles'
		return nil, fmt.Errorf("role '%s' is not allowed to add products", role)
	}
	if productData.Name == "" || productData.Price <= 0 || productData.Stock <= 0 {
		return nil, fmt.Errorf("all required columns must not be empty and valid")
	}

	productData.SellerID = userID

	productData.ID = helpers.GenerateNewUserID()

	err := s.productRepo.AddProduct(ctx, productData)
	if err != nil {
		return nil, fmt.Errorf("service: failed to add product: %w", err)
	}

	return productData, nil
}

func (s *productServiceImpl) UpdateProduct(ctx context.Context, productID string, productData *models.Product, sellerID string) (*models.Product, error) {
	// Validasi input
	if productID == "" || productData.Name == "" || productData.Price <= 0 || productData.Stock < 0 || productData.Type == "" {
		return nil, fmt.Errorf("all required columns must not be empty and valid for update")
	}

	productData.ID = productID

	// Check if the product exists and belongs to the same seller id
	existingProduct, err := s.productRepo.FindProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return nil, models.ErrProductNotFound
		}

		return nil, fmt.Errorf("service: failed to find product for update: %w", err)
	}

	if existingProduct.SellerID != sellerID {
		return nil, fmt.Errorf("service: product does not belong to this seller")
	}

	err = s.productRepo.UpdateProduct(ctx, productID, productData, sellerID)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return nil, models.ErrProductNotFound
		}

		return nil, fmt.Errorf("service: failed to update product: %w", err)
	}

	return productData, nil
}

func (s *productServiceImpl) DeleteProduct(ctx context.Context, productID string, sellerID string) error {
	// Check if the product exists and belongs to the same seller
	existingProduct, err := s.productRepo.FindProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return models.ErrProductNotFound
		}

		return fmt.Errorf("service: failed to find product for deletion: %w", err)
	}

	if existingProduct.SellerID != sellerID {
		return fmt.Errorf("service: product does not belong to this seller")
	}

	err = s.productRepo.DeleteProduct(ctx, productID, sellerID)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return models.ErrProductNotFound
		}

		return fmt.Errorf("service: failed to delete product: %w", err)
	}

	return nil
}
