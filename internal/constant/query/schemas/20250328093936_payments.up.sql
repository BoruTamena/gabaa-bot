CREATE TABLE payments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID REFERENCES orders(id) ON DELETE CASCADE,
    buyer_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    amount      DECIMAL(10,2) NOT NULL,
    payment_method TEXT CHECK (payment_method IN ('credit_card', 'paypal', 'crypto', 'mobile_money')),
    status      TEXT CHECK (status IN ('pending', 'completed', 'failed')) DEFAULT 'pending',
    transaction_id TEXT UNIQUE NOT NULL,
    created_at  TIMESTAMP DEFAULT now()
     updated_at  TIMESTAMP DEFAULT NULL,
     deleted_at  TIMESTAMP DEFAULT NULL
);
