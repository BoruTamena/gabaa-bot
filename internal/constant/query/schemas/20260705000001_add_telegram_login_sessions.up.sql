CREATE TABLE IF NOT EXISTS telegram_login_sessions (
    id               VARCHAR(64) PRIMARY KEY,
    status           VARCHAR(20) NOT NULL DEFAULT 'pending',
    telegram_user_id BIGINT,
    username         VARCHAR(255),
    expires_at       TIMESTAMPTZ NOT NULL,
    completed_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_telegram_login_sessions_expires_at ON telegram_login_sessions (expires_at);
