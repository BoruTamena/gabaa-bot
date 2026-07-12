-- Extend payments for LakiPay gateway
ALTER TABLE payments ADD COLUMN IF NOT EXISTS reference VARCHAR(100);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS transaction_id VARCHAR(100);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS amount NUMERIC;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS currency VARCHAR(10) DEFAULT 'ETB';
ALTER TABLE payments ADD COLUMN IF NOT EXISTS phone_number VARCHAR(20);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS medium VARCHAR(50);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS gateway_status VARCHAR(50);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS gateway_response JSONB;

CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_reference ON payments(reference) WHERE reference IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_transaction_id ON payments(transaction_id) WHERE transaction_id IS NOT NULL;

-- Redesign wallets: balance -> available_balance + escrow fields
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS currency VARCHAR(10) DEFAULT 'ETB';
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS pending_balance NUMERIC DEFAULT 0;
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS available_balance NUMERIC DEFAULT 0;
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS locked_balance NUMERIC DEFAULT 0;
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS total_earned NUMERIC DEFAULT 0;
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS total_withdrawn NUMERIC DEFAULT 0;

UPDATE wallets SET available_balance = balance WHERE available_balance = 0 AND balance > 0;
ALTER TABLE wallets DROP COLUMN IF EXISTS balance;

-- Payment webhook audit log
CREATE TABLE IF NOT EXISTS payment_webhooks (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT REFERENCES payments(id),
    transaction_id VARCHAR(100),
    event VARCHAR(50),
    status VARCHAR(50),
    payload JSONB NOT NULL,
    signature TEXT,
    verified BOOLEAN DEFAULT FALSE,
    processed BOOLEAN DEFAULT FALSE,
    received_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payment_webhooks_transaction_id ON payment_webhooks(transaction_id);

-- Escrow holds per paid order
CREATE TABLE IF NOT EXISTS escrows (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE REFERENCES orders(id),
    store_id BIGINT NOT NULL REFERENCES stores(id),
    amount NUMERIC NOT NULL,
    currency VARCHAR(10) DEFAULT 'ETB',
    status VARCHAR(50) NOT NULL DEFAULT 'held',
    release_at TIMESTAMP WITH TIME ZONE,
    released_at TIMESTAMP WITH TIME ZONE,
    refund_amount NUMERIC DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_escrows_store_id ON escrows(store_id);
