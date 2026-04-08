# Frontend Integration Guide

This guide outlines the integration process for the Gabaa E-Commerce Telegram Mini App.

## 1. Authentication & Onboarding Flow

All requests (except authentication) require a `Bearer` token in the `Authorization` header.

### Scenario: Authenticate User
**Endpoint:** `POST /auth/telegram`  
**Description:** Send the `initData` provided by Telegram to get a JWT token. The response tells the frontend how to route the user.

**Request:**
```json
{
  "initData": "user=%7B%22id%22%3A123456...&auth_date=162...&hash=..."
}
```

### Response Scenarios:

#### A. Customer (Shopper)
*Role is "customer". Frontend should redirect to the marketplace or a specific store's product list.*
```json
{
  "token": "eyJhbGciOiJIUzI1NiI...",
  "userId": 1,
  "username": "customer_joe",
  "role": "customer",
  "hasStore": false
}
```

#### B. Admin with NO Store
*Role is "admin" but `hasStore` is false. Frontend must guide the user to the **Store Setup** screen.*
```json
{
  "token": "eyJhbGciOiJIUzI1NiI...",
  "userId": 2,
  "username": "seller_new",
  "role": "admin",
  "hasStore": false
}
```

#### C. Admin WITH Store
*Role is "admin", `hasStore` is true, and `storeId` is provided. Frontend should redirect to the **Seller Dashboard** for that store.*
```json
{
  "token": "eyJhbGciOiJIUzI1NiI...",
  "userId": 3,
  "username": "seller_pro",
  "role": "admin",
  "hasStore": true,
  "storeId": 10
}
```

---

## 2. Store Management

### Scenario: First-time Seller (Determine State)
**Endpoint:** `GET /store/dashboard/:chat_id`  
**Description:** Use this to decide if the seller needs to set up a store or manage products.

**Response:**
```json
{
  "dashboard_type": "setup", // or "manage"
  "store": null // or Store object
}
```

### Scenario: Create Store (Setup)
**Endpoint:** `POST /store/from-chat`

**Request:**
```json
{
  "telegram_chat_id": 123456,
  "name": "My Great Store",
  "category": "Electronics",
  "description": "The best gadgets in town",
  "phone": "+251911223344",
  "email": "contact@mystore.com",
  "location": "Addis Ababa"
}
```

### Scenario: Get/Update Store Profile
**GET** `/store/:store_id`  
**PUT** `/store/:store_id`

---

## 3. Product Management

### Scenario: List Store Products
**Endpoint:** `GET /store/:store_id/products?page=1&page_size=10`

**Response:**
```json
{
  "total": 50,
  "data": [
    {
      "id": 1,
      "name": "Smartphone",
      "price": 500,
      "stock": 10,
      "images": "url1,url2"
    }
  ]
}
```

### Scenario: Add New Product
**Endpoint:** `POST /store/:store_id/product`

**Request:**
```json
{
  "name": "Laptop",
  "description": "Gaming laptop",
  "price": 1200,
  "stock": 5,
  "images": "https://example.com/image.jpg"
}
```

---

## 4. Shopping Cart & Checkout

### Scenario: Add Item to Cart
**Endpoint:** `POST /order/cart/add?product_id=1&quantity=2`  
**Description:** Items are stored in Redis cache per user.

**Response:**
```json
{
  "message": "added to cart"
}
```

### Scenario: Checkout (Create Order)
**Endpoint:** `POST /order/create?store_id=1`  
**Description:** Converts all items in the user's cart (for that store) into a single order.

**Response:**
```json
{
  "id": 101,
  "store_id": 1,
  "user_id": 123,
  "status": "pending",
  "total_price": 2400,
  "order_items": [
    {
      "product_id": 1,
      "quantity": 2,
      "price": 1200
    }
  ]
}
```

---

## 5. Order Management (Seller Side)

### Scenario: List Store Orders
**Endpoint:** `GET /store/:store_id/orders?page=1&page_size=20`

### Scenario: Verify Payment (Manual)
**Endpoint:** `POST /payment/verify`  
**Description:** Mark an order as paid and credit the seller's wallet.

**Request:**
```json
{
  "order_id": 101
}
```

---

## Error Handling

All error responses follow this structure:
```json
{
  "error": "BAD_REQUEST",
  "message": "Invalid email format"
}
```
**Common Error Codes:**
- `UNAUTHORIZED`: Invalid or missing token/initData.
- `FORBIDDEN`: No permission for this resource.
- `VALIDATION_ERROR`: Field validation failed.
- `NOT_FOUND`: Resource doesn't exist.
