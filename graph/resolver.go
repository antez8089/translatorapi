package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"fmt"
	generated1 "translatorapi/graph/generated"
	"translatorapi/graph/model"
	"translatorapi/models"

	// "github.com/lib/pq"
	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

// CreateWord creates a new Polish word.
func (r *mutationResolver) CreateWord(ctx context.Context, polishWord string, englishWord *string, sentence *string) (*model.Word, error) {
	var word models.Word

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		word = models.Word{PolishWord: polishWord}

		result := tx.Where(&word).FirstOrCreate(&word)

		// If there was an error
		if result.Error != nil {
			return fmt.Errorf("failed to create word: %v", result.Error)
		}

		// No error, check if the word was created or already existed
		if result.RowsAffected == 0 {
			// No rows were affected, meaning the word already existed
			return fmt.Errorf("word already exists: %s", polishWord)
		}
		// Optionally add translation and example
		if englishWord != nil {
			translation := models.Translation{
				EnglishWord: *englishWord,
				WordID:      word.ID,
			}
			if err := tx.Create(&translation).Error; err != nil {
				return err
			}

			if sentence != nil {
				example := models.Example{
					Sentence:      *sentence,
					TranslationID: translation.ID,
				}
				if err := tx.Create(&example).Error; err != nil {
					return err
				}
			}
		}

		return nil // triggers commit
	})

	if err != nil {
		return nil, err // triggers rollback
	}

	return ToGraphQLWord(&word), nil
}

// CreateTranslation creates a new translation for a word.
func (r *mutationResolver) CreateTranslation(ctx context.Context, polishWord string, englishWord string, sentence *string) (*model.Translation, error) {
	var translation models.Translation
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		// Find the word by its PolishWord
		var word models.Word
		if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("polish word not found: %s", polishWord)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		// // Check if the translation already exists
		// var existingTranslation models.Translation
		// if err := tx.Where("english_word = ? AND word_id = ?", englishWord, word.ID).First(&existingTranslation).Error; err != nil {
		// 	fmt.Println(err)
		// 	if err == gorm.ErrRecordNotFound {
		// 		return  fmt.Errorf("translation '%s' already exists for this word", englishWord)
		// 	}
		// 	return  fmt.Errorf("an error occurred: %v", err)
		// }

		// Create the translation for the found word
		translation = models.Translation{
			WordID:      word.ID,
			EnglishWord: englishWord,
		}

		result := tx.Where(&translation).FirstOrCreate(&translation)

		// If there was an error
		if result.Error != nil {
			return fmt.Errorf("failed to create word: %v", result.Error)
		}

		// No error, check if the word was created or already existed
		if result.RowsAffected == 0 {
			// No rows were affected, meaning the word already existed
			return fmt.Errorf("translation already exists: %s", englishWord)
		}

		if sentence != nil {
			example := models.Example{
				TranslationID: translation.ID,
				Sentence:      *sentence,
			}

			if err := tx.Create(&example).Error; err != nil {
				return err
			}
		}

		return nil // triggers commit
	})

	if err != nil {
		return nil, err // triggers rollback
	}

	return ToGraphQLTranslation(&translation), nil
}

// CreateExample creates a new example sentence for a translation.
func (r *mutationResolver) CreateExample(ctx context.Context, polishWord string, englishWord string, sentence string) (*model.Example, error) {
	var example models.Example
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		if tx.Error != nil {
			return fmt.Errorf("could not begin transaction: %w", tx.Error)
		}

		// Ensure rollback in case of panic
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		var word models.Word
		if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("polish word not found: %s", polishWord)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		var translation models.Translation
		if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("translation maching polish word not found: %s", englishWord)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		// // Check if the example already exists
		// var existingExample models.Example
		// if err := tx.Where("sentence = ? AND translation_id = ?", sentence, translation.ID).First(&existingExample).Error; err == nil {
		// 	if err == gorm.ErrRecordNotFound {
		// 		return fmt.Errorf("example '%s' already exists for this translation", sentence)
		// 	}
		// 	return fmt.Errorf("an error occurred: %v", err)
		// }

		// Create the example for the found translation
		example = models.Example{
			TranslationID: translation.ID,
			Sentence:      sentence,
		}

		result := tx.Where(&example).FirstOrCreate(&example)

		// If there was an error
		if result.Error != nil {
			return fmt.Errorf("failed to create example: %v", result.Error)
		}

		// No error, check if the word was created or already existed
		if result.RowsAffected == 0 {
			// No rows were affected, meaning the word already existed
			return fmt.Errorf("example already exists: %s", sentence)
		}

		return nil
	})

	if err != nil {
		return nil, err // triggers rollback
	}

	return ToGraphQLExample(&example), nil
}

