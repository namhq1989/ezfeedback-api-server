-- Create words table
CREATE TABLE words (
                       id VARCHAR(50) PRIMARY KEY,
                       word VARCHAR(50) NOT NULL,
                       level varchar(20) NOT NULL,
                       parts_of_speech TEXT[] NOT NULL,
                       ipa VARCHAR(50) NOT NULL,
                       definitions JSONB NOT NULL, -- JSON Format: {"pos": "string", "definition": "string"}
                       noun_form JSONB, -- JSON Format: {"base": "string", "plural": "string"}
                       verb_form JSONB, -- JSON Format: {"base": "string", "past": "string", "pastParticiple": "string", "gerund": "string", "presentThirdPerson": "string"}
                       created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                       last_fetched_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Create an index on the word column for faster lookups
CREATE INDEX idx_words_word ON words(word);

-- Create an index on the level and last_fetched_at columns for faster lookups
CREATE INDEX idx_words_level_pos_fetched ON words
    USING btree (level, last_fetched_at)
    INCLUDE (parts_of_speech);


-- Create word_examples table
CREATE TABLE word_examples (
                               id VARCHAR(50) PRIMARY KEY,
                               word_id VARCHAR(50) NOT NULL REFERENCES words(id) ON DELETE CASCADE,
                               example TEXT NOT NULL,
                               main_word VARCHAR(50) NOT NULL,
                               level varchar(20) NOT NULL,
                               created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Create index for faster word_id lookups
CREATE INDEX idx_word_examples_word_id ON word_examples(word_id);

-- Create word_news table
CREATE TABLE word_news (
                           id VARCHAR(50) PRIMARY KEY,
                           word_id VARCHAR(50) NOT NULL REFERENCES words(id) ON DELETE CASCADE,
                           categories TEXT[] NOT NULL,
                           source_url TEXT NOT NULL,
                           title TEXT NOT NULL,
                           summary TEXT NOT NULL,
                           image_url TEXT NOT NULL,
                           created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           published_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster word_id lookups
CREATE INDEX idx_word_news_word_id ON word_news(word_id);
-- Create index for faster source_url lookups
CREATE INDEX idx_word_news_source_url ON word_news(source_url);

-- Create users table
CREATE TABLE users (
                       id VARCHAR(50) PRIMARY KEY,
                       belongs_to VARCHAR(50) REFERENCES users(id) ON DELETE SET NULL,
                       name VARCHAR(255) NOT NULL,
                       subscription JSONB NOT NULL DEFAULT '{}', -- JSON Format: {"plan": "free|pro", "expiry": "ISO8601 date string", "customerId": "string"}
                       auth_providers JSONB NOT NULL DEFAULT '{}', -- JSON Format: [{"provider": "extension|app|google", "id": "string", "name": "string", "email": "string"}]
                       created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                       updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Create index for faster belongs_to lookups
CREATE INDEX idx_users_belongs_to ON users(belongs_to);

-- Create index for subscription plan queries using JSONB operators
CREATE INDEX idx_users_subscription_plan ON users((subscription->>'plan'));