DROP INDEX IF EXISTS idx_withdrawals_store_id;
ALTER TABLE payment_webhooks DROP COLUMN IF EXISTS withdrawal_id;
DROP TABLE IF EXISTS withdrawals;
