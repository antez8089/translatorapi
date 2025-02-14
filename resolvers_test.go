package graph_test

import (
	"context"
	"testing"
	"translatorapi/graph"
	"github.com/stretchr/testify/assert"
	"translatorapi/graph/model"
	"translatorapi/mockdatabase"
)

func TestCreate(t *testing.T) {
	// Inicjalizujemy bazę danych w pamięci
	gormDB, err := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Nie udało się zainicjować mockowej bazy danych: %v", err)
	}

	// Tworzymy resolver z wykorzystaniem mockowanej bazy
	resolver := &graph.Resolver{DB: gormDB}
	mutationResolver := resolver.Mutation()
	quadResolver := resolver.Query()

	// Wywołujemy funkcję mutacji
	word, err := mutationResolver.CreateWord(context.TODO(), "a", nil, nil)
	if err != nil {
		t.Fatalf("CreateWord nie powiodło się: %v", err)
	}

	// Sprawdzamy zawartość tabeli "words"
	var words []model.Word
	if err := gormDB.Find(&words).Error; err != nil {
		t.Fatalf("Nie udało się pobrać danych z tabeli 'words': %v", err)
	}

	// Logujemy znalezione słowa
	for _, w := range words {
		t.Logf("Słowo: %+v", w)
	}

	// Definiujemy oczekiwany wynik
	expectedWord := model.Word{
		ID:           "1",
		PolishWord:   "a",
		Translations: []*model.Translation{}, // Pusta lista tłumaczeń
	}

	// Sprawdzamy, czy zwrócone słowo odpowiada oczekiwanemu
	assert.Equal(t, &expectedWord, word)

	// Sprawdzamy, czy zwrócone słowo zostało zapisane w bazie
	assert.Equal(t, 1, len(words))



	// Tworzymy tłumaczenie
	translation , err := mutationResolver.CreateTranslation(context.TODO(), "a", "b", nil) 
	if err != nil {
		t.Fatalf("CreateTranslation nie powiodło się: %v", err)
	}

	expectedTranslation := model.Translation{
		ID:              "1",
		WordID:          "1",
		EnglishWord:     "b",
		Examples: []*model.Example{},
	}

	assert.Equal(t, &expectedTranslation, translation)

	expectedWord.Translations = append(expectedWord.Translations, &expectedTranslation)


	// Sprawdzamy zawartość tabeli "words"
	if err := gormDB.Find(&words).Error; err != nil {
		t.Fatalf("Nie udało się pobrać danych z tabeli 'words': %v", err)
	}

	currenword, _ := quadResolver.Words(context.Background()) 

	assert.Equal(t, &expectedWord, currenword[0])

	


	// Tworzymy przykład
	example, err := mutationResolver.CreateExample(context.TODO(), "a", "b", "c")
	if err != nil {
		t.Fatalf("CreateExample nie powiodło się: %v", err)
	}

	expectedExample := model.Example{
		ID:              "1",
		TranslationID:   "1",
		Sentence:        "c",
	}
	
	assert.Equal(t, &expectedExample, example)
	
}
