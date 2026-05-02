# Merchant Order Integration Guide

This guide covers all order management endpoints for the merchant dashboard. All endpoints are protected and require a valid JWT token. The `store_id` is automatically extracted from the token — no URL parameter needed.

---

## Authentication

```
Authorization: Bearer <your_token>
```

---

## Order Statuses (Lifecycle)

| Status | Description |
| :--- | :--- |
| `pending` | Order placed by customer, awaiting action |
| `shipped` | Merchant has dispatched the order |
| `delivered` | Order received by customer; wallet credited |
| `cancelled` | Order was cancelled by the customer |

---

## 1. List My Store Orders

Fetch a paginated list of all orders for the authenticated merchant's store.

**Endpoint:** `GET /my-store/orders`

### Query Parameters

| Param | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `order_id` | integer | No | Search by exact order ID |
| `status` | string | No | Filter: `pending`, `shipped`, `delivered`, `cancelled` |
| `page` | integer | No | Page number (default: 1) |
| `page_size` | integer | No | Items per page (default: 10) |

### Examples

```
GET /my-store/orders
GET /my-store/orders?status=pending
GET /my-store/orders?status=shipped&page=2
GET /my-store/orders?order_id=42
```

### Response

```json
{
  "success": true,
  "data": {
    "total": 25,
    "data": [
      {
        "id": 42,
        "store_id": 2,
        "user_id": 7,
        "status": "pending",
        "total_price": 3500.00,
        "created_at": "2026-05-02T20:15:00Z",
        "customer": {
          "id": 7,
          "username": "BoruTamena"
        },
        "shipping_address": {
          "id": 42,
          "label": "home",
          "recipient_name": "Boru Tamena",
          "phone": "+251911000000",
          "street": "Bole Road, XYZ Building",
          "city": "Addis Ababa",
          "country": "Ethiopia"
        },
        "order_items": [
          {
            "id": 1,
            "order_id": 42,
            "product_id": 3,
            "quantity": 2,
            "price": 1750.00,
            "product": {
              "id": 3,
              "name": "Vintage Leather Bag",
              "images": ["https://res.cloudinary.com/.../image1.jpg"],
              "category": "Fashion"
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

## 2. Get Order Detail

Retrieve the full detail of a single order. Returns complete item breakdown, product info, and customer identity.

**Endpoint:** `GET /my-store/orders/:order_id`

> **Security:** Returns `403 Forbidden` if the order belongs to a different store.

### Response

```json
{
  "success": true,
  "data": {
    "id": 42,
    "store_id": 2,
    "user_id": 7,
    "status": "pending",
    "total_price": 3500.00,
    "created_at": "2026-05-02T20:15:00Z",
    "customer": {
      "id": 7,
      "username": "BoruTamena"
    },
    "shipping_address": {
      "id": 42,
      "label": "home",
      "recipient_name": "Boru Tamena",
      "phone": "+251911000000",
      "street": "Bole Road, XYZ Building",
      "city": "Addis Ababa",
      "country": "Ethiopia"
    },
    "order_items": [
      {
        "id": 1,
        "order_id": 42,
        "product_id": 3,
        "quantity": 2,
        "price": 1750.00,
        "product": {
          "id": 3,
          "name": "Vintage Leather Bag",
          "images": ["https://res.cloudinary.com/.../image1.jpg"],
          "category": "Fashion"
        }
      }
    ]
  },
  "error": null
}
```

---

## 3. Update Order Status

Advance an order to `shipped` or `delivered`. Only these two values are accepted by the merchant.

> **💰 Wallet Credit:** When you set an order to `delivered`, the order's `total_price` is **automatically credited** to your store wallet.

**Endpoint:** `PUT /my-store/orders/:order_id/status`

**Request Body:**
```json
{
  "status": "shipped"
}
```

Valid values: `shipped`, `delivered`

**Success Response:**
```json
{
  "success": true,
  "data": { "message": "order status updated to shipped" },
  "error": null
}
```

**Error — Invalid Status (400):**
```json
{
  "success": false,
  "data": null,
  "error": "invalid status: must be 'shipped' or 'delivered'"
}
```

**Error — Wrong Store (403):**
```json
{
  "success": false,
  "data": null,
  "error": "order does not belong to your store"
}
```

---

## 4. Recommended Implementation Flow

```
Order placed by customer
        ↓
status: "pending"  ← Merchant sees this in the dashboard
        ↓
Merchant ships the goods
        ↓
PUT /my-store/orders/:id/status  { "status": "shipped" }
        ↓
status: "shipped"  ← Customer can see their order is on the way
        ↓
Customer confirms receipt
        ↓
PUT /my-store/orders/:id/status  { "status": "delivered" }
        ↓
status: "delivered"  ← Wallet auto-credited with total_price
```

---

## 5. Merchant Dashboard UI Checklist

- [ ] On load, call `GET /my-store/orders?status=pending` to show new orders requiring attention.
- [ ] Display a **badge count** for `pending` orders in the sidebar/nav.
- [ ] Show order list with `customer.username`, `total_price`, `status`, and `created_at`.
- [ ] On order row click, call `GET /my-store/orders/:order_id` to show full item breakdown.
- [ ] Provide "Mark as Shipped" button → calls `PUT` with `{ "status": "shipped" }`.
- [ ] Provide "Mark as Delivered" button → calls `PUT` with `{ "status": "delivered" }` and notify merchant that wallet has been credited.
- [ ] Provide a search input for the `order_id` query param and a status filter dropdown.
- [ ] Hide `shipped`/`delivered` action buttons for orders already in those terminal states.
