package handlers

import (
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/pkg/authclient"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/repositories"
)

// API struct yang akan memegang dependensi repository
type API struct {
	// UserRepo    repositories.UserRepository
	// SessionRepo repositories.SessionsRepository
	ProductRepo repositories.ProductRepository
	CartRepo    repositories.CartRepository
	AuthClient  *authclient.AuthClient
}

// NewHandler adalah konstruktor untuk API
// Perbarui agar menerima semua repository
func NewHandler(
	// userRepo repositories.UserRepository,
	// sessionRepo repositories.SessionsRepository,
	productRepo repositories.ProductRepository,
	cartRepo repositories.CartRepository,
	authClient *authclient.AuthClient,
) *API {
	return &API{
		// UserRepo:    userRepo,
		// SessionRepo: sessionRepo,
		ProductRepo: productRepo,
		CartRepo:    cartRepo,
		AuthClient:  authClient,
	}
}
