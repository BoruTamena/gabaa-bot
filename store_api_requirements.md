# Backend API Requirements for Store Page

This document outlines the API endpoints required to fully populate the Store Details page (`/stores/[storeName]`) and its associated tabs (Products, Stories, and Sells).

---

## 1. Get Store Details
Retrieves the general information for a specific store based on the store's name or slug.

**Endpoint:** `GET /api/v1/stores/{storeName}`

### Request
- **Path Parameters:**
  - `storeName` (string): The slug or name of the store.

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "id": "str-12345",
    "name": "TechHaven",
    "description": "Premium electronics and gadgets for everyday use. We provide the highest quality products with exceptional customer service. Based in the heart of the city, we deliver nationwide.",
    "avatar": "https://example.com/avatar.jpg",
    "cover": "https://example.com/cover.jpg",
    "rating": 4.8,
    "reviews": 124,
    "location": "Addis Ababa",
    "joinedDate": "2023-01-15T00:00:00Z",
    "categories": ["Electronics", "Gadgets", "Accessories"]
  }
}
```

---

## 2. Get Store Stories (with Pagination & Filtering)
Retrieves the stories (images or videos) posted by the store. This endpoint powers the interactive story modal and the "Stories" tab.

**Endpoint:** `GET /api/v1/stores/{storeName}/stories`

### Request
- **Path Parameters:**
  - `storeName` (string): The slug or name of the store.
- **Query Parameters:**
  - `page` (integer, optional): The page number for pagination. Default `1`.
  - `limit` (integer, optional): The number of stories per page. Default `10`.
  - `type` (string, optional): Filter by story type. Enum: `video`, `image`.
  - `search` (string, optional): Search term to filter stories by title.
  - `sort_by` (string, optional): Sorting criteria. Options: `newest`, `popular`.

### Example Request
`GET /api/v1/stores/techhaven/stories?page=1&limit=20&type=video&sort_by=newest`

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "stories": [
      {
        "id": "sty-987",
        "type": "video",
        "title": "Unboxing Earbuds Pro",
        "thumbnail": "https://example.com/thumb1.jpg",
        "videoUrl": "https://example.com/vid1.mp4",
        "views": 12500,
        "createdAt": "2023-10-25T10:30:00Z"
      },
      {
        "id": "sty-988",
        "type": "image",
        "title": "New Arrivals Setup",
        "thumbnail": "https://example.com/thumb2.jpg",
        "imageUrl": "https://example.com/img1.jpg",
        "views": 8500,
        "createdAt": "2023-10-22T14:15:00Z"
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 5,
      "totalItems": 45,
      "hasNextPage": true
    }
  }
}
```

---

## 3. Get Store Products
Retrieves the list of products sold by the store.

**Endpoint:** `GET /api/v1/stores/{storeName}/products`

### Request
- **Query Parameters:**
  - `page` (integer, optional): Page number. Default `1`.
  - `limit` (integer, optional): Items per page. Default `20`.
  - `search` (string, optional): Search query for product name.
  - `category` (string, optional): Filter by product category.

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "products": [
      {
        "id": "prod-101",
        "name": "Wireless Earbuds Pro",
        "price": 1200,
        "currency": "ETB",
        "category": "Electronics",
        "images": ["https://example.com/prod1.jpg"]
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 10,
      "totalItems": 195,
      "hasNextPage": true
    }
  }
}
```

---

## 4. Get Store Recent Sales
Retrieves the recent sales for the store, shown in the "Sells" tab.

**Endpoint:** `GET /api/v1/stores/{storeName}/sales`

### Request
- **Query Parameters:**
  - `page` (integer, optional): Page number. Default `1`.
  - `limit` (integer, optional): Items per page. Default `10`.

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "sales": [
      {
        "id": "ORD-12345",
        "productName": "Wireless Earbuds Pro",
        "amount": 1200,
        "currency": "ETB",
        "status": "Delivered",
        "buyerName": "John D.",
        "purchasedAt": "2023-10-26T09:15:00Z"
      }
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 3,
      "totalItems": 25,
      "hasNextPage": true
    }
  }
}
```
