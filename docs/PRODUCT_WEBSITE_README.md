# Gabaa Place — Product Profile Website Design & Development Guide

> **Purpose of this document**  
> This README is the single source of truth for designing and building the **public marketing / profile website** for **Gabaa Place** — a Telegram-native e-commerce Mini App that lets sellers launch a storefront inside their existing Telegram groups and channels with zero coding and zero extra configuration.

Use this document when writing copy, planning page layouts, choosing visuals, and scoping the website build. It reflects what is **implemented today** in the Gabaa backend (`gabaa-bot`) and clearly separates **planned** capabilities so the website stays honest while still selling the vision.

---

## Table of Contents

1. [Product Overview](#1-product-overview)
2. [Brand & Positioning](#2-brand--positioning)
3. [Core Value Propositions](#3-core-value-propositions)
4. [Who It Is For](#4-who-it-is-for)
5. [Complete Feature Catalog](#5-complete-feature-catalog)
6. [How It Works — Seller Journey](#6-how-it-works--seller-journey)
7. [How It Works — Buyer Journey](#7-how-it-works--buyer-journey)
8. [Website Information Architecture](#8-website-information-architecture)
9. [Page-by-Page Design Specifications](#9-page-by-page-design-specifications)
10. [Visual & UX Design Guidelines](#10-visual--ux-design-guidelines)
11. [Messaging & Copy Bank](#11-messaging--copy-bank)
12. [Calls to Action & Deep Links](#12-calls-to-action--deep-links)
13. [FAQ Content](#13-faq-content)
14. [SEO & Metadata](#14-seo--metadata)
15. [Technical Context for Developers](#15-technical-context-for-developers)
16. [Roadmap — Planned Features](#16-roadmap--planned-features)
17. [Appendix — API & Integration References](#17-appendix--api--integration-references)

---

## 1. Product Overview

### What is Gabaa Place?

**Gabaa Place** is a free Telegram e-commerce platform. Sellers connect **@GabaaPlaceBot** to a Telegram group or channel they already run, set up a digital storefront through the Mini App, and start selling to the community they have spent years building — without leaving Telegram, without building a separate website, and without writing a single line of code.

### One-line pitch

> Turn your Telegram group into a store in minutes — free setup, real orders, real inventory, right where your customers already are.

### Product components

| Component | Role | URL / Handle |
|-----------|------|--------------|
| **Gabaa Place** | Product brand | Marketing site (this document) |
| **@GabaaPlaceBot** | Telegram bot — store linking, product posts, onboarding | `https://t.me/GabaaPlaceBot` |
| **Gabaa Mini App** | Web storefront & merchant dashboard inside Telegram | `https://gabaa-web.vercel.app` |
| **Gabaa API** | Go backend — auth, stores, products, orders, stories, wallet | Deployed separately (see backend repo) |

### Problem we solve

Traditional e-commerce forces sellers to:

- Build and maintain a separate website or app
- Pay for hosting, plugins, and payment integrations
- Rebuild audience trust on a new platform
- Ask customers to leave the community channel they already use daily

Gabaa Place inverts this: **the store lives inside Telegram**, attached to the group or channel the seller already owns.

---

## 2. Brand & Positioning

### Brand name usage

| Context | Use |
|---------|-----|
| Marketing website, product posts footer | **Gabaa Place** |
| Telegram bot handle | **@GabaaPlaceBot** |
| Technical / repo context | **gabaa-bot**, **gabaa-web** |

### Tagline options (pick one for hero)

1. *Sell where you already connect.*
2. *Your Telegram community. Your store. Zero setup.*
3. *E-commerce built for Telegram sellers.*
4. *From group chat to checkout — in minutes.*

### Tone of voice

- **Clear** — no jargon; sellers are merchants, not developers
- **Confident** — emphasize how fast and free setup is
- **Community-first** — respect that sellers built their audience over years
- **Local** — Ethiopia-first (ETB currency, +251 phone examples, Addis Ababa in copy)

### Positioning statement

For **Telegram group and channel owners** who want to sell products online, **Gabaa Place** is a **free, no-code storefront** that runs inside Telegram. Unlike standalone shop builders or marketplace apps, Gabaa lets sellers **keep their audience in the same chat**, publish products directly to the channel, and manage orders without technical skills.

---

## 3. Core Value Propositions

These five pillars map directly to product capabilities. Each should have a dedicated section on the homepage (feature grid or scroll narrative).

### 1. Free, zero-configuration setup

**Headline:** *Launch your store for free — no coding, no plugins, no hosting.*

**What we mean:**

- No monthly platform fee for basic storefront setup (position as free tier)
- Store creation through the Mini App — name, category, description, logo, contact info
- **Silent bot linking**: add @GabaaPlaceBot as admin in your group; the store connects automatically
- No DNS, no SSL certificates, no theme files — Telegram *is* the distribution layer

**Implemented today:**

- `POST /store/from-chat` — create store from Mini App
- Dashboard routing: `setup` → `manage` → `storefront` via `GET /store/dashboard/:chat_id`
- Store statuses: `pending` (created) → `launched` (bot linked to group)
- Deep link flow: `https://t.me/GabaaPlaceBot?start=link_store_{store_id}`

**Website copy angle:**

> Open the Mini App, name your store, tap Connect to Bot, add the bot as admin in your group. Done. Your storefront is live for every member.

---

### 2. Track daily sales and earnings

**Headline:** *See every order. Know what you earned.*

**What we mean:**

- Merchants see all orders with status, customer info, shipping address, and line items
- Filter orders by status: pending, shipped, delivered, cancelled
- Store **wallet balance** reflects completed sales revenue
- Order timestamps enable daily/weekly sales views in the Mini App UI

**Implemented today:**

- `GET /my-store/orders` — paginated order list with filters
- `PUT /my-store/orders/:id` — update status to `shipped` or `delivered`
- `GET /store/:store_id/wallet` — current wallet balance (ETB)
- Wallet credited when order is marked `delivered` or payment verified
- Background **store view** tracking (analytics foundation — see Roadmap)

**Website copy angle:**

> Every sale lands in your order dashboard. Mark orders shipped and delivered, and your earnings update in your store wallet — so you always know how today went.

**UI note for website mockups:** Show a dashboard card with "Today's orders: 12", "Revenue: 4,500 ETB", "Pending: 3" — even if the dedicated analytics API is on the roadmap, orders + wallet already support this in the Mini App.

---

### 3. Manage stock in real time

**Headline:** *Inventory that stays in sync with every sale.*

**What we mean:**

- Set stock quantity when creating or editing a product
- Stock automatically decrements at checkout
- Stock restored if a pending order is cancelled
- Filter products by stock level (`min_stock`, `max_stock`) to spot low inventory
- Product lifecycle: `draft` → `published` → `archived`

**Implemented today:**

- Product CRUD at `/my-store/products` with `stock` field
- Checkout validates stock and decrements atomically
- Order cancellation restores stock
- Merchant filters: `?status=published&min_stock=1` for in-stock items; `?max_stock=0` for out-of-stock

**Website copy angle:**

> Add your products once. When a customer checks out, stock updates automatically. Get low-stock visibility before you run out.

---

### 4. Sell inside the community you already built

**Headline:** *Your customers never have to leave Telegram.*

**What we mean:**

- Store is bound to one Telegram group or channel via `telegram_chat_id`
- Buyers browse and checkout inside the **Telegram Mini App**
- Publishing a product posts a rich message to the linked group/channel with an **Order Now** button
- Deep links open product detail: `?startapp=product_{id}`
- Supports both **groups** and **channels**

**Implemented today:**

- Telegram `initData` authentication — seamless login inside Mini App
- Product publish pushes HTML product card + images to Telegram
- Mini App URL buttons for channel compatibility
- Dashboard type `storefront` for buyer catalog view in group context
- Customer cart, addresses, favorites, and checkout — all inside Mini App

**Website copy angle:**

> You spent years growing your Telegram community. Gabaa lets you sell to them right there — product posts in the channel, checkout in the Mini App, no new app to download.

**Visual suggestion:** Split-screen mockup — left: Telegram channel with product post; right: Mini App checkout screen.

---

### 5. Market with product stories, channel posts, and promotions

**Headline:** *Promote products like you promote content.*

**What we mean:**

- **Product Stories** — Instagram-style timed ads linking to a product (image or video, caption, schedule)
- **Channel posts** — publishing a product automatically creates a marketing post in the group
- **Discounts & ads** — part of the product vision (see Roadmap; not yet in API)

**Implemented today:**

| Marketing tool | Status | Details |
|----------------|--------|---------|
| Product Stories | ✅ Live | Create, schedule (`starts_at` / `ends_at`), track views |
| Telegram product posts | ✅ Live | Auto-post on `status: published` |
| Story view counter | ✅ Live | Incremented when buyers view a story |
| Discounts / coupon codes | 🔜 Planned | Not yet in backend |
| Product boosting (`is_boosted`) | 🔜 Planned | DB field exists; no logic yet |

**Story features (for website feature page):**

- Attach story to a specific product
- Upload image or video media
- Set active date range
- Public story feed: `GET /stories`
- Merchant management: `GET/POST/PUT/DELETE /my-store/stories`

**Website copy angle:**

> Run timed product stories your customers swipe through. Publish a product and it goes straight to your channel with a buy button. Discount campaigns are coming soon.

---

## 4. Who It Is For

### Primary persona — The Telegram Merchant

| Attribute | Detail |
|-----------|--------|
| **Who** | Group/channel admin selling fashion, electronics, beauty, food, or handmade goods |
| **Location** | Ethiopia (primary); expandable to other Telegram-heavy markets |
| **Tech skill** | Uses Telegram daily; may not have a website |
| **Pain** | Taking orders via DMs is chaotic; wants professionalism without cost |
| **Goal** | Organized catalog, order tracking, payments, less manual work |

**Website section:** "Built for Telegram sellers" with icons for Fashion, Electronics, Beauty, Home & Garden, Toys (matches seeded categories).

### Secondary persona — The Community Buyer

| Attribute | Detail |
|-----------|--------|
| **Who** | Member of a seller's Telegram group |
| **Behavior** | Already trusts the seller; shops on phone |
| **Needs** | Easy browse, cart, saved addresses, order history |
| **Gabaa provides** | Mini App storefront, cart, favorites, shipping addresses, order tracking |

**Website section:** Brief "For shoppers" strip — "Browse, cart, and checkout without leaving Telegram."

### Tertiary persona — The Platform Partner / Investor

- Cares about Telegram Mini App ecosystem growth
- Wants metrics: zero-setup onboarding, order volume, Ethiopian market fit
- Website should include a lightweight "Why Telegram commerce?" or "About" section

---

## 5. Complete Feature Catalog

Everything below is grouped for website feature pages. Status labels: **Live**, **Partial**, **Planned**.

### 5.1 Store & Onboarding

| Feature | Status | Description |
|---------|--------|-------------|
| Telegram Mini App login | **Live** | HMAC-validated `initData` → JWT session |
| Create store from chat context | **Live** | Auto-detects `telegram_chat_id` when opened in group |
| Store profile (name, logo, cover, description) | **Live** | Full CRUD on store |
| Category selection | **Live** | Global categories + custom store categories |
| Contact info (phone, email, location) | **Live** | Displayed on storefront |
| Silent bot linking | **Live** | Bot admin promotion → store `launched` |
| Connect-to-bot deep link | **Live** | `link_store_{id}` start parameter |
| PM confirmation on link | **Live** | Bot messages merchant privately |
| Store status indicator | **Live** | `pending` / `launched` |
| Email/password registration | **Planned** | Schema exists; no API routes |

### 5.2 Product Catalog & Inventory

| Feature | Status | Description |
|---------|--------|-------------|
| Create / edit / delete products | **Live** | Merchant product CRUD |
| Multi-image products | **Live** | Cloudinary-hosted image arrays |
| Draft / published / archived status | **Live** | Controlled visibility |
| Publish to Telegram channel | **Live** | Rich post with Order Now button |
| Stock management | **Live** | Set, filter, auto-decrement on sale |
| Product search & filters | **Live** | By query, category, status, stock |
| Public product catalog | **Live** | `GET /products` with pagination |
| Product detail page | **Live** | Single product with images, price, description |
| Image upload | **Live** | `POST /upload/images` → Cloudinary |
| Product favorites (wishlist) | **Live** | Customers save products |
| Product boosting | **Planned** | `is_boosted` column only |

### 5.3 Orders & Fulfillment

| Feature | Status | Description |
|---------|--------|-------------|
| Shopping cart | **Live** | Add, update, remove, clear |
| Shipping addresses | **Live** | CRUD + default address |
| Checkout | **Live** | Creates order, decrements stock, clears cart |
| Customer order history | **Live** | List, detail, cancel pending |
| Merchant order dashboard | **Live** | Filter by status, search by order ID |
| Order lifecycle | **Live** | pending → shipped → delivered / cancelled |
| Customer info on orders | **Live** | Username + shipping address on merchant view |
| Line items with product snapshot | **Live** | Quantity, price, product name/images |

### 5.4 Payments & Wallet

| Feature | Status | Description |
|---------|--------|-------------|
| Store wallet | **Live** | Balance per store |
| Wallet credit on delivery | **Live** | Merchant marks delivered → funds credited |
| Manual payment verification | **Live** | `POST /payment/verify` |
| ArifPay online payments | **Planned** | Stub initialized; not wired |

### 5.5 Marketing & Growth

| Feature | Status | Description |
|---------|--------|-------------|
| Product Stories (timed ads) | **Live** | Image/video, caption, schedule, views |
| Telegram channel product posts | **Live** | On publish |
| Store view tracking | **Partial** | Tracked in DB; no public analytics API yet |
| Story view analytics | **Live** | Per-story view count |
| Discount codes / promotions | **Planned** | — |
| Sponsored / boosted listings | **Planned** | — |

### 5.6 Telegram Integration

| Feature | Status | Description |
|---------|--------|-------------|
| @GabaaPlaceBot | **Live** | Onboarding, linking, product posts |
| Mini App deep links | **Live** | `?startapp=product_{id}` |
| Group & channel support | **Live** | URL buttons for channels |
| Webhook & long-polling modes | **Live** | Production-flexible bot delivery |
| Inline cart buttons | **Partial** | Code exists; publish flow uses URL buttons |

### 5.7 Customer Experience

| Feature | Status | Description |
|---------|--------|-------------|
| Browse all products | **Live** | Cross-store catalog |
| Browse by category | **Live** | Category filter |
| Product stories feed | **Live** | Swipeable story discovery |
| Favorites / wishlist | **Live** | Save products for later |
| Cart persistence | **Live** | PostgreSQL-backed cart |

---

## 6. How It Works — Seller Journey

Use this as the basis for a **"How it works"** section (4–6 steps with illustrations).

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│  Open Mini  │───▶│ Create Store │───▶│ Connect Bot │───▶│ Add Products │───▶│  Publish &  │
│  App in TG  │    │  (2 minutes) │    │  to Group   │    │  & Set Stock │    │    Sell     │
└─────────────┘    └──────────────┘    └─────────────┘    └──────────────┘    └─────────────┘
```

### Step 1 — Open & authenticate

Merchant opens the Gabaa Mini App from Telegram. Authentication is automatic via Telegram identity — no separate signup form.

### Step 2 — Create your store

If no store exists for this chat, dashboard shows **Setup**. Merchant enters store name, category, description, and contact details.

### Step 3 — Connect the bot

Merchant taps **Connect to Bot**, which opens `https://t.me/GabaaPlaceBot?start=link_store_{store_id}` and taps START. This silently registers the merchant with the bot.

### Step 4 — Link your group or channel

Merchant adds @GabaaPlaceBot as **Administrator** in their target group or channel. The bot links the chat to the store and sends a private confirmation message. **Nothing is posted publicly during linking.**

### Step 5 — Add products

Merchant uploads images, sets price and stock, saves as **draft**, then **publishes**. Publishing pushes the product to the Telegram channel with an Order Now button.

### Step 6 — Fulfill orders

Orders appear in the merchant dashboard. Merchant ships, marks **delivered**, and earnings credit to the store wallet.

---

## 7. How It Works — Buyer Journey

Shorter flow for a "Shopper experience" subsection.

1. **See** a product post or story in the Telegram channel
2. **Tap** Order Now or open the Mini App from the group menu
3. **Browse** catalog, view stories, save favorites
4. **Add to cart** and enter shipping address
5. **Checkout** — stock reserved, order created
6. **Track** order status; cancel if still pending

---

## 8. Website Information Architecture

### Recommended sitemap

```
/                          → Home (hero, value props, how it works, CTA)
/features                  → Full feature breakdown (seller + buyer)
/how-it-works              → Detailed seller onboarding steps
/pricing                   → Free tier (+ future paid plans if any)
/for-sellers               → Merchant-focused landing (optional A/B vs home)
/stories                   → Explain Product Stories marketing feature
/roadmap                   → Planned features (discounts, analytics, payments)
/faq                       → Common questions
/about                     → Mission, Telegram commerce thesis
/contact                   → Support email, Telegram link
/blog                      → Optional — seller tips, launch stories
/legal/privacy             → Privacy policy
/legal/terms               → Terms of service
```

### Primary navigation

```
Logo | Features | How It Works | Pricing | FAQ | [Open in Telegram] CTA
```

### Footer columns

| Product | Resources | Legal | Connect |
|---------|-----------|-------|---------|
| Features | Docs (link to repo docs) | Privacy | @GabaaPlaceBot |
| How it works | Merchant guide | Terms | Telegram channel |
| Pricing | API / Swagger | | Contact email |
| Roadmap | Mini App link | | |

---

## 9. Page-by-Page Design Specifications

### 9.1 Homepage `/`

**Goal:** Convert Telegram group owners to try Gabaa within 30 seconds of landing.

#### Section order

| # | Section | Content |
|---|---------|---------|
| 1 | **Hero** | Headline, subhead, primary CTA ("Start selling on Telegram"), secondary CTA ("See how it works"). Hero visual: phone mockup showing Telegram channel + Mini App overlay. |
| 2 | **Social proof strip** | "Free to start" · "No coding required" · "Built for Ethiopia" · "Powered by Telegram Mini Apps" |
| 3 | **5 value pillars** | Cards for Free setup, Sales tracking, Stock management, Community selling, Marketing tools — each linking to `/features#anchor` |
| 4 | **How it works** | 4-step condensed version with icons |
| 5 | **Product Stories highlight** | Visual of story carousel; "Market like Instagram, sell like a store" |
| 6 | **Seller dashboard preview** | Mockup: orders list, wallet balance, product grid |
| 7 | **Buyer experience** | Short strip — cart, checkout, addresses |
| 8 | **Categories we support** | Electronics, Fashion, Home & Garden, Beauty, Toys |
| 9 | **CTA banner** | "Ready to turn your group into a store?" + Telegram deep link |
| 10 | **Footer** | Standard |

#### Hero copy (recommended)

**Headline:** Turn your Telegram group into a store — free, in minutes.  
**Subhead:** Gabaa Place lets you sell to the community you already built. Add the bot, list products, take orders, and track earnings — without coding or leaving Telegram.

---

### 9.2 Features `/features`

**Goal:** Exhaustive but scannable feature reference for evaluators.

Structure as **tabbed or anchored sections**:

1. **Store setup** — silent linking, profile, categories
2. **Products & inventory** — CRUD, draft/publish, stock, images
3. **Orders & fulfillment** — cart, checkout, lifecycle, addresses
4. **Sales & wallet** — order dashboard, wallet balance, earnings
5. **Marketing** — stories, channel posts, (planned) discounts
6. **Telegram native** — Mini App, deep links, group/channel posts
7. **For buyers** — browse, favorites, cart, orders

Each feature block:

```
[Icon] Feature name
One-sentence benefit
2–3 bullet details
[Screenshot or illustration placeholder]
Status badge: Live | Coming soon
```

---

### 9.3 How It Works `/how-it-works`

**Goal:** Remove fear of technical setup.

- Full 6-step seller journey (Section 6) with screenshots
- Troubleshooting accordion (from merchant onboarding guide):
  - Bot not linking? → Did you tap START in private chat first?
  - Wrong store linked? → Use the correct `link_store_{id}` link
  - Mini App not loading? → Bot must remain admin
- Embedded video placeholder (future: 60-second setup walkthrough)

---

### 9.4 Pricing `/pricing`

**Goal:** Emphasize free entry; leave room for future tiers.

| Tier | Price | Includes |
|------|-------|----------|
| **Starter** | Free | Store setup, unlimited products, orders, stories, wallet |
| **Growth** | TBD | Analytics dashboard, boosted listings, discount campaigns |
| **Payments** | TBD | ArifPay / online payment integration |

> Only list **Starter** features as live today. Mark Growth and Payments as "Coming soon" to align with Roadmap.

---

### 9.5 Product Stories `/stories` (or section under Features)

**Goal:** Showcase the marketing differentiator.

Content blocks:

- What is a Product Story? (timed, swipeable product ad)
- How merchants create one (pick product → upload media → set dates → publish)
- How buyers discover stories (public story feed in Mini App)
- View count as engagement metric
- Comparison: Story vs channel post (story = campaign; post = catalog listing)

---

### 9.6 FAQ `/faq`

See [Section 13](#13-faq-content) — use accordion UI.

---

### 9.7 About `/about`

- Mission: democratize e-commerce for Telegram communities
- Why Telegram: 900M+ users, built-in trust, Mini App platform
- Ethiopia focus: ETB, local phone formats, Addis Ababa logistics context
- Open-source backend, commercial Mini App (adjust per actual licensing)

---

## 10. Visual & UX Design Guidelines

### Design direction

| Principle | Guidance |
|-----------|----------|
| **Feel** | Modern, mobile-first, trustworthy — not flashy crypto aesthetic |
| **Primary audience device** | Mobile (merchants manage from phone) |
| **Reference apps** | Telegram UI patterns, Instagram Stories (for story feature), clean SaaS landing pages |
| **Color palette suggestion** | Telegram blue (`#2AABEE`) as accent; warm neutral background; ETB/green for success states |
| **Typography** | Clean sans-serif (Inter, DM Sans, or similar); large headlines on mobile |
| **Imagery** | Real phone mockups in Telegram dark/light mode; Ethiopian merchant context where possible |
| **Icons** | Simple line icons per feature pillar |

### Component patterns

- **Feature cards** — icon, title, 2-line description, "Learn more" link
- **Step timeline** — numbered vertical steps on mobile; horizontal on desktop
- **Dashboard mockup** — framed browser/Telegram Mini App window; blur sensitive data
- **Story carousel mockup** — rounded rectangle mimicking IG/TG story UI
- **Status badges** — `Live` (green), `Coming soon` (amber)

### Accessibility

- WCAG AA contrast on all text
- Alt text on all product/UI mockups
- Mobile tap targets ≥ 44px for CTAs

### Responsive breakpoints

| Breakpoint | Layout |
|------------|--------|
| `< 640px` | Single column, sticky bottom CTA |
| `640–1024px` | Two-column feature grids |
| `> 1024px` | Full hero with side-by-side mockup |

---

## 11. Messaging & Copy Bank

### Headlines by page

| Page | Headline |
|------|----------|
| Home | Turn your Telegram group into a store |
| Features | Everything you need to sell on Telegram |
| How it works | Live in 5 minutes. No code. No new app. |
| Stories | Product ads your customers actually see |
| Pricing | Start free. Grow when you're ready. |

### Short feature blurbs (for cards)

| Feature | Blurb |
|---------|-------|
| Free setup | Create your store and connect your group without paying a cent or writing code. |
| Silent linking | Add the bot as admin — your store links automatically. No spam posts. |
| Channel posts | Publish a product and it appears in your channel with a buy button. |
| Stock sync | Inventory updates on every order. Cancel restores stock automatically. |
| Order dashboard | See pending, shipped, and delivered orders with customer details. |
| Store wallet | Track earnings as you fulfill orders. |
| Product stories | Run timed swipeable ads for any product in your catalog. |
| Mini App checkout | Customers browse and pay inside Telegram — no external site. |
| Favorites | Shoppers save products and come back later. |
| Addresses | One-tap checkout with saved delivery addresses. |

### Words to use

store, group, channel, sell, orders, stock, Telegram, free, minutes, community, publish, story, wallet, checkout

### Words to avoid

blockchain, Web3, enterprise-grade, synergy, leverage (as buzzword), deploy, repository, API (on marketing pages)

---

## 12. Calls to Action & Deep Links

| CTA label | Destination | Use on |
|-----------|-------------|--------|
| **Start selling on Telegram** | `https://t.me/GabaaPlaceBot` | Hero, footer, pricing |
| **Open Mini App** | `https://gabaa-web.vercel.app` | Secondary CTA (works best inside Telegram) |
| **Try the bot** | `https://t.me/GabaaPlaceBot?start=link_store_demo` | Demo flow (needs live demo store) |
| **View documentation** | Link to `docs/merchant_onboarding_guide.md` in repo | Footer / developers |

> **Note:** The Mini App is designed to run inside Telegram. On the marketing site, prefer the bot link as primary CTA; add a QR code on desktop for "Scan to open in Telegram."

---

## 13. FAQ Content

**Is Gabaa Place really free?**  
Yes. You can create a store, list products, take orders, and use product stories at no cost on the Starter plan.

**Do I need a website or coding skills?**  
No. Everything runs through Telegram and the Gabaa Mini App. If you can manage a Telegram group, you can run a store.

**How do I connect my group?**  
Create your store in the Mini App, tap Connect to Bot, press START in the private chat with @GabaaPlaceBot, then add the bot as an administrator in your group or channel.

**Will the bot post spam in my group?**  
Linking is silent. The bot only posts when you publish a product.

**Can I use a channel, not just a group?**  
Yes. Gabaa supports both Telegram groups and channels.

**How do customers pay?**  
Today, orders are placed through the Mini App and merchants coordinate payment and fulfillment. Integrated online payments (ArifPay) are on the roadmap.

**How do I track sales?**  
Your merchant dashboard lists every order with status and total. Your store wallet shows accumulated earnings from delivered orders.

**What happens when stock runs out?**  
Checkout blocks orders when stock is insufficient. You can filter your catalog to see low-stock and out-of-stock items.

**What are Product Stories?**  
Timed, swipeable product ads — like stories on social media — that link directly to a product in your store. You set the media, caption, and active dates.

**Are discounts supported?**  
Not yet. Coupon and discount campaigns are planned; sign up via the bot to be notified.

**Is my data secure?**  
Authentication uses Telegram's official WebApp validation. API access is protected with JWT tokens.

---

## 14. SEO & Metadata

### Primary keywords

- Telegram store
- Telegram e-commerce
- Telegram Mini App shop
- sell on Telegram Ethiopia
- Telegram group storefront
- no-code Telegram shop

### Page titles & descriptions

| Page | Title | Meta description |
|------|-------|------------------|
| Home | Gabaa Place — Free Telegram Store for Your Group | Turn your Telegram group or channel into a store. Free setup, inventory, orders, and product stories. No coding required. |
| Features | Features — Gabaa Place | Product catalog, stock management, order tracking, wallet, stories, and Mini App checkout for Telegram sellers. |
| How it works | How It Works — Gabaa Place | Launch your Telegram store in minutes: create, connect bot, add products, sell. |

### Open Graph

- `og:image` — 1200×630 hero graphic with phone mockup and logo
- `og:site_name` — Gabaa Place
- `twitter:card` — summary_large_image

---

## 15. Technical Context for Developers

This section helps the website team integrate honestly with the live product.

### Stack (existing)

| Layer | Technology |
|-------|------------|
| Backend API | Go, Gin, GORM, PostgreSQL |
| Bot | telebot v4, @GabaaPlaceBot |
| Mini App frontend | Hosted at `gabaa-web.vercel.app` (separate repo) |
| Media | Cloudinary (`gabaa_products` folder) |
| Auth | JWT from `POST /auth/telegram` |

### What the marketing site does NOT need

- Direct API integration for v1 (static site is fine)
- User authentication
- Product catalog embedding (optional future: public `GET /products` widget)

### Optional live integrations (v2)

| Integration | Endpoint | Use on website |
|-------------|----------|----------------|
| Public product showcase | `GET /products` | "Explore stores" section |
| Active stories demo | `GET /stories` | Live story feed widget |
| Swagger docs link | `/swagger/index.html` | Developer footer link |

### Environment references

```
Bot:      https://t.me/GabaaPlaceBot
Mini App: https://gabaa-web.vercel.app
API:      (deployment URL — configure per environment)
```

---

## 16. Roadmap — Planned Features

List on `/roadmap` and mark as **Coming soon** on the website. Do not claim these as live.

| Feature | User benefit | Current state |
|---------|--------------|---------------|
| **Discounts & coupon codes** | Run promotions and sales campaigns | Not implemented |
| **Analytics dashboard API** | Daily/weekly sales charts, store views | Store views tracked in DB only |
| **ArifPay payments** | Customers pay online at checkout | Stub only |
| **Product boosting** | Promote listings in catalog | DB column only |
| **Email/password login** | Web dashboard without Telegram | Schema only |
| **Telegram Login Widget** | Web auth via Telegram | Documented, not built |
| **Multi-chat linking** | One store → multiple groups | `linked_chats` table unused |

---

## 17. Appendix — API & Integration References

For the website team and technical writers. Full docs live in this repository.

| Document | Path | Contents |
|----------|------|----------|
| Merchant onboarding | `docs/merchant_onboarding_guide.md` | Store setup + bot linking |
| Frontend integration | `docs/frontend_integration_guide.md` | Full API overview |
| Store setup (UI) | `docs/store_setup_frontend_integration.md` | Dashboard types |
| Products | `docs/merchant_product_integration.md` | Upload, draft/publish |
| Orders | `docs/merchant_order_integration.md` | Merchant order lifecycle |
| Checkout | `docs/customer_checkout_integration.md` | Addresses + checkout |
| Cart | `docs/cart_integration_guide.md` | Cart API |
| OpenAPI | `docs/swagger.yaml` | Generated API spec |

### Key API surface (marketing summary)

**Public:** auth via Telegram, product catalog, categories, stories, image upload  
**Merchant (JWT):** store CRUD, products, orders, stories, wallet  
**Customer (JWT):** cart, addresses, orders, favorites  

### Order status vocabulary

`pending` → `shipped` → `delivered` | `cancelled`

### Product status vocabulary

`draft` → `published` → `archived`

### Store status vocabulary

`pending` → `launched`

---

## Website Build Checklist

Use this before launch:

- [ ] Hero CTA links to `https://t.me/GabaaPlaceBot`
- [ ] All five value pillars represented on homepage
- [ ] "How it works" matches actual 6-step onboarding flow
- [ ] Product Stories section reflects live feature (schedule, views, media)
- [ ] Discounts marked as **Coming soon**, not live
- [ ] Pricing page only promises free-tier live features
- [ ] FAQ answers match real product behavior
- [ ] Mobile layout tested (primary user device)
- [ ] Open Graph image and meta descriptions set
- [ ] Privacy policy and terms pages exist
- [ ] Footer links to Mini App and bot
- [ ] Ethiopia / ETB context in examples where relevant
- [ ] No claims of analytics dashboard API until shipped
- [ ] Screenshots/mockups use "Gabaa Place" branding consistently

---

## Document Maintenance

| When | Action |
|------|--------|
| New API feature ships | Update Section 5 status table and relevant page spec |
| Feature removed | Mark deprecated; update FAQ |
| Rebrand | Update Section 2 and all copy bank entries |
| New market | Add locale section under Brand |

---

**Gabaa Place** — *Sell where you already connect.*

*Backend repository: `gabaa-bot` · Bot: [@GabaaPlaceBot](https://t.me/GabaaPlaceBot) · Mini App: [gabaa-web.vercel.app](https://gabaa-web.vercel.app)*
