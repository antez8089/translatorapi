package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"errors"
	"fmt"
	generated1 "translatorapi/graph/generated"
	"translatorapi/graph/model"
	"translatorapi/models"
	"gorm.io/gorm"
)

type Resolver struct{
	DB *gorm.DB
}

// CreateWord creates a new Polish word.
func (r *mutationResolver) CreateWord(ctx context.Context, polishWord string, englishWord *string, sentence *string) (*model.Word, error) {
	if (englishWord == nil) && (sentence != nil) {
		return nil, errors.New("sentence cannot be provided without an English word")
	}
	var ifAlreadyExists models.Word
	if err := r.DB.Where("polish_word = ?", polishWord).First(&ifAlreadyExists).Error; err == nil {
		return nil, errors.New("word already exists")
	}

	word := models.Word{
		PolishWord: polishWord,
	}

	if err := r.DB.Create(&word).Error; err != nil {
		return nil, err
	}
	if (englishWord != nil) && (sentence == nil) {
		translation := models.Translation{
			EnglishWord: *englishWord,
			WordID:      word.ID,
		}

		if err := r.DB.Create(&translation).Error; err != nil {
			return nil, err
		}

	} else if (englishWord != nil) && (sentence != nil) {
		translation := models.Translation{
			EnglishWord: *englishWord,
			WordID:      word.ID,
		}
		if err := r.DB.Create(&translation).Error; err != nil {
			return nil, err
		}

		example := models.Example{
			Sentence:      *sentence,
			TranslationID: translation.ID,
		}
		if err := r.DB.Create(&example).Error; err != nil {
			return nil, err
		}

	}

	return ToGraphQLWord(&word), nil
}

// CreateTranslation creates a new translation for a word.
func (r *mutationResolver) CreateTranslation(ctx context.Context, polishWord string, englishWord string, sentence *string) (*model.Translation, error) {
	// Find the word by its PolishWord
	var word models.Word
	if err := r.DB.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, fmt.Errorf("word not found: %v", err)
	}

	// Check if the translation already exists
	var existingTranslation models.Translation
	if err := r.DB.Where("english_word = ? AND word_id = ?", englishWord, word.ID).First(&existingTranslation).Error; err == nil {
		return nil, fmt.Errorf("translation '%s' already exists for this word", englishWord)
	}

	// Create the translation for the found word
	translation := models.Translation{
		WordID:      word.ID,
		EnglishWord: englishWord,
	}

	if err := r.DB.Create(&translation).Error; err != nil {
		return nil, err
	}

	if sentence != nil {
		example := models.Example{
			TranslationID: translation.ID,
			Sentence:      *sentence,
		}

		if err := r.DB.Create(&example).Error; err != nil {
			return nil, err
		}
	}

	return ToGraphQLTranslation(&translation), nil
}

// CreateExample creates a new example sentence for a translation.
func (r *mutationResolver) CreateExample(ctx context.Context, polishWord string, englishWord string, sentence string) (*model.Example, error) {
	var word models.Word
	if err := r.DB.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, err
	}

	var translation models.Translation
	if err := r.DB.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
		return nil, fmt.Errorf("translation not found: %v", err)
	}

	// Check if the example already exists
	var existingExample models.Example
	if err := r.DB.Where("sentence = ? AND translation_id = ?", sentence, translation.ID).First(&existingExample).Error; err == nil {
		return nil, fmt.Errorf("example '%s' already exists for this translation", sentence)
	}

	// Create the example for the found translation
	example := models.Example{
		TranslationID: translation.ID,
		Sentence:      sentence,
	}

	if err := r.DB.Create(&example).Error; err != nil {
		return nil, err
	}

	return ToGraphQLExample(&example), nil
}

// DeleteWord is the resolver for the deleteWord field.
func (r *mutationResolver) DeleteWord(ctx context.Context, polishWord string) (bool, error) {
	if err := r.DB.Where("polish_word = ?", polishWord).Delete(&models.Word{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

// DeleteTranslation is the resolver for the deleteTranslation field.
func (r *mutationResolver) DeleteTranslation(ctx context.Context, polishWord string, englishWord string) (bool, error) {
	var word models.Word
	if err := r.DB.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return false, err
	}
	if err := r.DB.Where("word_id = ? AND english_word = ?", word.ID, englishWord).Delete(models.Translation{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

// DeleteExample is the resolver for the deleteExample field.
func (r *mutationResolver) DeleteExample(ctx context.Context, polishWord string, englishWord string, exampleSentence string) (bool, error) {
	var word models.Word
	if err := r.DB.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return false, err
	}
	var translation models.Translation
	if err := r.DB.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
		return false, err
	}
	if err := r.DB.Where("translation_id = ? AND sentence = ?", translation.ID, exampleSentence).Delete(models.Example{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Words is the resolver for the words field.
func (r *queryResolver) Words(ctx context.Context) ([]*model.Word, error) {
	var words []*models.Word
	if err := r.DB.Preload("Translations.Examples").Find(&words).Error; err != nil {
		return nil, err
	}

	var gqlWords []*model.Word
	for _, word := range words {
		gqlWords = append(gqlWords, ToGraphQLWord(word))
	}

	return gqlWords, nil
}

// Translations retrieves translations by PolishWord.
func (r *queryResolver) Translations(ctx context.Context, polishWord string) ([]*model.Translation, error) {
	var word models.Word
	if err := r.DB.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, fmt.Errorf("word not found: %v", err)
	}

	var translations []*models.Translation
	if err := r.DB.Preload("Examples").Where("word_id = ?", word.ID).Find(&translations).Error; err != nil {
		return nil, fmt.Errorf("could not fetch translations: %v", err)
	}

	var gqlTranslations []*model.Translation
	for _, translation := range translations {
		// Convert translation to GraphQL format
		gqlTranslations = append(gqlTranslations, ToGraphQLTranslation(translation))
	}

	return gqlTranslations, nil
}

// Examples retrieves examples by EnglishWord.
func (r *queryResolver) Examples(ctx context.Context,polishWord string, englishWord string) ([]*model.Example, error) {
	var word models.Word
	if err := r.DB.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, fmt.Errorf("word not found: %v", err)
	}

	var translation models.Translation
	if err := r.DB.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
		return nil, err
	}

	var examples []*models.Example
	if err := r.DB.Where("translation_id = ?", translation.ID).Find(&examples).Error; err != nil {
		return nil, err
	}

	var gqlExamples []*model.Example
	for _, example := range examples {
		gqlExamples = append(gqlExamples, ToGraphQLExample(example))
	}

	return gqlExamples, nil
}

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated1.MutationResolver { return &mutationResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	type Resolver struct{
	DB *gorm.DB
}
*/
