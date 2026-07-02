# Product Recommendation Guide

This guide explains how customers subscribe to category-based product recommendations and how the bot delivers alerts when matching products are published.

---

## Overview

When a merchant publishes a product, Gabaa Place can notify opted-in customers in their **private chat** with @GabaaPlaceBot if:

1. The customer has enabled recommendations
2. The customer has started the bot (`/start`)
3. The product category matches one of the customer's saved preferences
4. The customer is not the seller who published the product
5. The customer has not already been notified about that product

Notifications run in the **background** and do not block the publish API response.

---

## Customer Setup

### Step 1 — Start the bot

Customers must open @GabaaPlaceBot and send `/start` at least once. Telegram only allows private messages to users who have initiated a chat with the bot.

### Step 2 — Choose categories

**Option A: Bot commands**

```
/preferences
```

Tap categories in the inline keyboard to subscribe (✓ = selected).

```
/recommendations on
/recommendations off
```

**Option B: Mini App API**

Requires JWT from `POST /auth/telegram`.

#### Get preferences

```
GET /user/preferences
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "enabled": true,
    "categories": ["Electronics", "Fashion"]
  },
  "error": null
}
```

#### Update preferences

```
PUT /user/preferences
Authorization: Bearer <token>
Content-Type: application/json
```

**Request:**
```json
{
  "enabled": true,
  "categories": ["Electronics", "Beauty"]
}
```

- `categories` replaces the full list (max 10)
- `enabled` controls whether alerts are sent
- Category matching is **case-insensitive**

---

## What Customers Receive

When a matching product is published, the bot sends a private message with:

- Product name, category, price (ETB), and store name
- Product image (first image, if available)
- **Order Now** button linking to the Mini App (`?startapp=product_{id}`)
- Footer with `/preferences` reminder

---

## Merchant Behavior

Publishing a product (`PUT /my-store/product/:id` with `"status": "published"`) triggers:

1. Existing channel/group product post (unchanged)
2. Background recommendation DMs to matching customers

The seller who published the product never receives a recommendation for their own listing.

---

## Matching Rules

| Rule | Detail |
|------|--------|
| Scope | Global — any store on the platform |
| Category field | Matches `products.category` string |
| Opt-in | `recommendations_enabled` must be `true` |
| Bot access | `bot_started` must be `true` |
| Dedup | One notification per user per product |
| Rate limit | ~30ms delay between sends |

---

## Database Tables

| Table | Purpose |
|-------|---------|
| `user_category_preferences` | One row per user; `categories` JSONB list of subscribed category names |
| `users.recommendations_enabled` | Opt-in flag |
| `users.bot_started` | Telegram `/start` tracking |
| `product_recommendations` | Sent notification log (dedup) |

---

## Mini App UI Checklist

- [ ] Preferences screen with category multi-select
- [ ] Toggle for "Enable product recommendations"
- [ ] CTA: "Open @GabaaPlaceBot and tap Start" if alerts are enabled but bot not started
- [ ] Link to bot: `https://t.me/GabaaPlaceBot`

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| No alerts received | Ensure `/recommendations on`, categories selected, and `/start` sent to bot |
| Alerts stopped | User may have blocked the bot; check server logs for 403 errors |
| Wrong categories | Categories must match product `category` field (e.g. `Fashion`, not `fashion` — matching is case-insensitive) |
| Duplicate alerts | Should not occur; `product_recommendations` prevents resends for same product |

---

## API Summary

| Method | Route | Auth | Description |
|--------|-------|------|-------------|
| `GET` | `/user/preferences` | JWT | Get preferences |
| `PUT` | `/user/preferences` | JWT | Update preferences |

## Bot Commands

| Command | Description |
|---------|-------------|
| `/start` | Register bot access for private messages |
| `/preferences` | Category picker (inline keyboard) |
| `/recommendations on` | Enable alerts |
| `/recommendations off` | Disable alerts |
