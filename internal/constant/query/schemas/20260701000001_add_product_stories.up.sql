CREATE TABLE IF NOT EXISTS product_stories (
    id          BIGSERIAL    PRIMARY KEY,
    store_id    BIGINT       NOT NULL REFERENCES stores(id)   ON DELETE CASCADE,
    product_id  BIGINT       NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    caption     TEXT,
    media_urls  JSONB        NOT NULL DEFAULT '[]',
    media_type  VARCHAR(10)  NOT NULL DEFAULT 'image',
    starts_at   TIMESTAMPTZ  NOT NULL,
    ends_at     TIMESTAMPTZ  NOT NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    views       BIGINT       NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_product_stories_store_id    ON product_stories(store_id);
CREATE INDEX IF NOT EXISTS idx_product_stories_product_id  ON product_stories(product_id);
CREATE INDEX IF NOT EXISTS idx_product_stories_active_range ON product_stories(is_active, starts_at, ends_at)
    WHERE deleted_at IS NULL;
