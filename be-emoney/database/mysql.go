package database

import (
	"fmt"
	"log"

	"emoney-2fa/config"
	"emoney-2fa/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.OTP{},
		&models.Account{},
		&models.Transaction{},
		&models.Product{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed products if table is empty
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count == 0 {
		products := []models.Product{
			{Name: "Nasi Goreng Spesial", Description: "Nasi goreng dengan telur, ayam, dan sayuran", Price: 25000, Stock: 50, Category: "Makanan", ImageURL: "https://via.placeholder.com/300", IsActive: true},
			{Name: "Mie Ayam Bakso", Description: "Mie ayam dengan bakso dan pangsit", Price: 20000, Stock: 40, Category: "Makanan", ImageURL: "https://via.placeholder.com/300", IsActive: true},
			{Name: "Es Teh Manis", Description: "Teh manis dingin segar", Price: 5000, Stock: 100, Category: "Minuman", ImageURL: "https://via.placeholder.com/300", IsActive: true},
			{Name: "Jus Jeruk Segar", Description: "Jus jeruk peras tanpa gula tambahan", Price: 15000, Stock: 30, Category: "Minuman", ImageURL: "https://via.placeholder.com/300", IsActive: true},
			{Name: "Sate Ayam (10 tusuk)", Description: "Sate ayam dengan bumbu kacang", Price: 30000, Stock: 25, Category: "Makanan", ImageURL: "https://via.placeholder.com/300", IsActive: true},
			{Name: "Bakso Urat", Description: "Bakso urat dengan kuah kaldu sapi", Price: 22000, Stock: 35, Category: "Makanan", ImageURL: "https://via.placeholder.com/300", IsActive: true},
			{Name: "Kopi Susu", Description: "Kopi susu dengan gula aren", Price: 18000, Stock: 60, Category: "Minuman", ImageURL: "https://via.placeholder.com/300", IsActive: true},
			{Name: "Pisang Goreng (5 pcs)", Description: "Pisang goreng crispy dengan coklat", Price: 12000, Stock: 45, Category: "Snack", ImageURL: "https://via.placeholder.com/300", IsActive: true},
		}
		db.Create(&products)
		log.Println("Seeded 8 products")
	}

	log.Println("MySQL connected and migrated")
	return db
}
