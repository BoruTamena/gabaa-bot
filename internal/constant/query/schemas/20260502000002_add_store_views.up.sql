-- Create store_stats table for tracking views
CREATE TABLE IF NOT EXISTS store_stats (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT UNIQUE NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    views BIGINT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Initialize stats for existing stores
INSERT INTO store_stats (store_id)
SELECT id FROM stores
ON CONFLICT (store_id) DO NOTHING;
