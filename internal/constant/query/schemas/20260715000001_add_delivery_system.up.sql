CREATE TABLE IF NOT EXISTS delivery_agents (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    user_id BIGINT REFERENCES users(id),
    telegram_user_id BIGINT,
    loyalty_score INT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_invite',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_delivery_agents_username ON delivery_agents(LOWER(username)) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_delivery_agents_telegram_user_id ON delivery_agents(telegram_user_id);
CREATE INDEX IF NOT EXISTS idx_delivery_agents_user_id ON delivery_agents(user_id);

CREATE TABLE IF NOT EXISTS store_delivery_links (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES stores(id),
    delivery_agent_id BIGINT NOT NULL REFERENCES delivery_agents(id),
    share_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    connected_by_user_id BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(store_id, delivery_agent_id)
);

CREATE INDEX IF NOT EXISTS idx_store_delivery_links_store_id ON store_delivery_links(store_id);
CREATE INDEX IF NOT EXISTS idx_store_delivery_links_agent_id ON store_delivery_links(delivery_agent_id);

CREATE TABLE IF NOT EXISTS delivery_routes (
    id BIGSERIAL PRIMARY KEY,
    store_delivery_link_id BIGINT NOT NULL REFERENCES store_delivery_links(id) ON DELETE CASCADE,
    label VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_delivery_routes_link_id ON delivery_routes(store_delivery_link_id);

CREATE TABLE IF NOT EXISTS delivery_route_locations (
    id BIGSERIAL PRIMARY KEY,
    delivery_route_id BIGINT NOT NULL REFERENCES delivery_routes(id) ON DELETE CASCADE,
    location_type VARCHAR(20) NOT NULL,
    label VARCHAR(255),
    country VARCHAR(100) DEFAULT 'Ethiopia',
    region VARCHAR(100),
    city VARCHAR(100),
    street VARCHAR(255),
    landmark VARCHAR(255),
    notes TEXT,
    use_store_location BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_delivery_route_locations_route_id ON delivery_route_locations(delivery_route_id);
CREATE INDEX IF NOT EXISTS idx_delivery_route_locations_type ON delivery_route_locations(location_type);

CREATE TABLE IF NOT EXISTS delivery_agent_shares (
    id BIGSERIAL PRIMARY KEY,
    owner_store_id BIGINT NOT NULL REFERENCES stores(id),
    delivery_agent_id BIGINT NOT NULL REFERENCES delivery_agents(id),
    adopted_store_id BIGINT NOT NULL REFERENCES stores(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_delivery_agent_shares_adopted ON delivery_agent_shares(adopted_store_id);

CREATE TABLE IF NOT EXISTS delivery_area_presets (
    id BIGSERIAL PRIMARY KEY,
    region VARCHAR(100),
    city VARCHAR(100),
    street VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

INSERT INTO delivery_area_presets (region, city, street) VALUES
    ('Bole', 'Addis Ababa', NULL),
    ('Bole', 'Addis Ababa', 'Bole Road'),
    ('Bole', 'Addis Ababa', 'Atlas Avenue'),
    ('CMC', 'Addis Ababa', NULL),
    ('CMC', 'Addis Ababa', 'CMC Road'),
    ('Piassa', 'Addis Ababa', NULL),
    ('Kazanchis', 'Addis Ababa', NULL),
    ('Megenagna', 'Addis Ababa', NULL),
    ('Sarbet', 'Addis Ababa', NULL)
ON CONFLICT DO NOTHING;

ALTER TABLE orders ADD COLUMN IF NOT EXISTS delivery_agent_id BIGINT REFERENCES delivery_agents(id);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS delivery_route_id BIGINT REFERENCES delivery_routes(id);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS dispatched_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_orders_delivery_agent_id ON orders(delivery_agent_id);
