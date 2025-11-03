package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

var db *gorm.DB

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	fmt.Println("---db", os.Getenv("DB_USER"))
	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	// Retry database connection
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database. Retrying in 2 seconds... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after 30 attempts:", err)
	}

	// Auto migrate
	db.AutoMigrate(&Post{})

	// Echo setup
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", healthCheck)
	e.GET("/api/posts", getPosts)
	e.GET("/api/posts/:id", getPost)
	e.POST("/api/posts", createPost)
	e.PUT("/api/posts/:id", updatePost)
	e.DELETE("/api/posts/:id", deletePost)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func getPosts(c echo.Context) error {
	var posts []Post
	result := db.Order("created_at desc").Find(&posts)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}
	return c.JSON(http.StatusOK, posts)
}

func getPost(c echo.Context) error {
	id := c.Param("id")
	var post Post
	result := db.First(&post, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Post not found",
		})
	}
	return c.JSON(http.StatusOK, post)
}

func createPost(c echo.Context) error {
	post := new(Post)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	result := db.Create(post)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(http.StatusCreated, post)
}

func updatePost(c echo.Context) error {
	id := c.Param("id")
	var post Post
	if result := db.First(&post, id); result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Post not found",
		})
	}

	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	db.Save(&post)
	return c.JSON(http.StatusOK, post)
}

func deletePost(c echo.Context) error {
	id := c.Param("id")
	result := db.Delete(&Post{}, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Post deleted successfully",
	})
}
