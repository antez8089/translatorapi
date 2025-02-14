package graph_test

import (
	"context"
	"testing"
	"translatorapi/graph"
	"github.com/stretchr/testify/assert"
	"translatorapi/graph/model"
	"translatorapi/mockdatabase"
	"sync"
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

func TestCreateFull(t *testing.T) {

	// Inicjalizujemy bazę danych w pamięci
	gormDB, err := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Nie udało się zainicjować mockowej bazy danych: %v", err)
	}

	// Tworzymy resolver z wykorzystaniem mockowanej bazy
	resolver := &graph.Resolver{DB: gormDB}
	mutationResolver := resolver.Mutation()



	b := "b"
	c := "c"

	mutationResolver.CreateWord(context.TODO(), "a", &b, &c)

	// Sprawdzamy zawartość tabeli "words"
	var words []model.Word
	if err := gormDB.Find(&words).Error; err != nil {
		t.Fatalf("Nie udało się pobrać danych z tabeli 'words': %v", err)
	}
	assert.Equal(t, 1, len(words))

	var translations []model.Translation
	if err := gormDB.Find(&translations).Error; err != nil {
		t.Fatalf("Nie udało się pobrać danych z tabeli 'translations': %v", err)
	}
	assert.Equal(t, 1, len(translations))

	var examples []model.Example

	if err := gormDB.Find(&examples).Error; err != nil {
		t.Fatalf("Nie udało się pobrać danych z tabeli 'examples': %v", err)
	}	
	assert.Equal(t, 1, len(examples))

	_, err = mutationResolver.CreateWord(context.TODO(), "a", nil, nil)

	assert.Error(t, err)
	
}

func TestDelete(t *testing.T){
	

	// Inicjalizujemy bazę danych w pamięci
	gormDB, err := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Nie udało się zainicjować mockowej bazy danych: %v", err)
	}

	// Tworzymy resolver z wykorzystaniem mockowanej bazy
	resolver := &graph.Resolver{DB: gormDB}
	mutationResolver := resolver.Mutation()

	b := "b"
	c := "c"

	

	mutationResolver.CreateWord(context.TODO(), "a", &b, &c)
	mutationResolver.DeleteWord(context.TODO(), "a")

	// Sprawdzamy zawartość tabeli "words"
	var words []model.Word
	if err := gormDB.Find(&words).Error; err != nil {
		t.Fatalf("Nie udało się pobrać danych z tabeli 'words': %v", err)
	}
	assert.Equal(t, 0, len(words))

	var translations []model.Translation
	if err := gormDB.Find(&translations).Error; err != nil {
		t.Fatalf("Nie udało się pobrać danych z tabeli 'translations': %v", err)
	}	
	assert.Equal(t, 0, len(translations))


	var examples []model.Example
	if err := gormDB.Find(&examples).Error; err != nil {
		t.Fatalf("Nie udało się pobrać danych z tabeli 'examples': %v", err)
	}
	assert.Equal(t, 0, len(examples))


}


func TestConcurrentProcessing(t *testing.T) {
	// Set up the mock database
	gormDB, err := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Failed to set up mock database: %v", err)
	}

	// Migrate the schema to create the necessary tables
	if err := gormDB.AutoMigrate(&model.Word{}, &model.Translation{}, &model.Example{}); err != nil {
		t.Fatalf("Failed to auto-migrate: %v", err)
	}

	// Create the resolver
	resolver := &graph.Resolver{DB: gormDB}
	mutationResolver := resolver.Mutation()

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a slice of test data
	data := []string{"a", "b", "c", "d", "e"}

	// Start concurrent operations
	for _, word := range data {
		wg.Add(1) // Add to the WaitGroup for each goroutine
		go func(word string) {
			defer wg.Done() // Mark this goroutine as done when it finishes

			// Call your mutation or function that processes the data
			_, err := mutationResolver.CreateWord(context.TODO(), word, nil, nil)
			if err != nil {
				t.Errorf("Error creating word %s: %v", word, err)
			}
		}(word)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Now, check the results in the database (make assertions)
	var words []model.Word
	if err := gormDB.Find(&words).Error; err != nil {
		t.Fatalf("Failed to query words: %v", err)
	}

	// Verify that the correct number of words were created
	assert.Equal(t, len(data), len(words), "The number of words created should match the number of input words")
}
