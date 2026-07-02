DROP INDEX IF EXISTS idx_user_category_preferences_categories;
DROP INDEX IF EXISTS idx_user_category_preferences_user_id_unique;

ALTER TABLE user_category_preferences
    ADD COLUMN IF NOT EXISTS category VARCHAR(255);

UPDATE user_category_preferences
SET category = categories->>0
WHERE category IS NULL
  AND jsonb_array_length(categories) > 0;

ALTER TABLE user_category_preferences DROP COLUMN IF EXISTS categories;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_category_preferences_unique
    ON user_category_preferences (user_id, lower(category))
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_category_preferences_user_id
    ON user_category_preferences(user_id);

CREATE INDEX IF NOT EXISTS idx_user_category_preferences_category_lower
    ON user_category_preferences(lower(category));
