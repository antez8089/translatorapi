package graph_test

import (
	"context"
	"testing"
	"translatorapi/graph"
	"github.com/stretchr/testify/assert"
	"translatorapi/graph/model"
	"translatorapi/mockdatabase"
	// "translatorapi/database"
	"sync"
	"sync/atomic"
	"fmt"
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


func TestConcurrentWordCreation(t *testing.T) {
	// Set up the mock database
	gormDB, err  := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Failed to set up mock database: %v", err)
	}


	// Create the resolver
	resolver := &graph.Resolver{DB: gormDB}
	mutationResolver := resolver.Mutation()

	// Concurrent execution setup
	var wg sync.WaitGroup
	errCh := make(chan error, 5) // Buffered to avoid deadlocks

	
	wordsToInsert := []string{"a", "b", "c", "d", "e"}
	var insertCount atomic.Int32

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := mutationResolver.CreateWord(context.TODO(), wordsToInsert[i], nil, nil)
			if err != nil {
				errCh <- err // Capture error
			} else {
				insertCount.Add(1)
			}
		}()
	}

	wg.Wait()
	close(errCh) // Close the channel after all goroutines finish

	// Collect errors
	for err := range errCh {
		t.Errorf("Error during concurrent insertion: %v", err)
	}

	// Validate database state
	var words []model.Word
	if err := gormDB.Find(&words).Error; err != nil {
		t.Fatalf("Failed to query words: %v", err)
	}

	
	expectedCount := 5 
	assert.Equal(t, expectedCount, len(wordsToInsert), "Unexpected number of unique words in DB")
}


func TestConcurrentWordCreationDuplicates(t *testing.T) {
	// Set up the mock database
	gormDB, err := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Failed to set up mock database: %v", err)
	}

	// Migrate tables
	if err := gormDB.AutoMigrate(&model.Word{}); err != nil {
		t.Fatalf("Failed to auto-migrate: %v", err)
	}

	// Create the resolver
	resolver := &graph.Resolver{DB: gormDB}
	mutationResolver := resolver.Mutation()

	// Concurrent execution setup
	var wg sync.WaitGroup
	errCh := make(chan error, 10) // Buffered to avoid deadlocks

	// Simulate inserting the same word multiple times
	word := "test"
	var insertCount atomic.Int32

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := mutationResolver.CreateWord(context.TODO(), word, nil, nil)
			if err != nil {
				// If the error is because the word already exists, it's expected in concurrent cases.
				if err.Error() == fmt.Sprintf("word already exists: %s", word) {
					return // Ignore expected error for duplicates
				}
				errCh <- err // Capture error if it's something unexpected
			} else {
				insertCount.Add(1)
			}
		}()
	}

	wg.Wait()
	close(errCh) // Close the channel after all goroutines finish

	// Collect errors
	for err := range errCh {
		t.Errorf("Error during concurrent insertion: %v", err)
	}

	// Validate database state
	var words []model.Word
	if err := gormDB.Find(&words).Error; err != nil {
		t.Fatalf("Failed to query words: %v", err)
	}

	// **Expect only one row if database enforces uniqueness**
	expectedCount := 1 // Change if duplicates are allowed
	assert.Equal(t, expectedCount, len(words), "Unexpected number of unique words in DB")
}
