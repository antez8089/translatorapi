-- Ensure tables are created
CREATE TABLE IF NOT EXISTS words (
    id SERIAL PRIMARY KEY,
    polish_word VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS translations (
    id SERIAL PRIMARY KEY,
    word_id INT REFERENCES words(id) ON DELETE CASCADE,
    english_word VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS examples (
    id SERIAL PRIMARY KEY,
    translation_id INT REFERENCES translations(id) ON DELETE CASCADE,
    sentence TEXT NOT NULL
);

-- Add unique constraints safely
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'unique_polish_word'
    ) THEN
        ALTER TABLE words ADD CONSTRAINT unique_polish_word UNIQUE (polish_word);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'unique_english_word'
    ) THEN
        ALTER TABLE translations ADD CONSTRAINT unique_english_word UNIQUE (word_id, english_word);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'unique_sentence'
    ) THEN
        ALTER TABLE examples ADD CONSTRAINT unique_sentence UNIQUE (sentence, translation_id);
    END IF;
END $$;