func (r *mutationResolver) ReplaceTranslation(ctx context.Context, polishWord string, englishWord string, newTranslation string) (*model.Translation, error) {
	var translation models.Translation
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		// Find the word by its PolishWord
		var word models.Word
		if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("word not found: %v", err)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).Delete(models.Translation{}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create the translation for the found word
				translation := models.Translation{
					WordID:      word.ID,
					EnglishWord: newTranslation,
				}

				result := tx.Where(&translation).FirstOrCreate(&translation)

				// If there was an error
				if result.Error != nil {
					return fmt.Errorf("failed to create translation: %v", result.Error)
				}

				// No error, check if the word was created or already existed
				if result.RowsAffected == 0 {
					// No rows were affected, meaning the word already existed
					return fmt.Errorf("translation already exists: %s", newTranslation)
				}
				return nil

			}
			return fmt.Errorf("operation unsucesfull: %w", err)
		}

		// Create the translation for the found word
		translation = models.Translation{
			WordID:      word.ID,
			EnglishWord: newTranslation,
		}

		result := tx.Where(&translation).FirstOrCreate(&translation)

		// If there was an error
		if result.Error != nil {
			return fmt.Errorf("failed to create translation: %v", result.Error)
		}

		// No error, check if the word was created or already existed
		if result.RowsAffected == 0 {
			// No rows were affected, meaning the word already existed
			return fmt.Errorf("translation already exists: %s", newTranslation)
		}
		return nil
	})

	if err != nil {
		return nil, err // triggers rollback
	}

	return ToGraphQLTranslation(&translation), nil

}

// DeleteWord is the resolver for the deleteWord field.
func (r *mutationResolver) DeleteWord(ctx context.Context, polishWord string) (bool, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		var word models.Word
		if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("word not found: %s", polishWord)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		if err := tx.Where("polish_word = ?", polishWord).Delete(&models.Word{}).Error; err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return false, err // triggers rollback
	}

	return true, nil
}

// DeleteTranslation is the resolver for the deleteTranslation field.
func (r *mutationResolver) DeleteTranslation(ctx context.Context, polishWord string, englishWord string) (bool, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		var word models.Word
		if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("word not found: %s", polishWord)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		var translation models.Translation
		if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("translation not found: %s", englishWord)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).Delete(models.Translation{}).Error; err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return false, err // triggers rollback
	}

	return true, nil
}

// DeleteExample is the resolver for the deleteExample field.
func (r *mutationResolver) DeleteExample(ctx context.Context, polishWord string, englishWord string, exampleSentence string) (bool, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		var word models.Word
		if err := tx.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("word not found: %s", polishWord)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		var translation models.Translation
		if err := tx.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("translation not found: %s", englishWord)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		var example models.Example
		if err := tx.Where("translation_id = ? AND sentence = ?", translation.ID, exampleSentence).First(&example).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("example not found: %s", exampleSentence)
			}
			return fmt.Errorf("an error occurred: %v", err)
		}

		if err := tx.Where("translation_id = ? AND sentence = ?", translation.ID, exampleSentence).Delete(models.Example{}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return false, err // triggers rollback
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
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("word not found: %s", polishWord)
		}
		return nil, fmt.Errorf("an error occurred: %v", err)
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
func (r *queryResolver) Examples(ctx context.Context, polishWord string, englishWord string) ([]*model.Example, error) {

	var word models.Word
	if err := r.DB.Where("polish_word = ?", polishWord).First(&word).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("word not found: %v", err)
		}
		return nil, fmt.Errorf("an error occurred: %v", err)
	}

	var translation models.Translation
	if err := r.DB.Where("word_id = ? AND english_word = ?", word.ID, englishWord).First(&translation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("an error occurred: %v", err)
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
