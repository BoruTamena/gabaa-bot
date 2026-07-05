# Telegram Bot Web Auth — Frontend Integration Guide

This guide explains how to integrate **Login with Telegram** on the Gabaa web app (`gabaa-web.vercel.app`) using the bot-mediated session flow.

Use this flow when the user is **outside** the Telegram Mini App (regular browser). Inside the Mini App, keep using `POST /auth/telegram` with `initData`.

---

## Overview

```text
1. User clicks "Login with Telegram" on your website
2. Frontend calls POST /auth/telegram/session
3. User opens the returned botUrl in Telegram (@gabaaBot)
4. User taps START in the bot → bot confirms login
5. Frontend polls GET /auth/telegram/session/:sessionId
6. When status is "completed", save the JWT and use it for all API calls
```

**Bot:** [@gabaaBot](https://t.me/gabaaBot)  
**Session lifetime:** 5 minutes  
**Recommended poll interval:** every 2 seconds

---

## Standard response envelope

All endpoints return:

**Success**
```json
{
  "success": true,
  "data": { },
  "error": null
}
```

**Error**
```json
{
  "success": false,
  "data": null,
  "error": {
    "error": "NOT_FOUND",
    "message": "Session not found or expired"
  }
}
```

The examples below show only the `data` payload.

---

## Endpoints

### 1. Start login session

**`POST /auth/telegram/session`**

- **Auth required:** No
- **Request body:** None

**Response `data`:**
```json
{
  "sessionId": "a1b2c3d4e5f6789012345678abcdef01",
  "botUrl": "https://t.me/gabaaBot?start=login_a1b2c3d4e5f6789012345678abcdef01",
  "expiresAt": "2026-07-05T18:05:00Z"
}
```

| Field | Description |
|-------|-------------|
| `sessionId` | Store this locally; used for polling |
| `botUrl` | Open this link so the user starts the bot with the login payload |
| `expiresAt` | ISO timestamp; stop polling after this time |

---

### 2. Poll login session

**`GET /auth/telegram/session/:sessionId`**

- **Auth required:** No

**Pending** (user has not confirmed in Telegram yet):

```json
{
  "status": "pending"
}
```

**Completed** (login successful — session is consumed and deleted):

```json
{
  "status": "completed",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "userId": 1,
  "telegramUserId": 123456789,
  "username": "johndoe",
  "role": "customer",
  "hasStore": false,
  "storeId": 0
}
```

| Field | Description |
|-------|-------------|
| `token` | JWT — use as `Authorization: Bearer <token>` |
| `userId` | Internal Gabaa user ID |
| `telegramUserId` | Telegram user ID |
| `role` | `"customer"` or `"admin"` |
| `hasStore` | `true` if merchant has a linked store |
| `storeId` | Linked store ID (merchants only) |

**Expired or invalid session** → `404` with `error.error: "NOT_FOUND"`

---

### 3. Authenticated API calls (after login)

All protected endpoints expect:

```http
Authorization: Bearer <token>
```

JWT is valid for **7 days**. There is no refresh endpoint yet — re-run the bot login flow when the token expires.

---

## Choosing the right auth method

| Context | Method |
|---------|--------|
| User inside Telegram Mini App | `POST /auth/telegram` with `initData` |
| User on web browser (gabaa-web.vercel.app) | Bot session flow (this guide) |
| Future mobile app | Same bot session flow |

### Mini App login (existing)

**`POST /auth/telegram`**

```json
{
  "initData": "<window.Telegram.WebApp.initData>"
}
```

Response `data` has the same shape as the completed poll response (`token`, `userId`, `role`, etc.).

---

## Frontend implementation

### Step 1 — API base URL

Point requests at your backend, for example:

```ts
const API_BASE = import.meta.env.VITE_API_URL ?? "https://your-api.example.com";
```

### Step 2 — Start session and open Telegram

```ts
type LoginSession = {
  sessionId: string;
  botUrl: string;
  expiresAt: string;
};

type PollResult =
  | { status: "pending" }
  | {
      status: "completed";
      token: string;
      userId: number;
      telegramUserId: number;
      username: string;
      role: string;
      hasStore: boolean;
      storeId?: number;
    };

async function startTelegramLogin(): Promise<LoginSession> {
  const res = await fetch(`${API_BASE}/auth/telegram/session`, {
    method: "POST",
  });
  const json = await res.json();
  if (!json.success) throw new Error(json.error?.message ?? "Failed to start login");
  return json.data;
}
```

Open Telegram in a new tab/window:

```ts
function openTelegramLogin(botUrl: string) {
  window.open(botUrl, "_blank", "noopener,noreferrer");
}
```

On mobile, `window.open(botUrl)` opens the Telegram app via the universal link.

### Step 3 — Poll until completed

```ts
const POLL_INTERVAL_MS = 2000;

async function pollTelegramLogin(sessionId: string, expiresAt: string): Promise<PollResult> {
  const deadline = new Date(expiresAt).getTime();

  while (Date.now() < deadline) {
    const res = await fetch(`${API_BASE}/auth/telegram/session/${sessionId}`);
    const json = await res.json();

    if (res.status === 404) {
      throw new Error("Login session expired. Please try again.");
    }
    if (!json.success) {
      throw new Error(json.error?.message ?? "Login poll failed");
    }

    if (json.data.status === "completed") {
      return json.data;
    }

    await new Promise((r) => setTimeout(r, POLL_INTERVAL_MS));
  }

  throw new Error("Login timed out. Please try again.");
}
```

### Step 4 — Full login handler

```ts
async function loginWithTelegram() {
  const session = await startTelegramLogin();
  openTelegramLogin(session.botUrl);

  const result = await pollTelegramLogin(session.sessionId, session.expiresAt);

  if (result.status !== "completed") {
    throw new Error("Unexpected login state");
  }

  localStorage.setItem("gabaa_token", result.token);
  localStorage.setItem("gabaa_user", JSON.stringify({
    userId: result.userId,
    telegramUserId: result.telegramUserId,
    username: result.username,
    role: result.role,
    hasStore: result.hasStore,
    storeId: result.storeId ?? 0,
  }));

  return result;
}
```

### Step 5 — Attach token to API requests

```ts
async function apiFetch(path: string, options: RequestInit = {}) {
  const token = localStorage.getItem("gabaa_token");
  const headers = new Headers(options.headers);

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  if (!headers.has("Content-Type") && options.body) {
    headers.set("Content-Type", "application/json");
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });
  const json = await res.json();

  if (res.status === 401) {
    localStorage.removeItem("gabaa_token");
    localStorage.removeItem("gabaa_user");
    // redirect to login page
  }

  return json;
}
```

---

## React example (login button)

```tsx
import { useState } from "react";

export function TelegramLoginButton() {
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  async function handleLogin() {
    setLoading(true);
    setMessage("Opening Telegram… Confirm login in @gabaaBot");

    try {
      const result = await loginWithTelegram();
      setMessage(`Welcome, ${result.username || "user"}!`);
      window.location.href = "/"; // or your post-login route
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <button type="button" onClick={handleLogin} disabled={loading}>
        {loading ? "Waiting for Telegram…" : "Login with Telegram"}
      </button>
      {message && <p>{message}</p>}
    </div>
  );
}
```

---

## UX recommendations

1. **Clear instructions** while polling:  
   *"We opened Telegram. Tap START on @gabaaBot, then return here."*

2. **Show a spinner** during polling; disable the login button to avoid duplicate sessions.

3. **Handle timeout** — if 5 minutes pass, show "Try again" and call `POST /auth/telegram/session` again.

4. **Desktop fallback** — if pop-ups are blocked, show the `botUrl` as a clickable link or QR code.

5. **Do not poll after `completed`** — the session is deleted; a second poll returns `404`.

6. **Merchant routing** — after login:
   - `role === "admin" && hasStore` → merchant dashboard
   - otherwise → customer storefront

---

## Error handling

| HTTP | `error.error` | Frontend action |
|------|---------------|-----------------|
| 404 | `NOT_FOUND` | Session expired or invalid — restart login |
| 500 | `INTERNAL_ERROR` | Show retry; check backend logs |
| 401 (on API calls) | `UNAUTHORIZED` | Clear token; redirect to login |

Bot-side errors (user sees these in Telegram, not in your API):
- Invalid/expired link → *"This login link is invalid or has expired…"*
- Success → *"You're logged in to Gabaa! Return to the app to continue."*

---

## Integration checklist

- [ ] `POST /auth/telegram/session` on "Login with Telegram" click
- [ ] Open `botUrl` in new tab / Telegram app
- [ ] Poll `GET /auth/telegram/session/:sessionId` every 2s until `completed` or timeout
- [ ] Persist `token` (localStorage or secure cookie strategy)
- [ ] Send `Authorization: Bearer <token>` on all protected routes
- [ ] Handle `401` by clearing auth and showing login again
- [ ] Keep Mini App path: auto-login with `initData` when `window.Telegram?.WebApp` is available
- [ ] Backend migration applied: `20260705000001_add_telegram_login_sessions`

---

## Environment variables (frontend)

```env
VITE_API_URL=https://your-gabaa-api.example.com
```

No Telegram bot token is needed on the frontend. The backend validates the user via the bot session.

---

## Related docs

- [Frontend Integration Guide](./frontend_integration_guide.md) — general API patterns and Mini App auth
- [Store Setup Frontend Integration](./store_setup_frontend_integration.md) — merchant store linking via bot
