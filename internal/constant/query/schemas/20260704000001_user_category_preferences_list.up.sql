-- Refactor user_category_preferences: one row per user with a categories list (JSONB).

ALTER TABLE user_category_preferences
    ADD COLUMN IF NOT EXISTS categories JSONB NOT NULL DEFAULT '[]'::jsonb;

-- Migrate legacy per-category rows into the list column.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'user_category_preferences'
          AND column_name = 'category'
    ) THEN
        WITH aggregated AS (
            SELECT
                user_id,
                COALESCE(
                    jsonb_agg(DISTINCT category ORDER BY category)
                        FILTER (WHERE category IS NOT NULL AND category <> ''),
                    '[]'::jsonb
                ) AS cats
            FROM user_category_preferences
            WHERE deleted_at IS NULL
            GROUP BY user_id
        )
        UPDATE user_category_preferences ucp
        SET categories = aggregated.cats
        FROM aggregated
        WHERE ucp.user_id = aggregated.user_id;

        -- Keep a single row per user.
        DELETE FROM user_category_preferences a
        USING user_category_preferences b
        WHERE a.user_id = b.user_id
          AND a.id > b.id;

        DROP INDEX IF EXISTS idx_user_category_preferences_unique;
        DROP INDEX IF EXISTS idx_user_category_preferences_category_lower;

        ALTER TABLE user_category_preferences DROP COLUMN category;
    END IF;
END $$;

DROP INDEX IF EXISTS idx_user_category_preferences_user_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_category_preferences_user_id_unique
    ON user_category_preferences (user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_category_preferences_categories
    ON user_category_preferences USING GIN (categories);
