DROP INDEX IF EXISTS idx_orders_delivery_agent_id;
ALTER TABLE orders DROP COLUMN IF EXISTS dispatched_at;
ALTER TABLE orders DROP COLUMN IF EXISTS delivery_route_id;
ALTER TABLE orders DROP COLUMN IF EXISTS delivery_agent_id;

DROP TABLE IF EXISTS delivery_area_presets;
DROP TABLE IF EXISTS delivery_agent_shares;
DROP TABLE IF EXISTS delivery_route_locations;
DROP TABLE IF EXISTS delivery_routes;
DROP TABLE IF EXISTS store_delivery_links;
DROP TABLE IF EXISTS delivery_agents;
