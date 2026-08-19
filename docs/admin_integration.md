# Platform Admin Integration Guide

Integration reference for the Gabaa platform-admin API: authentication, store moderation, KYC review, global orders/users, and per-store drill-downs.

## Overview

| Item | Value |
|------|-------|
| Base URL | `https://<host>/api/v1` |
| Admin login | `POST /auth/admin/login` (public, no token) |
| Admin API prefix | `/admin/*` |
| Auth header | `Authorization: Bearer <jwt>` |
| Required JWT role | `platform_admin` |
| Content-Type | `application/json` |

Configure credentials via environment variables:

```env
ADMIN_LOGIN_USERNAME=admin
ADMIN_LOGIN_PASSWORD=changeme
```

Swagger UI (when server is running): `GET /swagger/index.html`

---

## Response Envelope

All endpoints return the same wrapper:

```json
{
  "success": true,
  "data": { },
  "error": null
}
```

**Success example**

```json
{
  "success": true,
  "data": { "token": "eyJ..." },
  "error": null
}
```

**Error example**

```json
{
  "success": false,
  "data": null,
  "error": {
    "error": "FORBIDDEN",
    "message": "Platform admin access required"
  }
}
```

| HTTP status | Error code | Typical cause |
|-------------|------------|---------------|
| 400 | `BAD_REQUEST` | Invalid JSON, query params, or store ID |
| 401 | `UNAUTHORIZED` | Missing/invalid/expired JWT, wrong admin credentials |
| 403 | `FORBIDDEN` | Valid JWT but role is not `platform_admin` |
| 404 | `NOT_FOUND` | Store or resource not found |
| 422 | `VALIDATION_ERROR` | Business rule or field validation failure |
| 500 | `INTERNAL_ERROR` | Server error |

---

## Pagination

All admin **list** endpoints support pagination via query params:

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `page` | int | `1` | Page number (1-based) |
| `page_size` | int | `10` | Items per page |

**Paginated response shape**

```json
{
  "success": true,
  "data": {
    "total": 42,
    "page": 1,
    "page_size": 10,
    "has_next": true,
    "has_previous": false,
    "data": []
  },
  "error": null
}
```

---

## 1. Authentication

### 1.1 Platform admin login

`POST /auth/admin/login`

Public endpoint. Returns a JWT usable on all `/admin/*` routes.

**Request**

```json
{
  "username": "admin",
  "password": "changeme"
}
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "userId": 0,
    "username": "admin",
    "role": "platform_admin",
    "hasStore": false,
    "isDelivery": false
  },
  "error": null
}
```

**Error responses**

| Status | Message |
|--------|---------|
| 400 | Missing or invalid request body |
| 401 | `invalid admin credentials` |

### 1.2 Using the token

Include on every admin request:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Non-admin tokens (merchant, customer, delivery) receive `403 FORBIDDEN` on `/admin/*`.

---

## 2. Store Moderation

### 2.1 List stores

`GET /admin/stores`

Paginated list of all stores with optional filters and text search.

**Query params**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `page_size` | int | Page size |
| `status` | string | `pending`, `launched`, `banned` |
| `verification_status` | string | `unverified`, `pending_review`, `verified`, `rejected` |
| `category` | string | Exact category match (ILIKE) |
| `query` | string | Search name, description, location, phone |

**Example request**

```http
GET /api/v1/admin/stores?page=1&page_size=20&status=launched&query=coffee
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "total": 1,
    "page": 1,
    "page_size": 20,
    "has_next": false,
    "has_previous": false,
    "data": [
      {
        "id": 3,
        "seller_id": 123456789,
        "telegram_chat_id": -1001234567890,
        "telegram_chat_title": "Coffee Shop",
        "status": "launched",
        "verificationStatus": "verified",
        "name": "Bole Coffee",
        "category": "Food & Drink",
        "description": "Specialty coffee",
        "logo_image": "https://res.cloudinary.com/.../logo.jpg",
        "cover_image": "https://res.cloudinary.com/.../cover.jpg",
        "phone": "+251911234567",
        "email": "shop@example.com",
        "location": "Addis Ababa, Bole"
      }
    ]
  },
  "error": null
}
```

---

### 2.2 Get store detail

`GET /admin/stores/:store_id`

Returns a single store profile (not paginated).

**Path params**

| Param | Type | Description |
|-------|------|-------------|
| `store_id` | int | Store ID |

