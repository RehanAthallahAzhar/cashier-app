package handlers

import "github.com/rehanazhar/cashier-app/repositories"

// API struct yang akan memegang dependensi repository
type API struct {
	// UserRepo    repositories.UserRepository
	// SessionRepo repositories.SessionsRepository
	ProductRepo repositories.ProductRepository
	CartRepo    repositories.CartRepository
}

// NewHandler adalah konstruktor untuk API
// Perbarui agar menerima semua repository
func NewHandler(
	// userRepo repositories.UserRepository,
	// sessionRepo repositories.SessionsRepository,
	productRepo repositories.ProductRepository,
	cartRepo repositories.CartRepository,
) *API {
	return &API{
		// UserRepo:    userRepo,
		// SessionRepo: sessionRepo,
		ProductRepo: productRepo,
		CartRepo:    cartRepo,
	}
}
