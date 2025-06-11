package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"
	"github.com/labstack/echo/v4"
)

// --- Implementasi Handler Cart ---

func (api *API) CartList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		userID := c.Get("user_id").(string) // Dapatkan userID dari context

		res, err := api.CartSvc.FindAllCarts(ctx, userID) // <-- Panggil service
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) { // Atau models.ErrCartEmpty jika Anda mendefinisikan
				return c.JSON(http.StatusOK, models.SuccessResponse{Message: "Your cart is still empty, let's go shopping"})
			}
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve cart items"})
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (api *API) AddCart(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	productID := c.Param("product_id") // Parameter dari URL
	userID := c.Get("user_id").(string)

	var req models.CartRequest // Asumsikan CartRequest memiliki Quantity dan Description
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
	}

	// Panggil service layer
	addedCartItem, err := api.CartSvc.AddCart(ctx, userID, productID, req.Quantity, req.Description) // <-- Panggil service
	if err != nil {
		switch {
		case errors.Is(err, models.ErrProductNotFound):
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		case err.Error() == fmt.Sprintf("insufficient stock: Only %s items available", strconv.Itoa(0)): // Perbaiki string perbandingan
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		default:
			c.Logger().Errorf("CartService.AddCart failed: %v", err)
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to add product to cart."})
		}
	}

	return c.JSON(http.StatusCreated, models.SuccessResponse{ // Changed to StatusCreated
		Username: username,
		Message:  "Product Successfully Added to Cart!",
		Data:     addedCartItem,
	})
}

func (api *API) DeleteCart(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	cartID := c.Param("id")
	productID := c.Param("product_id")

	// Panggil service layer
	err := api.CartSvc.DeleteCart(ctx, cartID, productID) // <-- Panggil service
	if err != nil {
		switch {
		case errors.Is(err, models.ErrCartItemNotFound):
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		default:
			c.Logger().Errorf("CartService.DeleteCart failed: %v", err)
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete product from cart."})
		}
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Product Successfully Deleted from Cart!"})
}

func (api *API) UpdateCart(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	cartID := c.Param("cart_id")

	var req models.CartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format or missing data"})
	}

	// Panggil service layer
	err := api.CartSvc.UpdateCart(ctx, cartID, req.Quantity) // <-- Panggil service
	if err != nil {
		switch {
		case errors.Is(err, models.ErrCartItemNotFound):
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Cart item for this product not found!"})
		case errors.Is(err, models.ErrInsufficientStock):
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		default:
			c.Logger().Errorf("CartService.UpdateCart failed: %v", err)
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update cart item."})
		}
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Cart item updated successfully!"})
}
