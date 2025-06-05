package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rehanazhar/shopeezy-inventory-cart/handlers"
)

func InitRoutes(e *echo.Echo, api *handlers.API) {
	// Static files
	e.Static("/static", "template")

	productGroup := e.Group("/product")
	{
		productGroup.GET("/list", api.ProductList())
		productGroup.POST("/add", api.AddProduct)
		productGroup.DELETE("/delete/:id", api.DeleteProduct)
		productGroup.PUT("/update/:id", api.UpdateProduct)
	}

	cartGroup := e.Group("/cart")
	{
		cartGroup.GET("/list", api.CartList())
		cartGroup.POST("/add/:id", api.AddCart)
		cartGroup.DELETE("/delete/:id/item/:product_id", api.DeleteCart)
		cartGroup.PUT("/update/:product_id", api.UpdateCart)
		// cartGroup.GET("/:id", api.GetCartItemByID)    // coming soon
	}

	// // Pages
	// e.GET("/", api.HomePage)
	// e.GET("/page/register", api.RegisterPage)
	// e.GET("/page/login", api.LoginPage)
	// e.GET("/page/dashboard", api.DashboardPage)

	// // Auth
	// e.POST("/user/register", h.Register)
	// e.POST("/user/login", h.Login)
	// e.GET("/user/session/valid", h.AuthMiddleware(h.SessionValid))
	// e.GET("/user/logout", h.AuthMiddleware(h.Logout))

	// // Profile
	// e.GET("/user/img/profile", h.AuthMiddleware(h.ImgProfileView))
	// e.POST("/user/img/update-profile", h.AuthMiddleware(h.ImgProfileUpdate))

}
