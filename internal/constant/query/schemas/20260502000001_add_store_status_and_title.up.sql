-- Add status and telegram_chat_title to stores table
ALTER TABLE stores ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'pending';
ALTER TABLE stores ADD COLUMN IF NOT EXISTS telegram_chat_title VARCHAR(255);

-- Remove unique constraint and not-null from telegram_chat_id to allow pending stores
ALTER TABLE stores ALTER COLUMN telegram_chat_id DROP NOT NULL;
DROP INDEX IF EXISTS idx_stores_telegram_chat_id;
CREATE INDEX IF NOT EXISTS idx_stores_telegram_chat_id ON stores(telegram_chat_id);
