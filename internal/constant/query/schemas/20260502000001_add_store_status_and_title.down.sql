ALTER TABLE stores DROP COLUMN IF EXISTS status;
ALTER TABLE stores DROP COLUMN IF EXISTS telegram_chat_title;

-- Re-add constraints (Note: This might fail if there are now duplicate chat IDs or NULLs)
-- ALTER TABLE stores ALTER COLUMN telegram_chat_id SET NOT NULL;
-- CREATE UNIQUE INDEX idx_stores_telegram_chat_id ON stores(telegram_chat_id);
