-- Reference is assigned after insert (needs withdrawal id), same as payments.
ALTER TABLE withdrawals ALTER COLUMN reference DROP NOT NULL;
