package mockdatabase

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"translatorapi/models"
	"fmt"
)

// MockDB ustawia mockową bazę danych z użyciem SQLite
func MockDB(t *testing.T) (*gorm.DB, error) {
	// Tworzymy bazę danych SQLite w pamięci
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Automatycznie migrujemy schemat (przykładowo tabela "words")
	// Upewnij się, że masz odpowiednie modele zdefiniowane w swojej aplikacji
	db.AutoMigrate(&models.Word{}, &models.Translation{}, &models.Example{})
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("Failed to enable foreign keys: %v", err)
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


	return db, nil
}
