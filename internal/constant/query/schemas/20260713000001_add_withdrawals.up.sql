CREATE TABLE IF NOT EXISTS withdrawals (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES stores(id),
    amount NUMERIC NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'ETB',
    phone_number VARCHAR(20) NOT NULL,
    medium VARCHAR(50) NOT NULL,
    reference VARCHAR(100) NOT NULL UNIQUE,
    transaction_id VARCHAR(100) UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'initiated',
    gateway_status VARCHAR(50),
    gateway_response JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_withdrawals_store_id ON withdrawals(store_id);

ALTER TABLE payment_webhooks ADD COLUMN IF NOT EXISTS withdrawal_id BIGINT REFERENCES withdrawals(id);
