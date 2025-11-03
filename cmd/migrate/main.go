package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/akifumiwatanabe/blog-app/migrations"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text"`
	ImageURL  string    `json:"image_url" gorm:"type:varchar(500)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "blog_user"),
		getEnv("DB_PASSWORD", "blog_password"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "blog_db"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	log.Println("Running migrations...")
	if err := db.AutoMigrate(&Post{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Migrations completed successfully")

	// Seed data
	log.Println("Seeding data...")
	if err := migrations.SeedData(db); err != nil {
		log.Fatal("Failed to seed data:", err)
	}
	log.Println("Data seeding completed successfully")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
