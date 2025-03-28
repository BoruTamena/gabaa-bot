CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_id BIGINT UNIQUE NOT NULL,
    username    TEXT UNIQUE,
    full_name   TEXT NOT NULL,
    phone       TEXT UNIQUE,
    role        TEXT CHECK (role IN ('buyer', 'seller')) DEFAULT 'buyer',
    created_at  TIMESTAMP DEFAULT now()
);
