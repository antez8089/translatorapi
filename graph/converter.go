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
		Translations: func() []*model.Translation {
			// Tworzenie pustej tablicy Translation
			translations := make([]*model.Translation, 0)
			// Iteracja po wszystkich Translation
			for _, t := range word.Translations {
				// Konwersja Translation na GraphQL Translation
				translations = append(translations, ToGraphQLTranslation(&t))
			}
			return translations
		}(),
	}
}

// Funkcja konwertująca Translation na GraphQL Translation
func ToGraphQLTranslation(t *models.Translation) *model.Translation {
	// Konwersja int na string, jeśli pole Translation.ID jest int
	return &model.Translation{
		ID:          strconv.Itoa(int(t.ID)),   // int na string
		WordID:      strconv.Itoa(int(t.WordID)), // Konwersja int na string
		EnglishWord: t.EnglishWord,
		Examples: func() []*model.Example {
			// Tworzenie pustej tablicy Example
			examples := make([]*model.Example, 0)
			// Iteracja po wszystkich Example
			for _, e := range t.Examples {
				// Konwersja Example na GraphQL Example
				examples = append(examples, ToGraphQLExample(&e))
			}
			return examples
		}(),
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
