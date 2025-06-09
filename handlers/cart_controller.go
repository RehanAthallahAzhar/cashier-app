package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"
	"github.com/labstack/echo/v4"
)

func (api *API) CartList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		userId := c.Get("userID").(string)

		res, err := api.CartRepo.FindAllCarts(ctx, userId)
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) {
				return c.JSON(http.StatusOK, models.SuccessResponse{Message: "Your cart is still empty, let's go shopping"})
			}
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve products"})
		}

		if len(res) == 0 {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Cart list not found!"})
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

	product_id := c.Param("product_id")

	userId := c.Get("userID").(string)
	newCartId := helpers.GenerateNewUserID()

	var req models.CartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
	}

	if product_id == "" || req.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Product ID and Quantity are required and must be valid."})
	}

	product, err := api.ProductRepo.FindProductByID(ctx, product_id)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve product details."})
	}

	// validation of available stock
	if product.Stock < req.Quantity {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Insufficient stock: Only " + strconv.Itoa(product.Stock) + " items available."})
	}

	err = api.CartRepo.AddCart(ctx, product, req.Quantity, newCartId, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to add product to cart."})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Product Successfully Added to Cart!"})
}

func (api *API) DeleteCart(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	id := c.Param("id")
	product_id := c.Param("product_id")

	idInt, _ := strconv.ParseUint(id, 10, 32)

	if product_id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Product ID is required."})
	}

	err := api.CartRepo.DeleteCart(ctx, uint(idInt), product_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete product"})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Product Successfully Delete to Cart!"})
}

func (api *API) UpdateCart(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	product_id := c.Param("product_id")

	// Bind request
	var req models.CartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format or missing data"})
	}

	if product_id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Product ID is required."})
	}
	if req.Quantity < 0 { // Defense in Depth
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Quantity cannot be negative."})
	}

	err := api.CartRepo.UpdateCart(ctx, product_id, req.Quantity)
	if err != nil {
		if errors.Is(err, models.ErrCartItemNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Cart item for this product not found!"})
		}
		if errors.Is(err, models.ErrInsufficientStock) {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update cart item."})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Cart item updated successfully!"})
}
