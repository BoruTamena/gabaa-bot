CREATE TABLE order_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID REFERENCES orders(id) ON DELETE CASCADE,
    product_id  UUID REFERENCES products(id) ON DELETE CASCADE,
    quantity    INT CHECK (quantity > 0) DEFAULT 1,
    price       DECIMAL(10,2) NOT NULL,
    created_at  TIMESTAMP DEFAULT now()
);