**Example request**

```http
GET /api/v1/admin/stores/3
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "id": 3,
    "seller_id": 123456789,
    "telegram_chat_id": -1001234567890,
    "telegram_chat_title": "Coffee Shop",
    "status": "launched",
    "verificationStatus": "verified",
    "name": "Bole Coffee",
    "category": "Food & Drink",
    "description": "Specialty coffee",
    "logo_image": "https://res.cloudinary.com/.../logo.jpg",
    "cover_image": "https://res.cloudinary.com/.../cover.jpg",
    "phone": "+251911234567",
    "email": "shop@example.com",
    "location": "Addis Ababa, Bole"
  },
  "error": null
}
```

**Error responses**

| Status | Message |
|--------|---------|
| 400 | `invalid store id` |
| 404 | `store not found` |

---

### 2.3 Update store status

`PATCH /admin/stores/:store_id/status`

Ban, unban, or change store lifecycle status.

**Path params**

| Param | Type | Description |
|-------|------|-------------|
| `store_id` | int | Store ID |

**Request body**

```json
{
  "status": "banned"
}
```

Allowed values: `pending`, `launched`, `banned`.

**Example request**

```http
PATCH /api/v1/admin/stores/3/status
Authorization: Bearer <token>
Content-Type: application/json

{"status": "banned"}
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "id": 3,
    "seller_id": 123456789,
    "telegram_chat_id": -1001234567890,
    "telegram_chat_title": "Coffee Shop",
    "status": "banned",
    "verificationStatus": "verified",
    "name": "Bole Coffee",
    "category": "Food & Drink",
    "description": "Specialty coffee",
    "logo_image": "https://res.cloudinary.com/.../logo.jpg",
    "cover_image": "https://res.cloudinary.com/.../cover.jpg",
    "phone": "+251911234567",
    "email": "shop@example.com",
    "location": "Addis Ababa, Bole"
  },
  "error": null
}
```

**Error responses**

| Status | Message |
|--------|---------|
| 400 | Invalid status or store ID |
| 422 | `invalid store status` |

---

## 3. Store KYC

### 3.1 List KYC verifications

`GET /admin/store-verifications`

Paginated queue of store KYC submissions.

**Query params**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `page_size` | int | Page size |
| `status` | string | `pending_review` (default), `verified`, `rejected` |
| `query` | string | Search store name or TIN |

**Example request**

```http
GET /api/v1/admin/store-verifications?page=1&page_size=10&status=pending_review&query=1234567890
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "total": 1,
    "page": 1,
    "page_size": 10,
    "has_next": false,
    "has_previous": false,
    "data": [
      {
        "storeId": 3,
        "storeName": "Bole Coffee",
        "verificationStatus": "pending_review",
        "tinNumber": "1234567890",
        "businessRegistrationNumber": "BR-2024-001",
        "tinCertificateUrl": "https://res.cloudinary.com/.../tin.pdf",
        "businessLicenseUrl": "https://res.cloudinary.com/.../license.pdf",
        "submittedAt": "2026-08-15T10:30:00Z"
      }
    ]
  },
  "error": null
}
```

---

### 3.2 Upsert store KYC (admin)

`POST /admin/stores/:store_id/kyc`

Create or update KYC documents on behalf of a store. Sets verification status to `pending_review`.

**Path params**

| Param | Type | Description |
|-------|------|-------------|
| `store_id` | int | Store ID |

**Request body**

```json
{
  "tinNumber": "1234567890",
  "businessRegistrationNumber": "BR-2024-001",
  "tinCertificateUrl": "https://res.cloudinary.com/.../tin.pdf",
  "businessLicenseUrl": "https://res.cloudinary.com/.../license.pdf"
}
```

