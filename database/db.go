package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

func InitDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using system environment variables.")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool options
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB object: %w", err)
	}

	// Set the maximum number of open connections (e.g., 10)
	sqlDB.SetMaxOpenConns(1)
	// Set the maximum number of idle connections (e.g., 5)
	sqlDB.SetMaxIdleConns(5)
	// Set the maximum connection lifetime (e.g., 1 hour)
	sqlDB.SetConnMaxLifetime(0)


	fmt.Println("Database connected!")
	return db, nil
}
