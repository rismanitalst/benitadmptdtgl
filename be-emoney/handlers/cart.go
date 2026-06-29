package handlers

import (
	"net/http"

	"emoney-2fa/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	db *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{db: db}
}

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

// GET /v1/cart
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var items []models.CartItem
	if err := h.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil cart",
		})
		return
	}

	var total float64
	var itemCount int
	for _, item := range items {
		subtotal := item.Product.Price * float64(item.Quantity)
		total += subtotal
		itemCount += item.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items":      items,
			"total":      total,
			"item_count": itemCount,
		},
	})
}

// POST /v1/cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "product_id diperlukan",
		})
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// Check if product exists
	var product models.Product
	if err := h.db.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Produk tidak ditemukan",
		})
		return
	}

	// Check if already in cart
	var cartItem models.CartItem
	result := h.db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&cartItem)

	if result.Error == gorm.ErrRecordNotFound {
		// Add new item
		cartItem = models.CartItem{
			UserID:    userID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := h.db.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal menambahkan ke cart",
			})
			return
		}
	} else {
		// Update quantity
		cartItem.Quantity += req.Quantity
		if err := h.db.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal update cart",
			})
			return
		}
	}

	// Reload with product
	h.db.Preload("Product").First(&cartItem, cartItem.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil ditambahkan ke cart",
		"data":    cartItem,
	})
}

// DELETE /v1/cart/:id
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	itemID := c.Param("id")

	if err := h.db.Where("id = ? AND user_id = ?", itemID, userID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus item dari cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item berhasil dihapus dari cart",
	})
}

// DELETE /v1/cart (clear all items)
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengosongkan cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cart berhasil dikosongkan",
	})
}