All fields are required. Certificate URLs must be valid URLs.

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "storeId": 3,
    "storeName": "Bole Coffee",
    "verificationStatus": "pending_review",
    "tinNumber": "1234567890",
    "businessRegistrationNumber": "BR-2024-001",
    "tinCertificateUrl": "https://res.cloudinary.com/.../tin.pdf",
    "businessLicenseUrl": "https://res.cloudinary.com/.../license.pdf",
    "submittedAt": "2026-08-19T18:00:00Z"
  },
  "error": null
}
```

---

### 3.3 Approve KYC

`POST /admin/store-verifications/:store_id/approve`

Marks the store as verified.

**Path params**

| Param | Type | Description |
|-------|------|-------------|
| `store_id` | int | Store ID |

**Request body**

None.

**Example request**

```http
POST /api/v1/admin/store-verifications/3/approve
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "storeId": 3,
    "storeName": "Bole Coffee",
    "verificationStatus": "verified",
    "tinNumber": "1234567890",
    "businessRegistrationNumber": "BR-2024-001",
    "tinCertificateUrl": "https://res.cloudinary.com/.../tin.pdf",
    "businessLicenseUrl": "https://res.cloudinary.com/.../license.pdf",
    "submittedAt": "2026-08-15T10:30:00Z",
    "reviewedAt": "2026-08-19T18:05:00Z"
  },
  "error": null
}
```

---

### 3.4 Reject KYC

`POST /admin/store-verifications/:store_id/reject`

Rejects KYC with an optional review note.

**Path params**

| Param | Type | Description |
|-------|------|-------------|
| `store_id` | int | Store ID |

**Request body**

```json
{
  "reviewNote": "TIN certificate is illegible. Please re-upload."
}
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "storeId": 3,
    "storeName": "Bole Coffee",
    "verificationStatus": "rejected",
    "tinNumber": "1234567890",
    "businessRegistrationNumber": "BR-2024-001",
    "reviewNote": "TIN certificate is illegible. Please re-upload.",
    "submittedAt": "2026-08-15T10:30:00Z",
    "reviewedAt": "2026-08-19T18:10:00Z"
  },
  "error": null
}
```

---

## 4. Orders

### 4.1 List global orders

`GET /admin/orders`

Cross-store order list for platform oversight.

**Query params**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `page_size` | int | Page size |
| `store_id` | int | Filter by store |
| `order_id` | int | Exact order ID |
| `status` | string | Order status filter |
| `query` | string | Search customer username (ILIKE) |

**Example request**

```http
GET /api/v1/admin/orders?page=1&page_size=20&store_id=3&status=shipped&query=john
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "total": 1,
    "page": 1,
    "page_size": 20,
    "has_next": false,
    "has_previous": false,
    "data": [
      {
        "id": 101,
        "store_id": 3,
        "user_id": 55,
        "status": "shipped",
        "total_price": 450.0,
        "created_at": "2026-08-18T14:22:00Z",
        "customer": {
          "id": 55,
          "username": "john_doe"
        },
        "shipping_address": {
          "id": 12,
          "user_id": 55,
          "label": "home",
          "recipient_name": "John Doe",
          "phone": "+251911234567",
          "street": "Bole Road",
          "city": "Addis Ababa",
          "region": "Bole",
          "country": "Ethiopia",
          "is_default": true,
          "created_at": "2026-07-01T08:00:00Z"
        },
        "order_items": [
          {
            "id": 201,
            "order_id": 101,
            "product_id": 7,
            "quantity": 2,
            "price": 225.0,
            "product": {
              "id": 7,
              "name": "Ethiopian Blend",
              "images": ["https://res.cloudinary.com/.../coffee.jpg"],
              "category": "Coffee"
            }
          }
        ]
      }
    ]
  },
  "error": null
}
```

---

### 4.2 List store orders

`GET /admin/stores/:store_id/orders`

Same order shape as global list, scoped to one store. `store_id` from the path is always applied; do not rely on a query `store_id` override.

**Query params**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `page_size` | int | Page size |
| `order_id` | int | Exact order ID |
| `status` | string | Order status filter |
| `query` | string | Search customer username |

**Example request**

```http
GET /api/v1/admin/stores/3/orders?page=1&page_size=10&status=delivered
Authorization: Bearer <token>
```

**Success response — `200`**

Same paginated `Order` array structure as [§4.1](#41-list-global-orders).

---

## 5. Users

### 5.1 List users

`GET /admin/users`

Paginated list of all platform users.

**Query params**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `page_size` | int | Page size |
| `role` | string | Filter by role (e.g. `admin`, `customer`) |
| `query` | string | Search username (ILIKE) or exact Telegram user ID |

**Example request**

```http
GET /api/v1/admin/users?page=1&page_size=20&role=admin&query=seller
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "total": 1,
    "page": 1,
    "page_size": 20,
    "has_next": false,
    "has_previous": false,
    "data": [
      {
        "id": 1,
        "telegram_user_id": 123456789,
        "email": "",
        "username": "sample_seller",
        "role": "admin"
      }
    ]
  },
  "error": null
}
```

---

## 6. Store Drill-Down

These endpoints load data for a single store. Start from [§2.1 List stores](#21-list-stores) or [§2.2 Get store detail](#22-get-store-detail), then use the store ID below.

### 6.1 Wallet summary

`GET /admin/stores/:store_id/wallet`

**Example request**

```http
GET /api/v1/admin/stores/3/wallet
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "store_id": 3,
    "currency": "ETB",
    "pending_balance": 1200.0,
    "available_balance": 8500.0,
    "locked_balance": 500.0,
    "total_earned": 25000.0,
    "total_withdrawn": 15000.0
  },
  "error": null
}
```

---

### 6.2 Wallet withdrawals

`GET /admin/stores/:store_id/wallet/withdrawals`

**Query params**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `page_size` | int | Page size |
| `status` | string | `initiated`, `pending`, `success`, `failed`, `cancelled` |

**Example request**

```http
GET /api/v1/admin/stores/3/wallet/withdrawals?page=1&page_size=10&status=success
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "total": 1,
    "page": 1,
    "page_size": 10,
    "has_next": false,
    "has_previous": false,
    "data": [
      {
        "id": 5,
        "store_id": 3,
        "amount": 2000.0,
        "currency": "ETB",
        "phone_number": "+251911234567",
        "medium": "TELEBIRR",
        "reference": "WD-20260818-001",
        "transaction_id": "LAKI-TXN-98765",
        "status": "success",
        "gateway_status": "SUCCESS",
        "created_at": "2026-08-18T09:00:00Z"
      }
    ]
  },
  "error": null
}
```

---

### 6.3 Payment transactions

`GET /admin/stores/:store_id/transactions`

Customer checkout payments for the store.

**Query params**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `page_size` | int | Page size |
| `status` | string | `initiated`, `pending`, `success`, `failed` |
| `medium` | string | `MPESA`, `TELEBIRR`, `CBE`, `ETHSWITCH` |
| `query` | string | Search payment reference or transaction ID |

**Example request**

```http
GET /api/v1/admin/stores/3/transactions?page=1&page_size=10&status=success&medium=TELEBIRR
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "total": 1,
    "page": 1,
    "page_size": 10,
    "has_next": false,
    "has_previous": false,
    "data": [
      {
        "id": 88,
        "order_id": 101,
        "status": "success",
        "method": "mobile_money",
        "reference": "PAY-20260818-001",
        "transaction_id": "LAKI-TXN-12345",
        "amount": 450.0,
        "currency": "ETB",
        "phone_number": "+251911234567",
        "medium": "TELEBIRR",
        "gateway_status": "SUCCESS",
        "created_at": "2026-08-18T14:22:00Z",
        "order_status": "shipped"
      }
    ]
  },
  "error": null
}
```

---

### 6.4 Story ads

`GET /admin/stores/:store_id/stories`

**Query params**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `page_size` | int | Page size |
| `is_active` | bool | Filter active/inactive stories |
| `type` | string | `image` or `video` |
| `query` | string | Search story caption |

**Example request**

```http
GET /api/v1/admin/stores/3/stories?page=1&page_size=10&is_active=true&type=image
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "total": 1,
    "page": 1,
    "page_size": 10,
    "has_next": false,
    "has_previous": false,
    "data": [
      {
        "id": 15,
        "store_id": 3,
        "product_id": 7,
        "caption": "Fresh roast this week!",
        "media_urls": ["https://res.cloudinary.com/.../story.jpg"],
        "media_type": "image",
        "starts_at": "2026-08-19T00:00:00Z",
        "ends_at": "2026-08-26T23:59:59Z",
        "is_active": true,
        "views": 142,
        "created_at": "2026-08-19T08:00:00Z"
      }
    ]
  },
  "error": null
}
```

---

### 6.5 Delivery agents

`GET /admin/stores/:store_id/deliveries`

Connected delivery agents and their routes for the store.

**Query params**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `page_size` | int | Page size |
| `status` | string | Agent link status (e.g. `active`, `pending_invite`) |
| `query` | string | Search agent username or full name |

**Example request**

```http
GET /api/v1/admin/stores/3/deliveries?page=1&page_size=10&status=active
Authorization: Bearer <token>
```

**Success response — `200`**

```json
{
  "success": true,
  "data": {
    "total": 1,
    "page": 1,
    "page_size": 10,
    "has_next": false,
    "has_previous": false,
    "data": [
      {
        "id": 2,
        "username": "courier_john",
        "full_name": "John Doe",
        "phone": "251911234567",
        "status": "active",
        "loyalty_score": 85,
        "share_enabled": true,
        "routes": [
          {
            "id": 4,
            "label": "Bole + Atlas Route",
            "is_active": true,
            "pickup_locations": [
              {
                "id": 10,
                "location_type": "pickup",
                "label": "My store",
                "country": "",
                "region": "",
                "city": "",
                "street": "",
                "landmark": "",
                "notes": "",
                "use_store_location": true
              }
            ],
            "delivery_locations": [
              {
                "id": 11,
                "location_type": "delivery",
                "label": "Bole area",
                "country": "Ethiopia",
                "region": "Bole",
                "city": "Addis Ababa",
                "street": "",
                "landmark": "",
                "notes": "",
                "use_store_location": false
              }
            ]
          }
        ]
      }
    ]
  },
  "error": null
}
```

---

## Endpoint Quick Reference

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/auth/admin/login` | No | Get platform admin JWT |
| GET | `/admin/stores` | Yes | List stores (paginated) |
| GET | `/admin/stores/:store_id` | Yes | Store detail |
| PATCH | `/admin/stores/:store_id/status` | Yes | Update store status |
| POST | `/admin/stores/:store_id/kyc` | Yes | Upsert store KYC |
| GET | `/admin/stores/:store_id/orders` | Yes | Store orders (paginated) |
| GET | `/admin/stores/:store_id/wallet` | Yes | Wallet summary |
| GET | `/admin/stores/:store_id/wallet/withdrawals` | Yes | Withdrawals (paginated) |
| GET | `/admin/stores/:store_id/transactions` | Yes | Payment transactions (paginated) |
| GET | `/admin/stores/:store_id/stories` | Yes | Story ads (paginated) |
| GET | `/admin/stores/:store_id/deliveries` | Yes | Delivery agents (paginated) |
| GET | `/admin/store-verifications` | Yes | KYC queue (paginated) |
| POST | `/admin/store-verifications/:store_id/approve` | Yes | Approve KYC |
| POST | `/admin/store-verifications/:store_id/reject` | Yes | Reject KYC |
| GET | `/admin/orders` | Yes | Global orders (paginated) |
| GET | `/admin/users` | Yes | All users (paginated) |

