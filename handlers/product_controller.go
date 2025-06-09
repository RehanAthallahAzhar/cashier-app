package handlers // Nama package harus 'handlers'

import (
	"errors"
	"net/http"

	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"
	"github.com/labstack/echo/v4"
)

func (api *API) ProductList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		res, err := api.ProductRepo.FindAllProducts(ctx)
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) {
				return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve products"})
		}

		if len(res) == 0 {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product list not found!"})
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (api *API) SellerProductList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		sellerId := c.Get("userID").(string)

		res, err := api.ProductRepo.FindProductBySellerID(ctx, sellerId)
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) {
				return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve products"})
		}

		if len(res) == 0 {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product list not found!"})
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (api *API) AddProduct(c echo.Context) error {
	ctx := c.Request().Context()

	var username string
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	userId := c.Get("userID").(string)

	role, ok := c.Get("role").(string)
	if !ok || role == "" {
		c.Logger().Warnf("Role not found in context for user %s (ID: %s), defaulting to 'guest'", username, userId)
		role = "guest" // Default atau handle sebagai error jika role adalah mandatory
	}

	if role == "users" || role == "guest" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Not allowed to add products"})
	}

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
	}

	product.SellerID = userId

	newProductId := helpers.GenerateNewUserID()
	product.ID = newProductId

	if product.Name == "" || product.Price <= 0 || product.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "All required coloumn must not be empty and valid"})
	}

	err := api.ProductRepo.AddProduct(ctx, &product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to add product"})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Product Added Successfully!"})
}

func (api *API) UpdateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	id := c.Param("id")
	sellerId := c.Get("userID").(string)

	product, err := api.ProductRepo.FindProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
	}

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
	}

	if product.ID == "" || product.Name == "" || product.Price <= 0 || product.Stock < 0 || product.Type == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "All required coloumn must not be empty and valid"})
	}

	err = api.ProductRepo.UpdateProduct(ctx, id, &product, sellerId)
	if err != nil {

		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update product"})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Product updated successfully"})
}

func (api *API) DeleteProduct(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	id := c.Param("id")
	sellerId := c.Get("userID").(string)

	err := api.ProductRepo.DeleteProduct(ctx, id, sellerId)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete product"})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Product deleted successfully"})
}
