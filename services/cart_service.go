// File: cashier-app/services/cart_service.go
package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/repositories"
)

// CartService mendefinisikan interface untuk logika bisnis keranjang.
type CartService interface {
	FindAllCarts(ctx context.Context, userID string) ([]models.CartResponse, error)
	AddCart(ctx context.Context, userID, productID string, quantity int, description string) (*models.Cart, error)
	DeleteCart(ctx context.Context, cartID, productID string) error
	UpdateCart(ctx context.Context, cartID string, quantity int) error
}

type cartServiceImpl struct {
	cartRepo    repositories.CartRepository
	productRepo repositories.ProductRepository // Diperlukan untuk cek stok
}

func NewCartService(cartRepo repositories.CartRepository, productRepo repositories.ProductRepository) CartService {
	return &cartServiceImpl{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartServiceImpl) FindAllCarts(ctx context.Context, userID string) ([]models.CartResponse, error) {
	res, err := s.cartRepo.FindAllCarts(ctx, userID)

	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("service: failed to retrieve all carts for user %s: %w", userID, err)
	}
	if len(res) == 0 {
		return nil, models.ErrCartItemNotFound
	}

	return res, nil
}

func (s *cartServiceImpl) AddCart(ctx context.Context, userID, productID string, quantity int, description string) (*models.Cart, error) {
	if productID == "" || quantity <= 0 {
		return nil, fmt.Errorf("product ID and quantity are required and must be valid")
	}

	product, err := s.productRepo.FindProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return nil, models.ErrProductNotFound
		}

		return nil, fmt.Errorf("service: failed to retrieve product details for product ID %s: %w", productID, err)
	}

	if product.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock: Only %s items available", strconv.Itoa(product.Stock))
	}

	newCartID := helpers.GenerateNewUserID()

	err = s.cartRepo.AddCart(ctx, product, quantity, description, productID, newCartID, userID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to add product to cart: %w", err)
	}

	// Asumsikan cartRepo.AddCart memodifikasi atau mengembalikan cart item yang ditambahkan
	// Jika tidak, Anda mungkin perlu mengambilnya dari repo lagi atau membuat struct di sini.
	return &models.Cart{
		ID:          newCartID,
		ProductID:   productID,
		UserID:      userID,
		Quantity:    float64(quantity),
		Description: description,
	}, nil
}

func (s *cartServiceImpl) DeleteCart(ctx context.Context, cartID, productID string) error {
	if cartID == "" || productID == "" {
		return fmt.Errorf("cart ID and product ID are required for deletion")
	}

	err := s.cartRepo.DeleteCart(ctx, cartID, productID)
	if err != nil {
		if errors.Is(err, models.ErrCartItemNotFound) {
			return models.ErrCartItemNotFound
		}

		return fmt.Errorf("service: failed to delete cart item %s for product %s: %w", cartID, productID, err)
	}

	return nil
}

func (s *cartServiceImpl) UpdateCart(ctx context.Context, cartID string, quantity int) error {
	if cartID == "" {
		return fmt.Errorf("cart ID is required for update")
	}
	if quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	err := s.cartRepo.UpdateCart(ctx, cartID, quantity)
	if err != nil {
		if errors.Is(err, models.ErrCartItemNotFound) {
			return models.ErrCartItemNotFound
		}
		if errors.Is(err, models.ErrInsufficientStock) {
			return models.ErrInsufficientStock
		}

		return fmt.Errorf("service: failed to update cart item %s: %w", cartID, err)
	}

	return nil
}
