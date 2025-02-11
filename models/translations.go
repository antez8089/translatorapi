package models


// Translation represents an English translation of a Polish word
type Translation struct {
	ID         uint   `gorm:"primaryKey"`
	WordID      uint        `gorm:"not null;index"`
	EnglishWord string      `gorm:"not null"`
	Examples    []Example   `gorm:"foreignKey:TranslationID;constraint:OnDelete:CASCADE"`
}
