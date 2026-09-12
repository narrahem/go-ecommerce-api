package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/narrahem/go-ecommerce-api/pkg/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: No .env file found, reading system environment variables")
	}

	// 1. Establish DB connection pool
	dbCfg := database.Config{
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root123"),
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		DBName:   getEnv("DB_NAME", "ecommerce"),
	}

	db, err := database.Connect(dbCfg)
	if err != nil {
		log.Fatalf("Database Connection failed: %v", err)
	}
	defer db.Close()

	// 2. Automatically Run Migrations (like php artisan migrate)
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 3. Start Gin Router
	r := gin.Default()

	r.GET("/health	", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})

	port := getEnv("PORT", "8080")
	log.Printf("Server running on :%s", port)
	r.Run(":" + port)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
