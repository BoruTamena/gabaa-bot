CREATE TABLE IF NOT EXISTS product_inquiries (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    product_id BIGINT NOT NULL REFERENCES products(id),
    store_id BIGINT NOT NULL REFERENCES stores(id),
    customer_id BIGINT NOT NULL REFERENCES users(id),
    status TEXT NOT NULL DEFAULT 'pending',
    quantity INTEGER NOT NULL,
    note TEXT,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    reviewed_by BIGINT NULL REFERENCES users(id),
    reviewed_at TIMESTAMPTZ NULL,
    review_note TEXT
);

CREATE INDEX IF NOT EXISTS idx_product_inquiries_product_id ON product_inquiries(product_id);
CREATE INDEX IF NOT EXISTS idx_product_inquiries_store_id ON product_inquiries(store_id);
CREATE INDEX IF NOT EXISTS idx_product_inquiries_customer_id ON product_inquiries(customer_id);
CREATE INDEX IF NOT EXISTS idx_product_inquiries_status ON product_inquiries(status);
CREATE INDEX IF NOT EXISTS idx_product_inquiries_deleted_at ON product_inquiries(deleted_at);
