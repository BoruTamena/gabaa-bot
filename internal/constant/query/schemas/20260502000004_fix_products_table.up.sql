-- Add seller_id and category to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS seller_id BIGINT REFERENCES users(id);
ALTER TABLE products ADD COLUMN IF NOT EXISTS category VARCHAR(100);

-- Update existing products to have a valid seller_id (from their store)
UPDATE products p
SET seller_id = s.seller_id
FROM stores s
WHERE p.store_id = s.id AND p.seller_id IS NULL;
