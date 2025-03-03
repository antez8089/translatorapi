package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"translatorapi/graph"
	generated "translatorapi/graph/generated"
	"translatorapi/graph/model"
	"translatorapi/mockdatabase"
	"translatorapi/models"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/ast"
	// "translatorapi/database"
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
	translation, err := mutationResolver.CreateTranslation(context.TODO(), "a", "b", nil)
	if err != nil {
		t.Fatalf("CreateTranslation nie powiodło się: %v", err)
	}

	expectedTranslation := model.Translation{
		ID:          "1",
		WordID:      "1",
		EnglishWord: "b",
		Examples:    []*model.Example{},
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
		ID:            "1",
		TranslationID: "1",
		Sentence:      "c",
	}

	assert.Equal(t, &expectedExample, example)

	gormDB.Exec("TRUNCATE words, translations, examples RESTART IDENTITY CASCADE;")

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

	gormDB.Exec("TRUNCATE words, translations, examples RESTART IDENTITY CASCADE;")

}

func TestDelete(t *testing.T) {

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

	gormDB.Exec("TRUNCATE words, translations, examples RESTART IDENTITY CASCADE;")

}

func TestCreateWordMutation(t *testing.T) {
	// Initialize mock database
	db, err := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Failed to initialize mock database: %v", err)
	}

	// Create GraphQL server with test database
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	// Prepare GraphQL mutation request
	query := `{ "query": "mutation { createWord(polishWord: \"a\") { polishWord id } }" }`
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte(query)))
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response body
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify response contains expected data
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "Response data missing")

	createWord, ok := data["createWord"].(map[string]interface{})
	assert.True(t, ok, "createWord response missing")
	assert.Equal(t, "a", createWord["polishWord"], "Unexpected polishWord value")

	// Verify object was actually saved in the database
	var word models.Word
	if err := db.First(&word).Error; err != nil {
		t.Fatalf("Word not found in database: %v", err)
	}
	assert.Equal(t, "a", word.PolishWord, "Unexpected value in database")

	db.Exec("TRUNCATE words, translations, examples RESTART IDENTITY CASCADE;")
}

func TestConcurrentCreateWordMutations(t *testing.T) {
	// Initialize mock database
	db, err := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Failed to initialize mock database: %v", err)
	}

	// Create GraphQL server with test database
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	// Prepare GraphQL mutation requests
	words := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew", "kiwi", "lemon"}
	var wg sync.WaitGroup

	for _, word := range words {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			query := `{ "query": "mutation { createWord(polishWord: \"` + w + `\") { polishWord id } }" }`
			resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte(query)))
			if err != nil {
				t.Errorf("Failed to execute request: %v", err)
				return
			}
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}(word)
	}

	wg.Wait()

	// Verify that 10 words were actually saved in the database
	var count int64
	if err := db.Model(&models.Word{}).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count words in database: %v", err)
	}
	assert.Equal(t, int64(10), count, "Unexpected number of words in database")

	db.Exec("TRUNCATE words, translations, examples RESTART IDENTITY CASCADE;")
}

func TestConcurrentLocking(t *testing.T) {
	// Initialize mock database
	db, err := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Failed to initialize mock database: %v", err)
	}

	// Create GraphQL server with test database
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	// Prepare GraphQL mutation requests
	words := []string{"a", "a", "a", "a", "b", "b", "a", "b", "b", "a", "a", "a", "a", "a", "b", "b", "a", "b", "b", "a"}
	var wg sync.WaitGroup

	for _, word := range words {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			query := `{ "query": "mutation { createWord(polishWord: \"` + w + `\") { polishWord id } }" }`
			resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte(query)))
			if err != nil {
				t.Errorf("Failed to execute request: %v", err)
				return
			}
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}(word)
	}

	wg.Wait()

	// Verify that 10 words were actually saved in the database
	var count int64
	if err := db.Model(&models.Word{}).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count words in database: %v", err)
	}
	assert.Equal(t, int64(2), count, "Unexpected number of words in database")

	db.Exec("TRUNCATE words, translations, examples RESTART IDENTITY CASCADE;")
}

func TestConcurrentTrnaslations(t *testing.T) {
	// Initialize mock database
	db, err := mockdatabase.MockDB(t)
	if err != nil {
		t.Fatalf("Failed to initialize mock database: %v", err)
	}

	// Create GraphQL server with test database
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{DB: db},
	}))

	ts := httptest.NewServer(srv)
	defer ts.Close()

	// Prepare GraphQL mutation request for creating the word
	query := `{ "query": "mutation { createWord(polishWord: \"a\") { polishWord id } }" }`
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte(query)))
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status for creating the word
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response body
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify response contains expected data for creating the word
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "Response data missing")

	createWord, ok := data["createWord"].(map[string]interface{})
	assert.True(t, ok, "createWord response missing")
	assert.Equal(t, "a", createWord["polishWord"], "Unexpected polishWord value")

	// Verify object was actually saved in the database
	var word models.Word
	if err := db.First(&word).Error; err != nil {
		t.Fatalf("Word not found in database: %v", err)
	}
	assert.Equal(t, "a", word.PolishWord, "Unexpected value in database")

	// Prepare GraphQL Translations

	// Prepare GraphQL mutation requests
	words := []string{"b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	var wg sync.WaitGroup

	for _, word := range words {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			// Properly format the query using fmt.Sprintf
			query := fmt.Sprintf(`{ "query": "mutation { createTranslation(polishWord: \"a\", englishWord:\"%s\") { id wordID englishWord } }" }`, w)
			resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte(query)))
			if err != nil {
				t.Errorf("Failed to execute request: %v", err)
				return
			}
			defer resp.Body.Close()

			// Check HTTP status for each translation mutation
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}(word)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Verify that 10 translations were added to the database
	var count int64
	if err := db.Model(&models.Translation{}).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count words in database: %v", err)
	}
	assert.Equal(t, int64(10), count, "Unexpected number of translations in database")

	db.Exec("TRUNCATE words, translations, examples RESTART IDENTITY CASCADE;")
}
