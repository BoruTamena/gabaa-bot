# LakiPay Payment Integration

This guide covers checkout payment initiation and order/payment status for the frontend.

## Overview

Checkout creates an order and initiates a LakiPay direct payment in one step. Payment status updates arrive via server-side webhooks — the frontend should poll order status after checkout.

## Wallet Balance Semantics

| Field | Meaning |
|-------|---------|
| `pending_balance` | Paid orders awaiting delivery (escrow) |
| `available_balance` | Released after delivery; withdrawable |
| `locked_balance` | Reserved during active withdrawal requests |
| `total_earned` | Lifetime released earnings |
| `total_withdrawn` | Lifetime withdrawn amount |

## Checkout

**Endpoint:** `POST /order/create` (JWT required)

**Request:**
```json
{
  "store_id": 1,
  "address_id": 5,
  "medium": "TELEBIRR",
  "phone_number": "251911234567"
}
```

**Supported mediums:** `MPESA`, `TELEBIRR`, `CBE`, `ETHSWITCH`

`phone_number` is optional — defaults to the shipping address phone. Must be Ethiopian format (`251...`).

**Response:**
```json
{
  "data": {
    "order": {
      "id": 42,
      "store_id": 1,
      "status": "pending",
      "total_price": 149.99,
      "order_items": [...]
    },
    "payment": {
      "id": 10,
      "order_id": 42,
      "status": "pending",
      "method": "lakipay",
      "reference": "ORDER-42-10",
      "transaction_id": "TXN-123456789",
      "amount": 149.99,
      "currency": "ETB",
      "medium": "TELEBIRR",
      "gateway_status": "PENDING"
    }
  }
}
```

**Store requirement:** Store must have `verification_status: verified` or checkout returns an error.

## Payment Lifecycle

| Stage | Payment status | Order status |
|-------|---------------|--------------|
| Checkout started | `initiated` | `pending` |
| LakiPay accepted | `pending` | `pending` |
| Webhook SUCCESS | `success` | `paid` |
| Webhook FAILED | `failed` | `cancelled` |
| Merchant delivers | `success` | `delivered` |

After payment success, funds go to the store's `pending_balance`. When the merchant marks the order `delivered`, funds move to `available_balance`.

## Polling Order Status

After checkout, poll:

**Endpoint:** `GET /orders/:order_id`

Stop polling when `status` is `paid`, `cancelled`, or after a timeout (~5 minutes).

## Wallet Summary

**Merchant endpoint:** `GET /my-store/wallet` (JWT with `store_id` in token)

**Legacy endpoint:** `GET /store/:store_id/wallet` (JWT required)

**Response:**
```json
{
  "data": {
    "id": 1,
    "store_id": 1,
    "currency": "ETB",
    "pending_balance": 500.00,
    "available_balance": 1200.00,
    "locked_balance": 0,
    "total_earned": 5000.00,
    "total_withdrawn": 0
  }
}
```

## Merchant Withdrawal

**Endpoint:** `POST /my-store/wallet/withdraw` (JWT required, store from token)

**Request:**
```json
{
  "amount": 500.00,
  "currency": "ETB",
  "phone_number": "251911234567",
  "medium": "TELEBIRR"
}
```

**Supported currencies:** `ETB`, `USD`

**Supported mediums:** `MPESA`, `TELEBIRR`, `CBE`, `ETHSWITCH`

**Response:**
```json
{
  "data": {
    "id": 3,
    "store_id": 1,
    "amount": 500.00,
    "currency": "ETB",
    "phone_number": "251911234567",
    "medium": "TELEBIRR",
    "reference": "WITHDRAW-1-3",
    "transaction_id": "TXN-987654321",
    "status": "pending",
    "gateway_status": "PENDING",
    "created_at": "2026-07-12T12:00:00Z"
  }
}
```

### Withdrawal Lifecycle

| Stage | Withdrawal status | Wallet effect |
|-------|------------------|---------------|
| Request submitted | `initiated` | `available_balance -= amount`, `locked_balance += amount` |
| LakiPay accepted | `pending` | funds stay locked |
| Webhook SUCCESS | `success` | `locked_balance -= amount`, `total_withdrawn += amount` |
| Webhook FAILED | `failed` | `locked_balance -= amount`, `available_balance += amount` |
| Webhook CANCELLED | `cancelled` | `locked_balance -= amount`, `available_balance += amount` |
| Webhook PENDING | `pending` | funds stay locked |

**List withdrawals:** `GET /my-store/wallet/withdrawals?page=1&page_size=20`

### Polling withdrawal status

After `POST /my-store/wallet/withdraw`, save the returned `id` and poll:

**Endpoint:** `GET /my-store/wallet/withdrawals/:withdrawal_id` (JWT required)

**Example:** `GET /my-store/wallet/withdrawals/3`

Poll every **3–5 seconds** until `status` is terminal: `success`, `failed`, or `cancelled`.

```json
{
  "data": {
    "id": 3,
    "store_id": 1,
    "amount": 500.00,
    "currency": "ETB",
    "phone_number": "251911234567",
    "medium": "TELEBIRR",
    "reference": "WITHDRAW-1-3",
    "transaction_id": "01d3a488-9b47-4a59-8d33-3ac78637888b",
    "status": "pending",
    "gateway_status": "PENDING",
    "created_at": "2026-07-12T12:00:00Z"
  }
}
```

Stop polling when `status` is `success`, `failed`, or `cancelled`. Optionally refresh `GET /my-store/wallet` after a terminal status to update balances.

Store must be verified (`verification_status: verified`) to withdraw.

## Webhook (Server Only)

LakiPay sends deposit and withdrawal updates to:

`POST /api/v1/webhook/lakipay`

The server routes by `event` field:
- `DEPOSIT` — order payment (checkout)
- `WITHDRAWAL` — merchant payout

Configure `LAKIPAY_CALLBACK_URL` to this endpoint. RSA signature verification is required for all events.

## Environment Variables

```env
LAKIPAY_SECRET_KEY=your_public_api_key   # LAKIPUB_...
LAKIPAY_PUB_KEY=your_secret_api_key      # LAKISEC_...
LAKIPAY_BASE_URL=https://api.lakipay.co
LAKIPAY_CALLBACK_URL=https://your-api-host/api/v1/webhook/lakipay
LAKIPAY_PUBLIC_KEY=-----BEGIN PUBLIC KEY-----...  # RSA PEM for webhook verification
```

## Frontend Checklist

- [ ] Add payment medium selector at checkout (MPESA, TELEBIRR, CBE, ETHSWITCH)
- [ ] Send `medium` and optional `phone_number` with checkout
- [ ] Show payment pending UI after checkout; poll order until `paid` or `cancelled`
- [ ] Display wallet summary with pending/available/locked balances
- [ ] Add merchant withdrawal form (amount, currency, phone, medium)
- [ ] Poll withdrawal list or wallet after withdraw until terminal status
- [ ] Poll `GET /my-store/wallet/withdrawals/:withdrawal_id` after withdraw until `success`, `failed`, or `cancelled`
- [ ] Handle store-not-verified error gracefully

## Removed Endpoint

`POST /payment/verify` (manual verification) has been removed. Payments are confirmed automatically via LakiPay webhooks.
