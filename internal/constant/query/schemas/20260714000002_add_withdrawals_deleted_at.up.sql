ALTER TABLE withdrawals ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX IF NOT EXISTS idx_withdrawals_deleted_at ON withdrawals(deleted_at);
