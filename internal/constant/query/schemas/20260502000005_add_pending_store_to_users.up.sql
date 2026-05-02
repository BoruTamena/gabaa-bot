-- Add pending_store_id to users table to track linking process
ALTER TABLE users ADD COLUMN IF NOT EXISTS pending_store_id BIGINT REFERENCES stores(id);
