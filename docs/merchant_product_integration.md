# Merchant Product Integration Guide

This guide details the workflow for merchants (store owners) to manage their products, from uploading images to publishing them on Telegram.

---

## 1. Image Upload Workflow
Before creating a product, you must upload its images to Cloudinary to get public URLs.

**Endpoint:** `POST /upload/images`  
**Content-Type:** `multipart/form-data`

**Request:**
- `files`: Multiple files (images).

**Response:**
```json
{
  "success": true,
  "data": [
    "https://res.cloudinary.com/.../image1.jpg",
    "https://res.cloudinary.com/.../image2.jpg"
  ],
  "error": null
}
```

---

## 2. Creating a Product
By default, new products are created with the status `draft`. They are **not** pushed to Telegram until published.

**Endpoint:** `POST /store/products` (Note: Ensure your `store_id` is in the context or path depending on your routing)  
**Auth Required:** Bearer Token

**Payload:**
```json
{
  "name": "Vintage Leather Bag",
  "description": "Handcrafted genuine leather bag.",
  "price": 2500.00,
  "stock": 10,
  "category": "Fashion",
  "images": ["https://res.cloudinary.com/.../image1.jpg"],
  "is_posted": false
}
```

---

## 3. Managing Products (Listing & Filtering)
Merchants can filter their inventory by status (draft vs published), stock levels, or search by title.

**Endpoint:** `GET /products` (Filter for owner's store)  
**Query Parameters:**
- `query`: Search by title/description (e.g., `?query=leather`).
- `status`: Filter by status (`draft`, `published`, `archived`).
- `min_stock`: Filter by minimum inventory (e.g., `?min_stock=1`).
- `page` & `page_size`: Pagination.

**Response:**
Includes the `status` field for each product to show in the UI.

---

## 4. Publishing to Telegram
To push a product to the linked Telegram Group/Channel, update its status to `published`.

**Endpoint:** `PUT /store/products/:product_id`

**Payload:**
```json
{
  "status": "published"
}
```

### What happens next?
1. The backend detects the status change to `published`.
2. The Bot automatically sends a rich message to the linked group.
3. The message includes:
   - **Product Photo**
   - **Price and Description**
   - **"🛒 Order Now" Button**: A direct link to the product's detail page on the web app.

---

## 5. Summary of Product Statuses

| Status | Behavior |
| :--- | :--- |
| `draft` | Visible only to the merchant in the dashboard. |
| `published` | Pushed to Telegram and visible to all customers on the storefront. |
| `archived` | Hidden from the storefront; kept in merchant history. |

---

## 6. Integration Checklist
- [ ] Use `POST /upload/images` first to get image URLs.
- [ ] Display a "Draft" badge for new products in the merchant UI.
- [ ] Provide a "Publish to Telegram" button that triggers a `PUT` request with `status: "published"`.
- [ ] Use `min_stock=1` filter to show "Out of Stock" items separately to the merchant.
