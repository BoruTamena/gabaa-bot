-- Empty strings are not NULL and violate the partial unique index on transaction_id.
UPDATE payments SET transaction_id = NULL WHERE transaction_id = '';
UPDATE withdrawals SET transaction_id = NULL WHERE transaction_id = '';
