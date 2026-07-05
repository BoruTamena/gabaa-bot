DROP INDEX IF EXISTS idx_store_kyc_store_id;
DROP TABLE IF EXISTS store_kyc;
ALTER TABLE stores DROP COLUMN IF EXISTS verification_status;
