# UI Integration Guide - New Endpoints

This guide provides the necessary information for the frontend to integrate with the newly added public product filtering, user orders, and cart status endpoints.

---

## 1. Public Product Listing

Get a paginated list of products with optional search and category filters.

### **Endpoint**
`GET /products`

### **Authorization**
None (Public)

### **Query Parameters**
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `category` | string | No | Filter products by category (e.g., `Electronics`) |
| `query` | string | No | Search term for product names and descriptions |
| `page` | int | No | Page number (default: `1`) |
| `page_size` | int | No | Number of items per page (default: `10`) |

### **Example Request**
`GET /products?category=Home&query=lamp&page=1&page_size=20`

### **Success Response**
```json
{
  "total": 45,
  "data": [
    {
      "id": 1,
      "store_id": 10,
      "name": "Modern Desk Lamp",
      "description": "Adjustable LED lamp for offices.",
      "price": 25.99,
      "stock": 100,
      "category": "Home",
      "images": "[\"image1.jpg\", \"image2.jpg\"]"
    }
  ]
}
```

---

## 2. User Cart Status

Retrieve the full status of the authenticated user's cart, including product details and total price.

### **Endpoint**
`GET /user/cart`

### **Authorization**
JWT Token required in `Authorization` header.

### **Example Request**
`GET /user/cart`

### **Success Response**
```json
{
  "items": [
    {
      "product": {
        "id": 1,
        "store_id": 10,
        "name": "Modern Desk Lamp",
        "category": "Home",
        "price": 25.99,
        "stock": 100,
        "images": "[\"image1.jpg\"]"
      },
      "quantity": 2
    }
  ],
  "total_price": 51.98
}
```

---

## 3. User Orders History

List all orders placed by the authenticated user with pagination.

### **Endpoint**
`GET /user/orders`

### **Authorization**
JWT Token required in `Authorization` header.

### **Query Parameters**
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `page` | int | No | Page number (default: `1`) |
| `page_size` | int | No | Number of items per page (default: `10`) |

### **Example Request**
`GET /user/orders?page=1&page_size=5`

### **Success Response**
```json
{
  "total": 12,
  "data": [
    {
      "id": 501,
      "store_id": 10,
      "user_id": 99,
      "status": "pending",
      "total_price": 75.50,
      "created_at": "2024-03-20T10:00:00Z",
      "order_items": [
        {
          "id": 1001,
          "order_id": 501,
          "product_id": 1,
          "quantity": 3,
          "price": 25.16
        }
      ]
    }
  ]
}
```

---

> [!TIP]
> **Performance Tip**: When searching, avoid sending a request for every keystroke. Use a debounce function (e.g., 300ms) to reduce API calls.

> [!IMPORTANT]
> **Authentication**: Ensure the `Authorization: Bearer <token>` header is sent for all `/user/*` and `POST /store/*` endpoints. Failure to do so will result in a `401 Unauthorized` response.

---

## 4. Product Categories (Catalog)

### **Get All Categories (Public)**
Used for global product filtering.
- **Endpoint**: `GET /categories`
- **Auth**: None
- **Query Params**: `page`, `page_size`

### **Get Store Catalog**
Used for a specific store's product management.
- **Endpoint**: `GET /store/:store_id/categories`
- **Auth**: JWT required
- **Response**: List of categories available to that store (Store Specific + General).

### **Create Store Category**
Add a new category to a store's local catalog.
- **Endpoint**: `POST /store/:store_id/category`
- **Auth**: JWT required (Admin)
- **Body**: `{"name": "Custom Category Name"}`
