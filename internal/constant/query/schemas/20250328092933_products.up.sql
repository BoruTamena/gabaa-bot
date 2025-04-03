CREATE TABLE products (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id   UUID REFERENCES users(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    description TEXT,
    price       DECIMAL(10,2) NOT NULL CHECK (price > 0),
    image_url   TEXT, 
    created_at  TIMESTAMP DEFAULT now()
    
);
