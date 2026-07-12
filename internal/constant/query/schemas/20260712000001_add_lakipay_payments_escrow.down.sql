DROP TABLE IF EXISTS escrows;
DROP TABLE IF EXISTS payment_webhooks;

ALTER TABLE wallets ADD COLUMN IF NOT EXISTS balance NUMERIC DEFAULT 0;
UPDATE wallets SET balance = available_balance;
ALTER TABLE wallets DROP COLUMN IF EXISTS currency;
ALTER TABLE wallets DROP COLUMN IF EXISTS pending_balance;
ALTER TABLE wallets DROP COLUMN IF EXISTS available_balance;
ALTER TABLE wallets DROP COLUMN IF EXISTS locked_balance;
ALTER TABLE wallets DROP COLUMN IF EXISTS total_earned;
ALTER TABLE wallets DROP COLUMN IF EXISTS total_withdrawn;

DROP INDEX IF EXISTS idx_payments_reference;
DROP INDEX IF EXISTS idx_payments_transaction_id;
ALTER TABLE payments DROP COLUMN IF EXISTS reference;
ALTER TABLE payments DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE payments DROP COLUMN IF EXISTS amount;
ALTER TABLE payments DROP COLUMN IF EXISTS currency;
ALTER TABLE payments DROP COLUMN IF EXISTS phone_number;
ALTER TABLE payments DROP COLUMN IF EXISTS medium;
ALTER TABLE payments DROP COLUMN IF EXISTS gateway_status;
ALTER TABLE payments DROP COLUMN IF EXISTS gateway_response;
