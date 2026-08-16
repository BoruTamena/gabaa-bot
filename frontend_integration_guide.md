# Frontend Integration Guide: Public API

Public marketplace and storefront APIs. All paths are under `/api/v1`.

## Base URL
Relative to the API host (e.g. `https://api.yourdomain.com`). Prefix: `/api/v1`.

## Suggested frontend flow

1. Browse **active stores** → pick a store  
2. Load **store details** by name  
3. On the store page: **products**, **stories**, **sales**  
4. Global discovery: **products**, **stories**, **categories**  
5. Detail views: **product by ID**, **story by ID**

## Endpoint index

| # | Method | Path | Description |
|---|--------|------|-------------|
| 1 | GET | `/stores` | List active (launched) stores |
| 2 | GET | `/stores/{storeName}` | Store public profile |
| 3 | GET | `/stores/{storeName}/products` | Products for that store |
| 4 | GET | `/stores/{storeName}/stories` | Stories for that store |
| 5 | GET | `/stores/{storeName}/sales` | Recent sales / social proof |
| 6 | GET | `/products` | Global product search |
| 7 | GET | `/product/{id}` | Single product |
| 8 | GET | `/stories` | Global story list / search |
| 9 | GET | `/stories/active` | Stories in current date window |
| 10 | GET | `/stories/{id}` | Single story |
| 11 | GET | `/categories` | Global categories |

Auth: none of the endpoints below require a JWT.

---

## 1. List Active Stores
Paginated list of stores with `status = launched`. Use this for the marketplace home / store directory.

**Endpoint:** `GET /api/v1/stores`

### Query Parameters
- `page` (int, optional): Page number (default: 1)
- `page_size` (int, optional): Items per page (default: 10)
- `query` (string, optional): Search store name, description, or location (ILIKE)
- `category` (string, optional): Filter by category (ILIKE)

### Examples
- `GET /api/v1/stores?page=1&page_size=12`
- `GET /api/v1/stores?query=bole&category=Electronics`

### Success Response
**Status:** `200 OK`
```json
{
  "success": true,
  "data": {
    "total": 24,
    "page": 1,
    "page_size": 12,
    "has_next": true,
    "has_previous": false,
    "data": [
      {
        "id": 12,
        "seller_id": 5,
        "name": "gabaa",
        "category": "Electronics",
        "description": "Premium goods and accessories",
        "logo_image": "https://...",
        "cover_image": "https://...",
        "phone": "+2519...",
        "email": "shop@example.com",
        "location": "Addis Ababa, Bole",
        "status": "launched",
        "verificationStatus": "verified",
        "telegram_chat_id": 123456789,
        "telegram_chat_title": "Gabaa Shop"
      }
    ]
  },
  "error": null
}
```

---

## 2. Get Store Details
Public profile for one store. Lookup is **case-insensitive** by name.

**Endpoint:** `GET /api/v1/stores/{storeName}`

### Path Parameters
- `storeName` (string, required)

