# Store Analytics Integration Documentation

This document describes how to integrate with the new store analytics API endpoints. All endpoints are secured, merchant-scoped, and read-only.

## Authentication & Headers

Every request to the analytics endpoints requires a JSON Web Token (JWT) matching the store owner's account.

- **Header name**: `Authorization`
- **Format**: `Bearer <JWT_TOKEN>`
- **Prerequisite**: The JWT token *must* contain the merchant's `store_id` (injected automatically in context claims during validation) and the user's `role` must be `"admin"`.

---

## Shared Query Parameters

Each analytics endpoint supports optional date-range filtering via query parameters.

| Parameter | Type | Format | Description | Default |
|-----------|------|--------|-------------|---------|
| `from` | string | RFC3339 | Start date for analytical data (inclusive) | `now - 30 days` |
| `to` | string | RFC3339 | End date for analytical data (inclusive) | `now` |

*Example URL with query parameters:*
`/store/analytics/sales?from=2026-07-01T00:00:00Z&to=2026-07-08T23:59:59Z`

---

## Endpoints Specification

### 1. Sales Analytics

Get sales performance metrics, period-by-period revenue breakdowns, and top-performing products.

- **Endpoint**: `GET /store/analytics/sales`
- **Response Format**: `BaseResponse{data=SalesAnalytics}`

#### Success Response Example (`200 OK`)

```json
{
  "success": true,
  "data": {
    "total_revenue": 2248.95,
    "revenue_change_pct": 15.42,
    "total_orders": 35,
    "average_order_value": 64.26,
    "revenue_by_period": [
      {
        "period": "2026-07-01",
        "revenue": 349.98,
        "orders": 5
      },
      {
        "period": "2026-07-02",
        "revenue": 520.12,
        "orders": 8
      }
    ],
    "top_selling_products": [
      {
        "product_id": 42,
        "product_name": "Smartphone X",
        "revenue": 1599.98,
        "units_sold": 2
      }
    ]
  },
  "error": null
}
```

---

### 2. Order Analytics

Analyze order flow, fulfillment state distributions, and cancellation rates.

- **Endpoint**: `GET /store/analytics/orders`
- **Response Format**: `BaseResponse{data=OrderAnalytics}`

#### Success Response Example (`200 OK`)

```json
{
  "success": true,
  "data": {
    "total_orders": 40,
    "orders_by_status": [
      {
        "status": "delivered",
        "count": 28,
        "percentage": 70.0
      },
      {
        "status": "pending",
        "count": 8,
        "percentage": 20.0
      },
      {
        "status": "cancelled",
        "count": 4,
        "percentage": 10.0
      }
    ],
    "average_order_value": 56.22,
    "recent_orders": 12,
    "cancellation_rate_pct": 10.0
  },
  "error": null
}
```

---

### 3. Product Analytics

Track stock statuses, low-stock warnings, and popular catalog items.

- **Endpoint**: `GET /store/analytics/products`
- **Response Format**: `BaseResponse{data=ProductAnalytics}`

#### Success Response Example (`200 OK`)

```json
{
  "success": true,
  "data": {
    "total_products": 120,
    "products_by_status": [
      {
        "status": "published",
        "count": 95
      },
      {
        "status": "draft",
        "count": 25
      }
    ],
    "low_stock_count": 4,
    "out_of_stock_count": 1,
    "top_viewed_products": [
      {
        "product_id": 10,
        "product_name": "Wireless Headphones",
        "views": 155,
        "category": "Electronics"
      }
    ]
  },
  "error": null
}
```

---

### 4. Stories Analytics

Track viewing engagement for product story ads created by the store.

- **Endpoint**: `GET /store/analytics/stories`
- **Response Format**: `BaseResponse{data=StoryAnalytics}`

#### Success Response Example (`200 OK`)

```json
{
  "success": true,
  "data": {
    "total_stories": 15,
    "active_stories": 3,
    "expired_stories": 12,
    "total_views": 845,
    "top_stories": [
      {
        "story_id": 3,
        "product_id": 12,
        "caption": "Summer Sale 20% Off!",
        "views": 320,
        "starts_at": "2026-07-01T00:00:00Z",
        "ends_at": "2026-07-08T00:00:00Z"
      }
    ]
  },
  "error": null
}
```

---

## Error Handling Specifications

When errors occur, they follow a standard error response body format.

### 401 Unauthorized

Returned when the request lacks a valid JWT token.

```json
{
  "success": false,
  "data": null,
  "error": {
    "error": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  }
}
```

### 403 Forbidden

Returned when the user is authenticated but not registered with an active store or role is not "admin".

```json
{
  "success": false,
  "data": null,
  "error": {
    "error": "FORBIDDEN",
    "message": "Merchant access required"
  }
}
```
