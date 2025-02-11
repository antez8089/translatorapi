package models


// Word represents a unique Polish word
type Word struct {
	ID         uint   `gorm:"primaryKey"`
	PolishWord   string         `gorm:"unique;not null"`
	Translations []Translation  `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"`

}
