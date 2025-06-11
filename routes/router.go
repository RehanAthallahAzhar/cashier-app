package routes

import (
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/handlers"    // Sesuaikan import path handler Anda
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/middlewares" // Impor middlewares

	// Impor authclient
	"github.com/labstack/echo/v4"
)

// InitRoutes -> menginisialisasi semua rute API
func InitRoutes(e *echo.Echo, api *handlers.API, authMiddlewareOpts middlewares.AuthMiddlewareOptions) {
	// Static files
	e.Static("/static", "template")

	productGroup := e.Group("/product")
	productGroup.Use(middlewares.AuthMiddleware(authMiddlewareOpts)) // using middleware
	{
		productGroup.GET("/list", api.ProductList())
		productGroup.GET("/mylist", api.SellerProductList())
		productGroup.POST("/add", api.AddProduct)
		productGroup.DELETE("/delete/:id", api.DeleteProduct)
		productGroup.PUT("/update/:id", api.UpdateProduct)
	}

	cartGroup := e.Group("/cart")
	cartGroup.Use(middlewares.AuthMiddleware(authMiddlewareOpts)) // using middleware
	{
		cartGroup.GET("/list", api.CartList())
		cartGroup.POST("/add/:product_id", api.AddCart)
		cartGroup.DELETE("/delete/:id/item/:product_id", api.DeleteCart)
		cartGroup.PUT("/update/:cart_id", api.UpdateCart)
	}

}
