package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rehanazhar/shopeezy-inventory-cart/models"
)

func (api *API) CartList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		res, err := api.CartRepo.FindAllCarts(ctx)
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) {
				return c.JSON(http.StatusOK, models.SuccessResponse{Message: "Your cart is still empty, let's go shopping"})
			}
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to retrieve products"})
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

	id := c.Param("id")

	idint, _ := strconv.ParseUint(id, 10, 32)

	var req models.CartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Invalid JSON format"})
	}

	if idint == 0 || req.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Product ID and Quantity are required and must be valid."})
	}

	product, err := api.ProductRepo.FindProductByID(ctx, uint(idint))
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to retrieve product details."})
	}

	// validation of available stock
	if product.Stock < req.Quantity {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Insufficient stock: Only " + strconv.Itoa(product.Stock) + " items available."})
	}

	err = api.CartRepo.AddCart(ctx, product, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to add product to cart."})
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
	productIdInt, _ := strconv.ParseUint(product_id, 10, 32)

	if productIdInt == 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Product ID is required."})
	}

	err := api.CartRepo.DeleteCart(ctx, uint(idInt), uint(productIdInt))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to delete product"})
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
	productIdInt, _ := strconv.ParseUint(product_id, 10, 32)

	// Bind request
	var req models.CartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Invalid JSON format or missing data"})
	}

	if productIdInt == 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Product ID is required."})
	}
	if req.Quantity < 0 { // Defense in Depth
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Quantity cannot be negative."})
	}

	err := api.CartRepo.UpdateCart(ctx, uint(productIdInt), req.Quantity)
	if err != nil {
		if errors.Is(err, models.ErrCartItemNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Cart item for this product not found!"})
		}
		if errors.Is(err, models.ErrInsufficientStock) {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to update cart item."})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Cart item updated successfully!"})
}
