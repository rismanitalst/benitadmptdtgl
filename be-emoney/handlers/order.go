package handlers

import (
	"net/http"
	"strconv"

	"emoney-2fa/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

type CheckoutRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	Notes           string `json:"notes"`
	PaymentMethod   string `json:"payment_method" binding:"required"`
}

// POST /v1/orders/checkout
func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "shipping_address dan payment_method diperlukan",
		})
		return
	}

	// Get cart items
	var cartItems []models.CartItem
	if err := h.db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil cart",
		})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Cart kosong",
		})
		return
	}

	// Calculate total and create order items
	var total float64
	var orderItems []models.OrderItem
	for _, item := range cartItems {
		subtotal := item.Product.Price * float64(item.Quantity)
		total += subtotal
		orderItems = append(orderItems, models.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Price:       item.Product.Price,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Create order
	order := models.Order{
		UserID:          userID,
		TotalAmount:     total,
		Status:          "pending",
		ShippingAddress: req.ShippingAddress,
		Notes:           req.Notes,
		PaymentMethod:   req.PaymentMethod,
		Items:           orderItems,
	}

	if err := h.db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat order",
		})
		return
	}

	// Clear cart
	h.db.Where("user_id = ?", userID).Delete(&models.CartItem{})

	// Reload with items
	h.db.Preload("Items").First(&order, order.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order berhasil dibuat",
		"data":    order,
	})
}

// GET /v1/orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var orders []models.Order
	if err := h.db.Preload("Items").Where("user_id = ?", userID).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
	})
}

// GET /v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID order tidak valid",
		})
		return
	}

	var order models.Order
	if err := h.db.Preload("Items").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Order tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}
