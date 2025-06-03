package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rehanazhar/cashier-app/handlers"
)

func InitRoutes(e *echo.Echo, api *handlers.API) {
	// Static files
	e.Static("/static", "template")

	productGroup := e.Group("/product")
	{
		productGroup.GET("/list", api.ProductList())
		productGroup.POST("/add", api.AddProduct)
		productGroup.DELETE("/delete", api.DeleteProduct)
		productGroup.PUT("/update", api.UpdateProduct)
	}

	cartGroup := e.Group("/cart")
	{
		cartGroup.GET("/list", api.CartList())
		cartGroup.POST("/add", api.AddCart)
		cartGroup.DELETE("/delete", api.DeleteCart)
		cartGroup.PUT("/update", api.UpdateCart)
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
