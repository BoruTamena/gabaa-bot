ALTER TABLE stores
    ADD COLUMN IF NOT EXISTS verification_status VARCHAR(20) NOT NULL DEFAULT 'unverified';

CREATE TABLE IF NOT EXISTS store_kyc (
    id                           BIGSERIAL PRIMARY KEY,
    store_id                     BIGINT NOT NULL UNIQUE REFERENCES stores(id) ON DELETE CASCADE,
    tin_number                   VARCHAR(50) NOT NULL,
    business_registration_number VARCHAR(100) NOT NULL,
    tin_certificate_url          TEXT NOT NULL,
    business_license_url         TEXT NOT NULL,
    review_note                  TEXT,
    submitted_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at                  TIMESTAMPTZ,
    created_at                   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_store_kyc_store_id ON store_kyc(store_id);
