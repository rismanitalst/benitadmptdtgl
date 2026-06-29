package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID          uint        `gorm:"not null;index" json:"user_id"`
	TotalAmount     float64     `gorm:"type:decimal(15,2)" json:"total_amount"`
	Status          string      `gorm:"type:varchar(50);default:pending" json:"status"`
	ShippingAddress string      `gorm:"type:text" json:"shipping_address"`
	Notes           string      `gorm:"type:text" json:"notes"`
	PaymentMethod   string      `gorm:"type:varchar(50)" json:"payment_method"`
	Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	User            User        `gorm:"foreignKey:UserID" json:"-"`
}

type OrderItem struct {
	gorm.Model
	OrderID     uint    `gorm:"not null;index" json:"order_id"`
	ProductID   uint    `gorm:"not null" json:"product_id"`
	ProductName string  `gorm:"type:varchar(255)" json:"product_name"`
	Price       float64 `gorm:"type:decimal(15,2)" json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `gorm:"type:decimal(15,2)" json:"subtotal"`
	Order       Order   `gorm:"foreignKey:OrderID" json:"-"`
}
