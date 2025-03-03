package models


// Example represents an example sentence using a translation
type Example struct {
	ID         uint   `gorm:"primaryKey"`
	TranslationID uint        `gorm:"not null;index"`
	Sentence      string      `gorm:"unique;not null"`
}
