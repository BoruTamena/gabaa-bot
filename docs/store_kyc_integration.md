# Store KYC Verification — Frontend Integration Guide

This guide covers merchant store verification (KYC) and platform admin review. Payment gating is **not active yet** — a `RequireStoreVerified` middleware exists for future payment gateway integration.

---

## Verification statuses

| Status | Meaning |
|--------|---------|
| `unverified` | Default — store has not submitted KYC |
| `pending_review` | KYC submitted, awaiting platform admin |
| `verified` | Admin approved — eligible for payments (when gated) |
| `rejected` | Admin rejected — merchant can resubmit |

This is separate from Telegram launch status (`pending` / `launched`).

---

## Standard response envelope

```json
{
  "success": true,
  "data": { },
  "error": null
}
```

---

## Merchant flow

### Step 1 — Upload documents (JWT required)

**`POST /upload/documents`**

- Multipart field: `files` (one or more)
- Allowed types: PDF, PNG, JPEG, WEBP
- Max size: 10 MB per file
- Returns array of Cloudinary URLs

```bash
curl -X POST "$API/upload/documents" \
  -H "Authorization: Bearer $TOKEN" \
  -F "files=@tin_certificate.pdf" \
  -F "files=@business_license.png"
```

Response `data`:
```json
[
  "https://res.cloudinary.com/.../tin_certificate.pdf",
  "https://res.cloudinary.com/.../business_license.png"
]
```

### Step 2 — Submit KYC (JWT, merchant `role=admin`)

**`POST /store/verification`**

Requires JWT with `store_id` claim set.

Request body:
```json
{
  "tinNumber": "1234567890",
  "businessRegistrationNumber": "BR-2024-001",
  "tinCertificateUrl": "https://res.cloudinary.com/.../tin_certificate.pdf",
  "businessLicenseUrl": "https://res.cloudinary.com/.../business_license.png"
}
```

Response `data`:
```json
{
  "storeId": 1,
  "storeName": "My Shop",
  "verificationStatus": "pending_review",
  "tinNumber": "1234567890",
  "businessRegistrationNumber": "BR-2024-001",
  "tinCertificateUrl": "https://...",
  "businessLicenseUrl": "https://...",
  "submittedAt": "2026-07-05T20:00:00Z"
}
```

**Submit rules:**
- Allowed when status is `unverified` or `rejected`
- Blocked when `pending_review` or `verified`

### Step 3 — Check KYC status

**`GET /store/verification`**

Response `data` — same shape as submit response. If never submitted, returns only `storeId`, `storeName`, and `verificationStatus: "unverified"`.

---

## Platform admin flow

Bootstrap a platform admin in the database:

```sql
UPDATE users SET role = 'platform_admin' WHERE id = <ops_user_id>;
```

That user must log in and use a JWT with `role: "platform_admin"`.

### List pending submissions

**`GET /admin/store-verifications?status=pending_review`**

Returns array of `StoreKYCResponse` objects including document URLs.

### Approve

**`POST /admin/store-verifications/:store_id/approve`**

Sets store to `verified`.

### Reject

**`POST /admin/store-verifications/:store_id/reject`**

Request body (optional):
```json
{
  "reviewNote": "TIN certificate is expired. Please upload a current document."
}
```

Sets store to `rejected`. Merchant can resubmit.

---

## Public store response

`GET /store/:store_id` now includes:

```json
{
  "verificationStatus": "verified"
}
```

Document URLs are **not** exposed on the public store endpoint.

---

## Future payment gating

`RequireStoreVerified()` middleware is implemented but **not wired** to any route yet.

When the payment gateway is integrated, attach it to payment routes on the backend:

```go
paymentGroup := api.Group("/")
paymentGroup.Use(authMiddleware.RequireStoreVerified())
RegisterPaymentRoutes(paymentGroup, paymentHandler)
```

Until then, checkout and order status updates work for all stores regardless of verification status.

---

## Frontend checklist

- [ ] Show verification badge from `store.verificationStatus`
- [ ] KYC form: TIN, business registration number, two file uploads
- [ ] Upload docs via `POST /upload/documents` before submitting KYC
- [ ] Submit via `POST /store/verification`
- [ ] Poll or refresh `GET /store/verification` for status updates
- [ ] Handle `rejected` state with `reviewNote` and allow resubmit
- [ ] Run migration `20260706000001_add_store_kyc` on the database

---

## Related docs

- [Frontend Integration Guide](./frontend_integration_guide.md)
- [Store Setup Frontend Integration](./store_setup_frontend_integration.md)
