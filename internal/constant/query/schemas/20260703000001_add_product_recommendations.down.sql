DROP TABLE IF EXISTS product_recommendations;
DROP TABLE IF EXISTS user_category_preferences;

ALTER TABLE users
    DROP COLUMN IF EXISTS recommendations_enabled,
    DROP COLUMN IF EXISTS bot_started;