---

## Typical Admin Flow

```mermaid
sequenceDiagram
    participant Admin as Admin Dashboard
    participant API as Gabaa API

    Admin->>API: POST /auth/admin/login
    API-->>Admin: JWT (role=platform_admin)

    Admin->>API: GET /admin/store-verifications?status=pending_review
    API-->>Admin: Paginated KYC queue

    Admin->>API: POST /admin/store-verifications/3/approve
    API-->>Admin: StoreKYCResponse (verified)

    Admin->>API: GET /admin/stores/3
    API-->>Admin: Store detail

    Admin->>API: GET /admin/stores/3/orders
    API-->>Admin: Paginated orders

    Admin->>API: GET /admin/stores/3/wallet
    API-->>Admin: Wallet balances

    Admin->>API: PATCH /admin/stores/5/status {"status":"banned"}
    API-->>Admin: Updated store
```

1. **Login** — `POST /auth/admin/login`, store the `token`.
2. **Review KYC** — `GET /admin/store-verifications`, then approve or reject.
3. **Moderate stores** — list, inspect detail, ban/unban via status patch.
4. **Investigate** — drill into orders, wallet, transactions, stories, deliveries per store.
5. **Platform overview** — `GET /admin/orders` and `GET /admin/users` for cross-store visibility.

---

## cURL Examples

**Login**

```bash
curl -s -X POST "https://<host>/api/v1/auth/admin/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"changeme"}'
```

**List stores**

```bash
TOKEN="<jwt>"
curl -s "https://<host>/api/v1/admin/stores?page=1&page_size=20&status=launched" \
  -H "Authorization: Bearer $TOKEN"
```

**Ban a store**

```bash
curl -s -X PATCH "https://<host>/api/v1/admin/stores/3/status" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"banned"}'
```

**Approve KYC**

```bash
curl -s -X POST "https://<host>/api/v1/admin/store-verifications/3/approve" \
  -H "Authorization: Bearer $TOKEN"
```
