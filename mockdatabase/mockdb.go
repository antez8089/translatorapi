package mockdatabase

import (
	"fmt"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MockDB(t *testing.T) (*gorm.DB, error) {
	// Connect to the test database in Docker
	dsn := "host=localhost user=test_user password=test_pass dbname=translatorapi_test port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}
	sqlSchema := `
	CREATE TABLE IF NOT EXISTS words (
    id SERIAL PRIMARY KEY,
    polish_word VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS translations (
		id SERIAL PRIMARY KEY,
		word_id INT REFERENCES words(id) ON DELETE CASCADE,
		english_word VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS examples (
		id SERIAL PRIMARY KEY,
		translation_id INT REFERENCES translations(id) ON DELETE CASCADE,
		sentence TEXT NOT NULL
	);
	`

	if err := db.Exec(sqlSchema).Error; err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}
	// Set connection pool options
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB object: %w", err)
	}
	// Set the maximum number of open connections (e.g., 10)
	sqlDB.SetMaxOpenConns(5)
	// Set the maximum number of idle connections (e.g., 5)
	sqlDB.SetMaxIdleConns(5)
	// Set the maximum connection lifetime (e.g., 1 hour)
	sqlDB.SetConnMaxLifetime(0)

	// Clean database before tests
	db.Exec("SET CONSTRAINTS ALL IMMEDIATE;")
	db.Exec("TRUNCATE words, translations, examples RESTART IDENTITY CASCADE;")

	return db, nil
}
