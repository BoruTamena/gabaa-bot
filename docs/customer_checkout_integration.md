# Customer Address & Checkout Integration Guide

This guide covers how a customer manages their shipping addresses and completes the checkout process. All these endpoints require a valid customer JWT token.

---

## 1. Address Management Endpoints

A customer can manage multiple addresses (e.g., Home, Work).

### 1.1 Add New Address
**Endpoint:** `POST /user/addresses`

**Request Body:**
```json
{
  "label": "work",
  "recipient_name": "John Doe",
  "phone": "+251911000000",
  "street": "Bole Road, XYZ Building, 3rd Floor",
  "city": "Addis Ababa",
  "region": "Addis Ababa",
  "country": "Ethiopia",
  "is_default": true
}
```
*Note: If this is the user's very first address, it will automatically become their default address, regardless of the `is_default` flag. If `is_default` is `true`, all other addresses will be unset as default.*

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 14,
    "user_id": 7,
    "label": "work",
    "recipient_name": "John Doe",
    "phone": "+251911000000",
    "street": "Bole Road, XYZ Building, 3rd Floor",
    "city": "Addis Ababa",
    "region": "Addis Ababa",
    "country": "Ethiopia",
    "is_default": true,
    "created_at": "2026-05-02T22:00:00Z"
  },
  "error": null
}
```

### 1.2 List Saved Addresses
**Endpoint:** `GET /user/addresses`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 14,
      "label": "work",
      "recipient_name": "John Doe",
      "phone": "+251911000000",
      "street": "Bole Road...",
      "city": "Addis Ababa",
      "is_default": true
    },
    {
      "id": 15,
      "label": "home",
      "recipient_name": "John Doe",
      "phone": "+251922111111",
      "street": "CMC, Block 4",
      "city": "Addis Ababa",
      "is_default": false
    }
  ],
  "error": null
}
```

### 1.3 Update Address
**Endpoint:** `PUT /user/addresses/:id`

**Request Body:**
```json
{
  "label": "office",
  "phone": "+251999999999"
}
```
*(Only provide the fields you want to change)*

### 1.4 Delete Address
**Endpoint:** `DELETE /user/addresses/:id`

**Response:**
```json
{
  "success": true,
  "data": { "message": "address deleted successfully" },
  "error": null
}
```

### 1.5 Set as Default Address
Quickly swap the active default address without sending the full payload.

**Endpoint:** `PUT /user/addresses/:id/default`

**Response:**
```json
{
  "success": true,
  "data": { "message": "default address updated" },
  "error": null
}
```

---

## 2. Checkout Process

The checkout flow requires the customer to have items in their cart for a specific store, and they must select a delivery address.

**Endpoint:** `POST /order/create`

> **Note on Cart Structure:** Carts are global but orders are store-specific. You must pass the `store_id` you are currently checking out from.

**Request Body:**
```json
{
  "store_id": 1,
  "address_id": 14
}
```

**Validation Checks performed by the Backend:**
1. Verifies the user has a valid cart with items.
2. Validates that the provided `address_id` actually belongs to the authenticated user attempting the checkout.

**Success Response:**
Returns the fully generated order with the embedded shipping address.
```json
{
  "success": true,
  "data": {
    "id": 102,
    "store_id": 1,
    "user_id": 7,
    "status": "pending",
    "total_price": 4500.00,
    "created_at": "2026-05-02T22:05:00Z",
    "shipping_address": {
      "id": 14,
      "label": "work",
      "recipient_name": "John Doe",
      "phone": "+251911000000",
      "street": "Bole Road, XYZ Building, 3rd Floor",
      "city": "Addis Ababa",
      "country": "Ethiopia"
    },
    "order_items": [
      {
        "id": 201,
        "order_id": 102,
        "product_id": 5,
        "quantity": 1,
        "price": 4500.00
      }
    ]
  },
  "error": null
}
```

---

## 3. Recommended Frontend UI Flow

1. **Cart Screen:** Display cart items. When the user taps "Checkout", fetch their saved addresses via `GET /user/addresses`.
2. **Address Selection:** 
   - If they have no addresses, show an "Add New Address" form.
   - If they have addresses, pre-select the one where `is_default === true`.
   - Allow them to swap the selected address or create a new one.
3. **Confirm & Pay:** When they confirm, fire the `POST /order/create` request passing the `store_id` and the selected `address_id`.
4. **Success:** Clear the local cart state for that store and navigate to an Order Success screen.
