# Frontend Integration Guide

This guide outlines the integration process for the Gabaa E-Commerce Telegram Mini App, including all updated endpoints, HTTP methods, and the new standardized JSON response structure.

## 0. Standard API Responses

All endpoints adhere to a standardized JSON response format.

**Success Response Envelope:**
```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

**Error Response Envelope:**
```json
{
  "success": false,
  "data": null,
  "error": {
    "error": "BAD_REQUEST",
    "message": "Detailed error message here"
  }
}
```
**Common Error Codes:**
- `UNAUTHORIZED`: Invalid or missing token/initData.
- `FORBIDDEN`: No permission for this resource.
- `VALIDATION_ERROR`: Field validation failed.
- `NOT_FOUND`: Resource doesn't exist.
- `INTERNAL_ERROR`: Unexpected server error.

*Note: All endpoints below show only the `data` payload or request bodies for brevity unless the full envelope is required for context. Every request (except `/auth/telegram`) must include a `Bearer` token in the `Authorization` header.*

---

## 1. Authentication & Onboarding

### Option A: Telegram Mini App Login
**Endpoint:** `POST /auth/telegram`  
**Description:** Validate Telegram `initData` and fetch JWT. Use this for automatic login within the Mini App.

**Request Body:**
```json
{
  "initData": "user=%7B%22id%22%3A123456...&auth_date=162..."
}
```

### Option B: Email & Password (Dashboard)
**Endpoint:** `POST /auth/register` (Register new merchant/user)  
**Endpoint:** `POST /auth/login` (Login existing merchant/user)

**Request Body:**
```json
{
  "email": "merchant@example.com",
  "password": "securepassword123"
}
```

### Option C: Telegram Login Widget (Dashboard)
**Description:** Use the official [Telegram Login Widget](https://core.telegram.org/widgets/login).
**Callback Endpoint:** `GET /auth/telegram/callback` (Handled by backend to issue JWT).

**Response Data (`data` field for all auth methods):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiI...",
  "userId": 1,
  "username": "customer_joe",
  "role": "merchant",  // or "customer"
  "hasStore": false,
  "storeId": 0
}
```

---

## 2. Store Setup & Silent Linking

For merchants to link their store to a Telegram Group/Channel silently:

1.  **Initialize Store:** Create the store via `POST /store/from-chat` (or similar).
2.  **Deep Link to Bot:** Redirect the merchant to the bot with a start payload:  
    `https://t.me/GabaaPlaceBot?start=link_store_{store_id}`
3.  **Bot PM Start:** When the merchant clicks "Start" in PM, the bot acknowledges but does nothing in public.
4.  **Add Bot to Group:** Merchant adds the bot to their target group/channel as an **Administrator**.
5.  **Silent Activation:** The backend detects the addition and silently links the group to the store. The merchant receives a **Private Message** confirmation: "✅ Store successfully linked!"

### Get Dashboard State
**Endpoint:** `GET /store/dashboard/:chat_id`  
**Description:** Determine if the user needs to set up a store or manage products.

**Response Data:**
```json
{
  "dashboard_type": "setup", // or "manage"
  "store": null // or Store details
}
```

### Create Store (Setup)
**Endpoint:** `POST /store/from-chat`  
**Description:** Create a new store (Admin only).

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

### Get / Update Store
- **GET** `/store/:store_id` (Fetch store details)
- **PUT** `/store/:store_id` (Update store details)

---

## 3. Categories Management

- **GET** `/categories?page=1&page_size=10` (List all global categories)
- **GET** `/store/:store_id/categories` (List specific store categories)
- **POST** `/store/:store_id/category` (Add custom store category)

---

## 4. Product Management

### Public Browsing
- **GET** `/products?page=1&page_size=10&category=Phones&query=iPhone` (Public list/search all products)
- **GET** `/products/:id` (Fetch single public product)

### Store Products (Admin / Store View)
- **GET** `/store/:store_id/products?page=1&page_size=10`
- **POST** `/store/:store_id/product` (Create new product)
- **PUT** `/store/:store_id/product/:id` (Update existing product)
- **DELETE** `/store/:store_id/product/:id` (Delete product)

**Product Request Body Example (POST/PUT):**
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

## 5. Shopping Cart Management

The user's cart is managed via Redis and persists per user.

- **POST** `/user/cart/add?product_id=1&quantity=2` (Add item to cart)
- **GET** `/user/cart` (Fetch active user cart)
- **PUT** `/user/cart/update?product_id=1&action=increment` (Increase quantity by 1)
- **PUT** `/user/cart/update?product_id=1&action=decrement` (Decrease quantity by 1. If quantity reaches 0, item is removed.)
- **DELETE** `/user/cart/remove?product_id=1` (Remove item completely)
- **DELETE** `/user/cart/clear` (Empty the cart)

---

## 6. Order Management

### Checkout (Create Order)
**Endpoint:** `POST /order/create?store_id=1`  
**Description:** Converts the user's cart items (filtered by the given store ID) into an order. Cart items for that store are then cleared.

**Response Data Example:**
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

### User Orders
- **GET** `/user/orders?page=1&page_size=10` (List all orders placed by the user)
- **GET** `/orders/:order_id` (Get specific order details)
- **PUT** `/user/orders/:order_id/cancel` (Cancel a 'pending' order - restores product stock automatically)

### Store Orders (Admin)
- **GET** `/store/:store_id/orders?page=1&page_size=20` (List store orders)
- **PUT** `/store/:store_id/orders/:order_id/status?status=completed` (Update an order status to 'processing', 'completed', 'shipped', etc.)

---

## 7. Payments & Wallet

### Verify Payment (Manual)
**Endpoint:** `POST /payment/verify`  
**Description:** Confirm offline/manual payment. This marks the order as completed and credits the store wallet.

**Request:**
```json
{
  "order_id": 101
}
```

### Get Store Wallet
**Endpoint:** `GET /store/:store_id/wallet`  
**Description:** Fetch the store's current verified wallet balance.

**Response Data:**
```json
{
  "balance": 2400.50
}
```
