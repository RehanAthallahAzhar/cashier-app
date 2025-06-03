package handlers // Nama package harus 'handlers'

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rehanazhar/cashier-app/models"
)

func (api *API) ProductList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		res, err := api.ProductRepo.ReadProducts(ctx)
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

	var req models.ProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Invalid JSON format"})
	}

	if req.Id == 0 || req.Name == "" || req.Price <= 0 || req.Stock < 0 || req.Type == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: All required coloumn must not be empty and valid"})
	}

	err := api.ProductRepo.UpdateProduct(ctx, uint(req.Id), req)
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

	var req models.DeleteProduct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Invalid JSON format"})
	}

	err := api.ProductRepo.DeleteProduct(ctx, uint(req.Id))
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to delete product"})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Delete Success"})
}
