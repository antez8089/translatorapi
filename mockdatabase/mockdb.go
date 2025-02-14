package mockdatabase

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"translatorapi/models"
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

	// Możesz dodać dodatkowe dane lub ustawienia, jeśli potrzebujesz

	return db, nil
}
