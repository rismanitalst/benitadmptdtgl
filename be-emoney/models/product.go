package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(255);not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(15,2);not null" json:"price"`
	Stock       int     `gorm:"default:0" json:"stock"`
	Category    string  `gorm:"type:varchar(100)" json:"category"`
	ImageURL    string  `gorm:"type:text" json:"image_url"`
	IsActive    bool    `gorm:"default:true" json:"is_active"`
}
