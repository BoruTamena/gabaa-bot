DROP INDEX IF EXISTS idx_withdrawals_deleted_at;
ALTER TABLE withdrawals DROP COLUMN IF EXISTS deleted_at;
