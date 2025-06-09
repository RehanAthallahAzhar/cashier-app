package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"
)

type CartRepository struct {
	db          *gorm.DB
	ProductRepo ProductRepository // karena cart membutuhkan fungsi product
}

func NewCartRepository(db *gorm.DB, productRepo ProductRepository) CartRepository {
	return CartRepository{db: db, ProductRepo: productRepo}
}

func (c *CartRepository) FindAllCarts(ctx context.Context, id string) ([]models.JoinCart, error) {
	var listCart []models.JoinCart
	err := c.db.WithContext(ctx).Table("carts").
		Select("carts.id, carts.product_id, products.name, carts.quantity, carts.total_price").
		Joins("left join products on products.id = carts.product_id").Where("carts.deleted_at is NULL AND user_id = ?", id).
		Scan(&listCart).Error

	return listCart, err
}

func (c *CartRepository) FindCartItemByID(ctx context.Context, cartItemID uint) (models.JoinCart, error) {
	var cartItem models.JoinCart
	err := c.db.WithContext(ctx).Table("carts").
		Select("carts.id, carts.product_id, products.name, carts.quantity, carts.total_price").
		Joins("left join products on products.id = carts.product_id").
		Where("carts.id = ?", cartItemID).
		Where("carts.deleted_at is NULL").
		Scan(&cartItem).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.JoinCart{}, models.ErrCartItemNotFound
		}
		return models.JoinCart{}, fmt.Errorf("failed to retrieve cart item: %w", err)
	}

	if cartItem.Id == 0 {
		return models.JoinCart{}, models.ErrCartItemNotFound
	}

	return cartItem, nil
}

func (c *CartRepository) AddCart(ctx context.Context, product models.Product, quantity int, newCartId string, userId string) error {
	var cart models.Cart

	// Search if the item is already in the cart
	cartExistErr := c.db.WithContext(ctx).First(&cart, "product_id = ?", product.ID).Error

	// Use transactions to ensure atomicity
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// If the product is not in your cart (gorm.ErrRecordNotFound)
		if cartExistErr == gorm.ErrRecordNotFound {
			totalPriceForNewItem := product.Price * float64(quantity) * (1 - (product.Discount / 100))

			var newCart = &models.Cart{
				ProductID:  newCartId,
				UserID:     userId,
				Quantity:   float64(quantity),
				TotalPrice: totalPriceForNewItem,
			}

			// new cart entry
			err := tx.WithContext(ctx).Create(newCart).Error
			if err != nil {
				return err // create an error to rollback
			}
		} else if cartExistErr != nil {
			// If there are other errors when searching for the basket (selain ErrRecordNotFound)
			return cartExistErr
		} else {
			// If the product is already in the cart, update the quantity and total price
			totalPriceForAddedItems := product.Price * float64(quantity) * (1 - (product.Discount / 100))

			// Update cart quantity and total price
			err := tx.WithContext(ctx).Model(&models.Cart{}).Where("product_id = ?", product.ID).
				Updates(map[string]interface{}{
					"quantity":    gorm.Expr("quantity + ?", quantity),
					"total_price": gorm.Expr("total_price + ?", totalPriceForAddedItems),
				}).Error
			if err != nil {
				return err // Rollback
			}
		}

		// Reduce product stock in the products table
		// Use gorm.Expr for safe arithmetic operations
		err := tx.WithContext(ctx).Model(&models.Product{}).Where("id = ?", product.ID).Update("stock", gorm.Expr("stock - ?", quantity)).Error
		if err != nil {
			return err // Rollback
		}

		return nil // Commit
	})
}

func (c *CartRepository) deleteCartItemInTransaction(ctx context.Context, tx *gorm.DB, cartItemID uint, productID string) error {
	var slctedCart models.Cart
	err := tx.WithContext(ctx).Table("carts").Select("*").Where("id = ?", cartItemID).Where("product_id = ?", productID).Scan(&slctedCart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrCartItemNotFound
		}
		return err
	}

	var slctedProduct models.Product
	err = tx.WithContext(ctx).Table("products").Select("*").Where("id = ?", productID).Scan(&slctedProduct).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Log warning atau info di sini jika produk tidak ditemukan tapi tetap hapus item keranjang
			log.Printf("WARNING: Product with ID %s associated with cart item %d not found in products table. Proceeding with cart item deletion.", productID, cartItemID)
		}
		return err // Tetap kembalikan error jika gagal mengambil produk (selain not found)
	}

	// Kembalikan stok produk
	err = tx.WithContext(ctx).Table("products").Where("id = ?", productID).Update("stock", gorm.Expr("stock + ?", slctedCart.Quantity)).Error
	if err != nil {
		return err
	}

	// Hapus item dari keranjang
	result := tx.WithContext(ctx).Table("carts").Where("id = ?", cartItemID).Delete(&models.Cart{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("failed to delete cart item: no rows affected")
	}

	return nil
}

// Internal helper function to delete cart items in transactions
func (c *CartRepository) DeleteCart(ctx context.Context, cartItemID uint, productID string) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return c.deleteCartItemInTransaction(ctx, tx, cartItemID, productID)
	})
}

func (c *CartRepository) UpdateCart(ctx context.Context, productID string, newQuantity int) error {
	if newQuantity < 0 {
		return errors.New("new quantity cannot be negative")
	}

	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingCartItem models.Cart
		err := tx.First(&existingCartItem, "product_id = ?", productID).Error
		if err != nil {
			// if cart not found
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrCartItemNotFound
			}
			return fmt.Errorf("failed to retrieve cart item by product ID: %w", err)
		}

		// Jika newQuantity adalah 0, ini berarti item dihapus dari keranjang
		if newQuantity == 0 {
			return c.deleteCartItemInTransaction(ctx, tx, existingCartItem.ID, existingCartItem.ProductID)
		}

		var product models.Product
		product, err = c.ProductRepo.FindProductByID(ctx, existingCartItem.ProductID)
		if err != nil {
			return fmt.Errorf("failed to retrieve associated product details: %w", err)
		}

		oldQuantity := int(existingCartItem.Quantity)
		quantityChange := newQuantity - oldQuantity

		// 3. Validasi stok (hanya jika menambah kuantitas)
		if quantityChange > 0 {
			if product.Stock < quantityChange {
				return models.ErrInsufficientStock
			}
		}

		err = tx.Model(&models.Product{}).Where("id = ?", product.ID).
			Update("stock", gorm.Expr("stock - ?", quantityChange)).Error
		if err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}

		// 5Hitung ulang TotalPrice untuk item keranjang
		newTotalPrice := product.Price * float64(newQuantity) * (1 - (product.Discount / 100))

		// 6. Perbarui item keranjang (kuantitas dan total harga)
		err = tx.Model(&models.Cart{}).Where("id = ?", existingCartItem.ID). // Update berdasarkan ID item keranjang yang ditemukan
											Updates(map[string]interface{}{
				"quantity":    float64(newQuantity),
				"total_price": newTotalPrice,
			}).Error
		if err != nil {
			return fmt.Errorf("failed to update cart item: %w", err)
		}

		return nil // Commit transaksi
	})
}
