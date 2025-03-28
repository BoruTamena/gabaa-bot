CREATE TABLE orders (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    seller_id   UUID REFERENCES users(id) ON DELETE CASCADE,
    status      TEXT CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')) DEFAULT 'pending',
    total_price DECIMAL(10,2) NOT NULL,
    created_at  TIMESTAMP DEFAULT now()
);
