package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"fmt"
	generated1 "translatorapi/graph/generated"
	"translatorapi/graph/model"
	"translatorapi/models"

	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

// CreateWord creates a new Polish word.
func (r *mutationResolver) CreateWord(ctx context.Context, polishWord string, englishWord *string, sentence *string) (*model.Word, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ensure the word doesn't already exist
	var ifAlreadyExists models.Word
	if err := tx.Where("polish_word = ?", polishWord).First(&ifAlreadyExists).Error; err == nil {
		tx.Rollback()
		return nil, fmt.Errorf("word already exists: %s", polishWord)
	}

	// Proceed to insert the new word
	word := models.Word{PolishWord: polishWord}
	if err := tx.Create(&word).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if englishWord != nil {
		translation := models.Translation{
			EnglishWord: *englishWord,
			WordID:      word.ID,
		}
		if err := tx.Create(&translation).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if sentence != nil {
			example := models.Example{
				Sentence:      *sentence,
				TranslationID: translation.ID,
			}
			if err := tx.Create(&example).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return ToGraphQLWord(&word), nil
}

// CreateTranslation creates a new translation for a word.
func (r *mutationResolver) CreateTranslation(ctx context.Context, polishWord string, englishWord string, sentence *string) (*model.Translation, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find the word by its PolishWord
	var word models.Word
	if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, fmt.Errorf("word not found: %v", err)
	}

	// Check if the translation already exists
	var existingTranslation models.Translation
	if err := tx.Where("english_word = ? AND word_id = ?", englishWord, word.ID).First(&existingTranslation).Error; err == nil {
		return nil, fmt.Errorf("translation '%s' already exists for this word", englishWord)
	}

	// Create the translation for the found word
	translation := models.Translation{
		WordID:      word.ID,
		EnglishWord: englishWord,
	}

	if err := tx.Create(&translation).Error; err != nil {
		return nil, err
	}

	if sentence != nil {
		example := models.Example{
			TranslationID: translation.ID,
			Sentence:      *sentence,
		}

		if err := tx.Create(&example).Error; err != nil {
			return nil, err
		}
	}

	// Everything succeeded, commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return ToGraphQLTranslation(&translation), nil
}

// CreateExample creates a new example sentence for a translation.
func (r *mutationResolver) CreateExample(ctx context.Context, polishWord string, englishWord string, sentence string) (*model.Example, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var word models.Word
	if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, err
	}

	var translation models.Translation
	if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
		return nil, fmt.Errorf("translation not found: %v", err)
	}

	// Check if the example already exists
	var existingExample models.Example
	if err := tx.Where("sentence = ? AND translation_id = ?", sentence, translation.ID).First(&existingExample).Error; err == nil {
		return nil, fmt.Errorf("example '%s' already exists for this translation", sentence)
	}

	// Create the example for the found translation
	example := models.Example{
		TranslationID: translation.ID,
		Sentence:      sentence,
	}

	if err := tx.Create(&example).Error; err != nil {
		return nil, err
	}

	// Everything succeeded, commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return ToGraphQLExample(&example), nil
}

func (r *mutationResolver) ReplaceTranslation(ctx context.Context, polishWord string, englishWord string, newTranslation string) (*model.Translation, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find the word by its PolishWord
	var word models.Word
	if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, fmt.Errorf("word not found: %v", err)
	}

	if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).Delete(models.Translation{}).Error; err != nil {
		// Create the translation for the found word
		translation := models.Translation{
			WordID:      word.ID,
			EnglishWord: newTranslation,
		}

		if err := tx.Create(&translation).Error; err != nil {
			return nil, err
		}
		return ToGraphQLTranslation(&translation), fmt.Errorf("nothing to replace, added new: %w", err)
	}

	// Create the translation for the found word
	translation := models.Translation{
		WordID:      word.ID,
		EnglishWord: newTranslation,
	}

	if err := tx.Create(&translation).Error; err != nil {
		return nil, fmt.Errorf("operation unsucesfull: %w", err)
	}

	// Everything succeeded, commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return ToGraphQLTranslation(&translation), nil

}

// DeleteWord is the resolver for the deleteWord field.
func (r *mutationResolver) DeleteWord(ctx context.Context, polishWord string) (bool, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return false, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("polish_word = ?", polishWord).Delete(&models.Word{}).Error; err != nil {
		return false, err
	}

	// Everything succeeded, commit transaction
	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("transaction commit failed: %w", err)
	}

	return true, nil
}

// DeleteTranslation is the resolver for the deleteTranslation field.
func (r *mutationResolver) DeleteTranslation(ctx context.Context, polishWord string, englishWord string) (bool, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return false, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var word models.Word
	if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return false, err
	}
	if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).Delete(models.Translation{}).Error; err != nil {
		return false, err
	}

	// Everything succeeded, commit transaction
	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("transaction commit failed: %w", err)
	}

	return true, nil
}

// DeleteExample is the resolver for the deleteExample field.
func (r *mutationResolver) DeleteExample(ctx context.Context, polishWord string, englishWord string, exampleSentence string) (bool, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return false, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var word models.Word
	if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return false, err
	}
	var translation models.Translation
	if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
		return false, err
	}
	if err := tx.Where("translation_id = ? AND sentence = ?", translation.ID, exampleSentence).Delete(models.Example{}).Error; err != nil {
		return false, err
	}
	// Everything succeeded, commit transaction
	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("transaction commit failed: %w", err)
	}

	return true, nil
}

// Words is the resolver for the words field.
func (r *queryResolver) Words(ctx context.Context) ([]*model.Word, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var words []*models.Word
	if err := tx.Preload("Translations.Examples").Find(&words).Error; err != nil {
		return nil, err
	}

	var gqlWords []*model.Word
	for _, word := range words {
		gqlWords = append(gqlWords, ToGraphQLWord(word))
	}

	// Everything succeeded, commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return gqlWords, nil
}

// Translations retrieves translations by PolishWord.
func (r *queryResolver) Translations(ctx context.Context, polishWord string) ([]*model.Translation, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var word models.Word
	if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, fmt.Errorf("word not found: %v", err)
	}

	var translations []*models.Translation
	if err := tx.Preload("Examples").Where("word_id = ?", word.ID).Find(&translations).Error; err != nil {
		return nil, fmt.Errorf("could not fetch translations: %v", err)
	}

	var gqlTranslations []*model.Translation
	for _, translation := range translations {
		// Convert translation to GraphQL format
		gqlTranslations = append(gqlTranslations, ToGraphQLTranslation(translation))
	}
	// Everything succeeded, commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return gqlTranslations, nil
}

// Examples retrieves examples by EnglishWord.
func (r *queryResolver) Examples(ctx context.Context, polishWord string, englishWord string) ([]*model.Example, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", tx.Error)
	}

	// Ensure rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var word models.Word
	if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		return nil, fmt.Errorf("word not found: %v", err)
	}

	var translation models.Translation
	if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
		return nil, err
	}

	var examples []*models.Example
	if err := tx.Where("translation_id = ?", translation.ID).Find(&examples).Error; err != nil {
		return nil, err
	}

	var gqlExamples []*model.Example
	for _, example := range examples {
		gqlExamples = append(gqlExamples, ToGraphQLExample(example))
	}
	// Everything succeeded, commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
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
