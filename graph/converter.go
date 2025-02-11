package graph

import (
	"strconv"
	"translatorapi/graph/model"
	"translatorapi/models"
)

// Funkcja konwertująca Word na GraphQL Word
func ToGraphQLWord(word *models.Word) *model.Word {
	// Konwersja int na string, jeśli pole Word.ID jest int
	return &model.Word{
		ID:         strconv.Itoa(int(word.ID)),  // int na string
		PolishWord: word.PolishWord,
	}
}

// Funkcja konwertująca Translation na GraphQL Translation
func ToGraphQLTranslation(t *models.Translation) *model.Translation {
	// Konwersja int na string, jeśli pole Translation.ID jest int
	return &model.Translation{
		ID:          strconv.Itoa(int(t.ID)),   // int na string
		WordID:      strconv.Itoa(int(t.WordID)), // Konwersja int na string
		EnglishWord: t.EnglishWord,
	}
}

// Funkcja konwertująca Example na GraphQL Example
func ToGraphQLExample(e *models.Example) *model.Example {
	// Konwersja int na string, jeśli pole Example.ID jest int
	return &model.Example{
		ID:            strconv.Itoa(int(e.ID)),     // int na string
		TranslationID: strconv.Itoa(int(e.TranslationID)), // int na string
		Sentence:      e.Sentence,
	}
}
