package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"

	"github.com/labstack/echo/v4"
)

func (api *API) ProductList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		c.Logger().Infof("Received request for product list from IP: %s", c.RealIP())

		res, err := api.ProductSvc.FindAllProducts(ctx)
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) {

				return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			}

			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve products"})
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (api *API) SellerProductList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		sellerId := c.Get("user_id").(string)

		res, err := api.ProductSvc.FindProductBySellerID(ctx, sellerId)
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) {

				return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			}

			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve products"})
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

	userID := c.Get("user_id").(string)

	role, ok := c.Get("role").(string)
	if !ok || role == "" {
		c.Logger().Warnf("Role not found in context for user %s (ID: %s), defaulting to 'guest'", username, userID)
		role = "guest" // Default role to prevent errors
	}

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
	}

	product.ID = helpers.GenerateNewUserID()

	addedProduct, err := api.ProductSvc.AddProduct(ctx, userID, username, role, &product)
	if err != nil {
		// spesific errors
		switch err.Error() {
		case fmt.Sprintf("role '%s' is not allowed to add products", role):
			return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: err.Error()})
		case "all required columns must not be empty and valid":
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		default:
			c.Logger().Errorf("ProductService.AddProduct failed: %v", err)
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal Server Error: Failed to add product"})
		}
	}

	return c.JSON(http.StatusCreated, models.SuccessResponse{
		Username: username,
		Message:  "Product Added Successfully!",
		Data:     addedProduct,
	})
}

func (api *API) UpdateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	productID := c.Param("id")
	sellerID := c.Get("user_id").(string)

	var productData models.Product
	if err := c.Bind(&productData); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
	}

	updatedProduct, err := api.ProductSvc.UpdateProduct(ctx, productID, &productData, sellerID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrProductNotFound):

			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		case err.Error() == "all required columns must not be empty and valid for update":

			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		case err.Error() == "service: product does not belong to this seller":

			return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: err.Error()})
		default:
			c.Logger().Errorf("ProductService.UpdateProduct failed: %v", err)

			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update product"})
		}
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Product updated successfully", Data: updatedProduct})
}

func (api *API) DeleteProduct(c echo.Context) error {
	ctx := c.Request().Context()

	username := ""
	if val := c.Get("username"); val != nil {
		if u, ok := val.(string); ok {
			username = u
		}
	}

	productID := c.Param("id")
	sellerID := c.Get("user_id")

	err := api.ProductSvc.DeleteProduct(ctx, productID, sellerID.(string))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrProductNotFound):
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		case err.Error() == "service: product does not belong to this seller":
			return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: err.Error()})
		default:
			c.Logger().Errorf("ProductService.DeleteProduct failed: %v", err)
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete product"})
		}
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{Username: username, Message: "Product deleted successfully"})
}
