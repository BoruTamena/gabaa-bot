
CREATE TABLE users (
    tel_id      BIGINT PRIMARY KEY,
    username    TEXT UNIQUE,
    first_name  TEXT NOT NULL,
    last_name   TEXT NOT NULL,
    phone       TEXT UNIQUE DEFAULT NULL,
    role        TEXT CHECK (role IN ('buyer', 'seller')) DEFAULT 'buyer',
    created_at  TIMESTAMP DEFAULT now(),
    updated_at  TIMESTAMP DEFAULT NULL,
    deleted_at  TIMESTAMP DEFAULT NULL
);