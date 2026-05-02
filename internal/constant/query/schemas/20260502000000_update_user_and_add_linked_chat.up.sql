-- Update users table
ALTER TABLE users ALTER COLUMN telegram_user_id DROP NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT;

-- Create linked_chats table
CREATE TABLE IF NOT EXISTS linked_chats (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT REFERENCES stores(id),
    chat_id BIGINT UNIQUE NOT NULL,
    chat_type VARCHAR(50),
    title VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
