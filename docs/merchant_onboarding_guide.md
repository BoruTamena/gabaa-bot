# Merchant Onboarding & Store Setup Guide

This guide explains how a merchant can set up their digital storefront using Gabaa Bot, from registration to linking their Telegram group or channel.

---

## 1. Authentication & Onboarding

### Option A: Telegram Mini App Login (Recommended)
**Endpoint:** `POST /auth/telegram`  
**Description:** Validate Telegram `initData` and fetch JWT. This is the primary method for merchants using the Mini App.

**Request Body:**
```json
{
  "initData": "user=%7B%22id%22%3A123456...&auth_date=162..."
}
```

### Option B: Email & Password / Login Widget (Web)
Use these for the standalone web dashboard. See the [Frontend Integration Guide](./frontend_integration_guide.md) for details.

---

### Step 2: Create Your Store

Once authenticated, you can create your store. 

#### Path A: From Telegram Mini App (Easiest)
If you are using the Mini App, the `chat_id` is automatically detected.
1. Open the Mini App in your Private Chat or a Group where you are an admin.
2. The dashboard will show a **"Setup Store"** button if no store is linked.
3. Fill in the details. The `telegram_chat_id` will be sent automatically.

#### Path B: From Web Dashboard
1. Fill in your store details.
2. You can leave the `telegram_chat_id` as `0` if you want to link it later via the Silent Linking flow (Step 3).

**API Action:** `POST /store/from-chat`

---

## Step 3: Initialize the Bot (Silent Link)

Before you can link your store to a group, the bot needs to know who you are.

1.  On your dashboard, click the **"Connect to Bot"** button.
2.  This will open a private chat with the bot in Telegram via a link like:  
    `https://t.me/GabaaPlaceBot?start=link_store_{your_store_id}`
3.  Click **"START"**.
4.  **Result:** The bot will silently register your Telegram ID and associate it with your store account. No message is posted in any group yet.

---

## Step 4: Link Your Group or Channel

Now you can link your store to your target Telegram community.

1.  Go to your Telegram Group or Channel.
2.  Add **Gabaa Bot** (`@GabaaPlaceBot`) as a member.
3.  **Crucial Step:** Promote the bot to an **Administrator**.
4.  The bot must have "Post Messages" or "Manage Chat" permissions.

---

## Step 5: Silent Confirmation

As soon as you add the bot as an admin, the following happens automatically:

1.  The bot detects it has been added by a registered merchant.
2.  It silently links the Group ID to your Store ID in our database.
3.  The bot **will NOT** post anything in the group during this process.
4.  You (the merchant) will receive a **Private Message (PM)** from the bot:  
    > "✅ Store 'Abebe's Electronics' successfully linked to group 'Abebe Fans Group'!"

---

## Step 6: Verify and Start Selling

1.  Open your Group/Channel.
2.  Click on the **"Menu"** button or the **"Mini App"** link in the chat.
3.  The Gabaa Mini App will open, automatically detecting your store context and showing your products to all members of the group.

---

## Troubleshooting

- **Bot not linking?** Ensure you clicked "Start" in the private chat (Step 3) before adding the bot to the group.
- **Wrong store linked?** If you have multiple stores, make sure you use the specific `/start` link generated for the store you wish to link.
- **Mini App not loading?** Verify the bot is still an administrator in the group.
