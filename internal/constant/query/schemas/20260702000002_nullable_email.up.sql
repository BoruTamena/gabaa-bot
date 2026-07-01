-- Make the email column nullable so Telegram-only users (with no email)
-- don't collide on the unique constraint. PostgreSQL allows multiple NULLs
-- on a unique index, but only one empty string.
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
UPDATE users SET email = NULL WHERE email = '';