### Success Response
**Status:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 12,
    "name": "gabaa",
    "category": "Electronics",
    "description": "Premium goods and accessories",
    "logo_image": "https://...",
    "cover_image": "https://...",
    "location": "Addis Ababa, Bole",
    "status": "launched",
    "verificationStatus": "verified",
    "telegram_chat_id": 123456789
  },
  "error": null
}
```

### Errors
- `400` — store name missing  
- `404` — store not found  

---

## 3. Get Store Products
Published products for a store (by name).

**Endpoint:** `GET /api/v1/stores/{storeName}/products`

### Path Parameters
- `storeName` (string, required)

### Query Parameters
- `page` / `page_size` — pagination  
- `category` — filter by category  
- `query` — search name or description  

### Success Response
Same paginated product shape as section 6.

---

## 4. Get Store Stories
Active stories for a store. The path `storeName` is treated as an **ILIKE search term** on store name (empty results if none match — not a 404).

**Endpoint:** `GET /api/v1/stores/{storeName}/stories`

### Query Parameters
- `page` / `page_size`  
- `type` — `image` | `video`  
- `search` — caption text  
- `sort_by` — `newest` | `popular`  

### Success Response
Same paginated story shape as section 8.

---

## 5. Get Store Recent Sales
Social proof for the store “Sells” tab.

**Endpoint:** `GET /api/v1/stores/{storeName}/sales`

### Query Parameters
- `page` / `page_size`  
- `status` (optional), e.g. `delivered`  

### Success Response
**Status:** `200 OK`
```json
{
  "success": true,
  "data": {
    "total": 15,
    "data": [
      {
        "id": "ORD-1052",
        "productName": "Nike Air Max",
        "amount": 2500,
        "currency": "ETB",
        "status": "delivered",
        "buyerName": "johndoe",
        "purchasedAt": "2026-08-13T09:30:00Z"
      }
    ]
  },
  "error": null
}
```

---

## 6. List / Search Public Products
Published products across all stores (or filtered to one).

**Endpoint:** `GET /api/v1/products`

### Query Parameters
- `page` / `page_size`  
- `store_id` — filter by store ID (from section 1/2)  
- `title` — ILIKE on product name  
- `min_price` / `max_price` — inclusive range  
- `category` — category filter  
- `query` — name **or** description  

### Examples
- `GET /api/v1/products?page=1&page_size=10`
- `GET /api/v1/products?store_id=12&title=nike&min_price=100&max_price=500`

### Success Response
**Status:** `200 OK`
```json
{
  "success": true,
  "data": {
    "total": 50,
    "page": 1,
    "page_size": 10,
    "has_next": true,
    "has_previous": false,
    "data": [
      {
        "id": 100,
        "store_id": 12,
        "seller_id": 5,
        "name": "Nike Air Max",
        "description": "Classic running shoes",
        "price": 2500,
        "stock": 10,
        "category": "Shoes",
        "images": ["https://..."],
        "status": "published",
        "is_posted": true,
        "is_boosted": false
      }
    ]
  },
  "error": null
}
```

---

## 7. Get Product by ID
**Endpoint:** `GET /api/v1/product/{id}`

### Path Parameters
- `id` (int, required)

### Success Response
Single `Product` object (same fields as list items).

---

## 8. List / Search Public Stories
Active stories across stores; optional store-name text search.

**Endpoint:** `GET /api/v1/stories`

### Query Parameters
- `page` / `page_size`  
- `store_name` — ILIKE on store name  
- `type` — `image` | `video`  
- `search` — caption  
- `sort_by` — `newest` | `popular`  

### Examples
- `GET /api/v1/stories?page=1&page_size=10`
- `GET /api/v1/stories?store_name=bole`

### Success Response
**Status:** `200 OK`
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
        "id": 45,
        "store_id": 12,
        "product_id": 100,
        "caption": "Limited time offer on these sneakers!",
        "media_urls": ["https://..."],
        "media_type": "image",
        "starts_at": "2026-08-13T00:00:00Z",
        "ends_at": "2026-08-15T23:59:59Z",
        "is_active": true,
        "views": 150,
        "created_at": "2026-08-13T10:00:00Z",
        "product": {
          "id": 100,
          "name": "Nike Air Max",
          "price": 2500,
          "images": ["https://..."],
          "category": "Shoes"
        }
      }
    ]
  },
  "error": null
}
```

---

## 9. Active Story Feed (date window)
Stories that are active **and** currently within `starts_at` … `ends_at`.

**Endpoint:** `GET /api/v1/stories/active`

### Query Parameters
- `page` / `page_size`

Pagination shape matches section 8.

---

## 10. Get Single Story
Fetches one story and increments its view count.

**Endpoint:** `GET /api/v1/stories/{id}`

### Path Parameters
- `id` (int, required)

---

## 11. List Categories
For filters and setup forms.

**Endpoint:** `GET /api/v1/categories`

### Query Parameters
- `page` / `page_size`

### Success Response
```json
{
  "success": true,
  "data": {
    "total": 5,
    "data": [
      { "id": 1, "name": "Electronics" },
      { "id": 2, "name": "Fashion" }
    ]
  },
  "error": null
}
```

---

## Standard Error Format

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "NOT_FOUND",
    "message": "Store not found"
  }
}
```

## Implementation Notes
- **Active stores** means `status = launched` only (pending stores are hidden).
- Prefer `has_next` / `has_previous` when present; otherwise `Math.ceil(total / page_size)`.
- Public product lists only return `status = published`.
- Use `title` for product name-only search; use `query` for name + description.
- Store path lookups are case-insensitive (`Gabaa` and `gabaa` resolve the same).
- Register static paths like `/stores` and `/stories/active` before parameterized routes on the client routers as well.
