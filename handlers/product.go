package handlers // Nama package harus 'handlers'

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rehanazhar/shopeezy-inventory-cart/models"
)

func (api *API) ProductList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		res, err := api.ProductRepo.FindAllProducts(ctx)
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) {
				return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to retrieve products"})
		}

		if len(res) == 0 {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product list not found!"})
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (api *API) AddProduct(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Invalid JSON format"})
	}

	if product.Name == "" || product.Price <= 0 || product.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Name, Price, and Stock are required and must be valid."})
	}

	err := api.ProductRepo.AddProduct(ctx, product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to add product"})
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

	idint, _ := strconv.ParseUint(id, 10, 32)

	product, err := api.ProductRepo.FindProductByID(ctx, uint(idint))
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
	}

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Invalid JSON format"})
	}

	if product.ID == 0 || product.Name == "" || product.Price <= 0 || product.Stock < 0 || product.Type == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: All required coloumn must not be empty and valid"})
	}

	err = api.ProductRepo.UpdateProduct(ctx, uint(product.ID), &product)
	if err != nil {

		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to update product"})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Update Success"})
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

	idint, _ := strconv.ParseUint(id, 10, 32)

	err := api.ProductRepo.DeleteProduct(ctx, uint(idint))
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to delete product"})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Delete Success"})
}
