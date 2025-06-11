package handlers

import (
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/pkg/authclient"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/services"
)

// API struct untuk mengelola handler dan dependensinya.
type API struct {
	ProductRepo repositories.ProductRepository
	CartRepo    repositories.CartRepository // Tetap butuh ini untuk NewHandler
	AuthClient  *authclient.AuthClient
	ProductSvc  services.ProductService
	CartSvc     services.CartService // <-- TAMBAHKAN INI!
}

// NewHandler membuat instance baru dari API.
func NewHandler(productRepo repositories.ProductRepository, cartRepo repositories.CartRepository, authClient *authclient.AuthClient, productSvc services.ProductService, cartSvc services.CartService) *API { // <-- SESUAIKAN SIGNATURE
	return &API{
		ProductRepo: productRepo,
		CartRepo:    cartRepo,
		AuthClient:  authClient,
		ProductSvc:  productSvc,
		CartSvc:     cartSvc, // <-- INISIALISASI INI!
	}
}
