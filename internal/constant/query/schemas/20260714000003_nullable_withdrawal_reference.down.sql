UPDATE withdrawals SET reference = 'WITHDRAW-' || store_id || '-' || id WHERE reference IS NULL;
ALTER TABLE withdrawals ALTER COLUMN reference SET NOT NULL;
