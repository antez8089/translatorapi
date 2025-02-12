package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"translatorapi/database"
	generated1 "translatorapi/graph/generated"
	"translatorapi/graph/model"
	"translatorapi/models"
)

// CreateWord is the resolver for the createWord field.
func (r *mutationResolver) CreateWord(ctx context.Context, polishWord string) (*model.Word, error) {
	word := models.Word{
		PolishWord: polishWord,
	}

	if err := database.DB.Create(&word).Error; err != nil {
		return nil, err
	}

	return ToGraphQLWord(&word), nil
}

// CreateTranslation is the resolver for the createTranslation field.
func (r *mutationResolver) CreateTranslation(ctx context.Context, wordID string, englishWord string) (*model.Translation, error) {
	id, err := strconv.ParseUint(wordID, 10, 64) // Konwersja z string na uint
	if err != nil {
		return nil, err
	}
	translation := models.Translation{
		WordID:      uint(id),
		EnglishWord: englishWord,
	}

	if err := database.DB.Create(&translation).Error; err != nil {
		return nil, err
	}

	return ToGraphQLTranslation(&translation), nil
}

// CreateExample is the resolver for the createExample field.
func (r *mutationResolver) CreateExample(ctx context.Context, translationID string, sentence string) (*model.Example, error) {
	id, err := strconv.ParseUint(translationID, 10, 64) // Konwersja z string na int
	if err != nil {
		return nil, err
	}
	example := models.Example{
		TranslationID: uint(id),
		Sentence:      sentence,
	}

	if err := database.DB.Create(&example).Error; err != nil {
		return nil, err
	}

	return ToGraphQLExample(&example), nil
}

// Create
func (r *mutationResolver) CreateWordWithTranslation(ctx context.Context, polishWord string, englishWord string) (*model.Word, error) {
	word := models.Word{
		PolishWord: polishWord,
	}

	if err := database.DB.Create(&word).Error; err != nil {
		return nil, err
	}

	if word.ID == 0 {
		return nil, errors.New("failed to retrieve Word ID after creation")
	}

	translation := models.Translation{
		EnglishWord: englishWord,
		WordID:      word.ID,
	}

	if err := database.DB.Create(&translation).Error; err != nil {
		return nil, err
	}
	// Preload translations so GraphQL can access them
	var completeWord models.Word
	if err := database.DB.Preload("Translations").First(&completeWord, word.ID).Error; err != nil {
		return nil, err
	}

	return ToGraphQLWord(&completeWord), nil
}

// Words is the resolver for the words field.
func (r *queryResolver) Words(ctx context.Context) ([]*model.Word, error) {
	var words []*models.Word
	// Preload translations here
	if err := database.DB.Preload("Translations").Find(&words).Error; err != nil {
		return nil, err
	}

	var gqlWords []*model.Word
	for _, word := range words {
		gqlWords = append(gqlWords, ToGraphQLWord(word))
	}

	return gqlWords, nil
}

// In your queryResolver
func (r *queryResolver) Translations(ctx context.Context, wordID string) ([]*model.Translation, error) {
	var translations []*model.Translation
	// Convert the wordID from string to uint (assuming it's stored as uint in the database)
	wordIDInt, err := strconv.ParseUint(wordID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid wordID: %v", err)
	}

	// Fetch translations related to the wordID
	if err := database.DB.Where("word_id = ?", wordIDInt).Find(&translations).Error; err != nil {
		return nil, fmt.Errorf("could not fetch translations: %v", err)
	}

	return translations, nil
}

// Examples is the resolver for the examples field.
func (r *queryResolver) Examples(ctx context.Context, translationID string) ([]*model.Example, error) {
	var examples []*models.Example
	if err := database.DB.Find(&examples).Error; err != nil {
		return nil, err
	}

	var gqlExamples []*model.Example
	for _, e := range examples {
		gqlExamples = append(gqlExamples, ToGraphQLExample(e))
	}

	return gqlExamples, nil
}

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated1.MutationResolver { return &mutationResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
