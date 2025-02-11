-- init.sql

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
