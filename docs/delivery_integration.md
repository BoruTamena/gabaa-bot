# Delivery Integration

Store-connected delivery network: merchants connect couriers with route definitions (pickup + delivery locations), auto-suggest couriers on ship, and couriers complete deliveries via Telegram + the delivery WebApp.

## Overview

| Actor | Auth | Base path |
|-------|------|-----------|
| Merchant | JWT (`role: admin`, `store_id` in token) | `/my-store/delivery/*` |
| Delivery person | JWT (`role: delivery`, `delivery_agent_id` in token) | `/delivery/*` |

Run migration `20260715000001_add_delivery_system` before enabling routes.

## Merchant Flow

### 1. Connect a delivery person

`POST /my-store/delivery/agents`

```json
{
  "username": "courier_john",
  "full_name": "John Doe",
  "phone": "251911234567",
  "share_enabled": true,
  "routes": [
    {
      "label": "Bole + Atlas Route",
      "pickup_locations": [
        { "label": "My store", "use_store_location": true },
        {
          "label": "Warehouse",
          "city": "Addis Ababa",
          "region": "Bole",
          "street": "Bole Road",
          "landmark": "Near Edna Mall"
        }
      ],
      "delivery_locations": [
        { "label": "Bole area", "city": "Addis Ababa", "region": "Bole" },
        {
          "label": "Atlas street",
          "city": "Addis Ababa",
          "region": "Bole",
          "street": "Atlas Avenue"
        }
      ]
    }
  ]
}
```

- Creates or reuses agent by Telegram username.
- Agent status is `pending_invite` until they message the bot `/start`.
- If no pickup locations are provided, a default pickup with `use_store_location: true` is added.

### 2. List / update / disconnect

| Method | Path | Description |
|--------|------|-------------|
| GET | `/my-store/delivery/agents` | List connected agents with routes |
| PUT | `/my-store/delivery/agents/:id` | Update agent, `share_enabled`, routes |
| POST | `/my-store/delivery/agents/:id/routes` | Add route |
| PUT | `/my-store/delivery/routes/:route_id` | Update route label/locations |
| DELETE | `/my-store/delivery/routes/:route_id` | Remove route |
| DELETE | `/my-store/delivery/agents/:id` | Disconnect agent |
| GET | `/my-store/delivery/area-presets` | Autocomplete presets (region/city/street) |

### 3. Cross-store sharing

| Method | Path | Description |
|--------|------|-------------|
| GET | `/my-store/delivery/shared-agents` | Browse shareable agents |
| POST | `/my-store/delivery/agents/:id/adopt` | Adopt shared agent (copies routes; `use_store_location` resolves to adopting store) |

Toggle sharing: `PUT /my-store/delivery/agents/:id` with `{ "share_enabled": true }`.

### 4. Dispatch suggestions

For a paid order ready to ship:

`GET /my-store/orders/:order_id/delivery-suggestions`

Returns ranked couriers with `match_summary`, `score`, and `suggested: true` on the top match. Matching uses store location (pickup) and customer shipping address (delivery) against route locations.

### 5. Ship with dispatch

`PUT /my-store/orders/:order_id/status`

```json
{
  "status": "shipped",
  "delivery_agent_id": 12,
  "delivery_route_id": 5
}
```

- `delivery_agent_id` and `delivery_route_id` are optional; when omitted, the top suggestion is used.
- Returns `400` with `no_matching_courier` if no agent matches both pickup and delivery.
- Sets `delivery_agent_id`, `delivery_route_id`, `dispatched_at`, and notifies the courier on Telegram.

Mark delivered (merchant): same endpoint with `{ "status": "delivered" }` — releases escrow as before.

## Delivery Person Flow

### 1. Bot onboarding

When a merchant connects a courier by `@username`, the person must open the bot and send `/start`. The bot:

1. Links `telegram_user_id` and `user_id`
2. Sets agent status to `active`
3. Sends a welcome message with the delivery WebApp link (`startapp=delivery`)

### 2. Login

Use existing Telegram Mini App auth (`POST /auth/telegram` with init data). Users who are **not** store owners and are linked to an active delivery agent receive:

```json
{
  "role": "delivery",
  "deliveryAgentId": 12,
  "isDelivery": true
}
```

JWT includes `delivery_agent_id` claim. Store owners who are also couriers keep `role: admin`.

### 3. APIs

All require `Authorization: Bearer <token>` and `role: delivery`.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/delivery/profile` | Agent info, loyalty, routes |
| GET | `/delivery/orders` | Assigned orders (`?status=shipped`, pagination) |
| GET | `/delivery/orders/:order_id` | Order detail + pickup + delivery address |
| PUT | `/delivery/orders/:order_id/status` | `delivered` or `failed` |

**Delivered:** releases escrow, sets order `delivered`, `loyalty_score += 1`.

**Failed:** sets order `cancelled`, `loyalty_score -= 1` (floor 0).

### 4. Telegram dispatch notification

On ship, couriers receive a message with pickup store/location, delivery address, items, and a WebApp button: `startapp=delivery_order_<orderId>`.

## Frontend Polling

### Merchant dispatch UI

1. Load order detail (`GET /my-store/orders/:id`).
2. If status is `paid`, call `GET /my-store/orders/:id/delivery-suggestions`.
3. Show ranked list; pre-select `suggested: true` entry.
4. On ship, `PUT /my-store/orders/:id/status` with `shipped` + optional `delivery_agent_id` / `delivery_route_id`.
5. Poll order status until `delivered` or `cancelled`.

### Delivery app

1. Auth via Telegram; verify `role === 'delivery'`.
2. Poll `GET /delivery/orders?status=shipped` on an interval (e.g. 15–30s) or on WebApp focus.
3. Open order via deep link `delivery_order_<id>` → `GET /delivery/orders/:id`.
4. On completion, `PUT /delivery/orders/:id/status` with `delivered` or `failed`.

## Loyalty

- `+1` on successful delivery (courier marks `delivered`)
- `-1` on failed delivery (courier marks `failed`)
- Score never displays below `0`

## Location Matching (v1)

Case-insensitive substring match on `street`, `city`, `region`, `landmark`:

| Level | Points |
|-------|--------|
| Street | 40 |
| Landmark | 30 |
| City | 20 |
| Region | 10 |
| Country-only | 5 |

Total suggestion score = pickup score + delivery score + `loyalty_score`.

## Bot Deep Links

| Link | Purpose |
|------|---------|
| `startapp=delivery` | Delivery person home |
| `startapp=delivery_order_<orderId>` | Open specific assigned order |

## Deploy Checklist

1. Apply migration `20260715000001_add_delivery_system`
2. Redeploy API + bot
3. Merchants connect couriers; couriers `/start` the bot
4. Enable dispatch UI on paid orders
