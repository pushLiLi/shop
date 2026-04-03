# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

BYCIGAR e-commerce platform — a cigar shop migrated from a static HTML site to Vue 3 + Go. The project is feature-complete with products, categories, cart, orders, favorites, user auth, admin panel, CMS pages, site settings, chat system, notifications, and payment proof management.

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Vue 3 (Composition API, `<script setup>`) + Vite + Vue Router + Pinia |
| Backend | Go 1.23 + Gin + GORM + JWT + Swagger |
| Database | MySQL 8.4 (Docker) |
| File Storage | Local filesystem (`/opt/bycigar/uploads/`) |
| Deployment | Hybrid: MySQL (Docker) + Go binary (systemd) + nginx (host) |
| Language | JavaScript only (no TypeScript) |

## Common Commands

### Frontend (`bycigar-vue/`)
```bash
cd bycigar-vue
npm run dev      # Dev server at localhost:5173, proxies /api -> localhost:3000
npm run build    # Production build to dist/
npm run preview  # Preview production build
npm run test     # Run tests with vitest
npm run test:watch  # Run tests in watch mode
```

### Backend (`server-go/`)
```bash
cd server-go
go run ./cmd/main.go    # Start API server at localhost:3000
go build -o server ./cmd/main.go   # Build binary
```

### Docker (MySQL only)
```bash
docker compose up -d         # MySQL only (127.0.0.1:3306)
```

### Database
```bash
docker compose up -d         # Start MySQL (root/123456, database: bycigar)
```
Auto-migrates all tables on backend startup. Seeds admin user and default data on first run.

## Architecture

### Backend (`server-go/`)

- **Entry**: `cmd/main.go` — loads config → connects MySQL → migrations → seeds → InitStorage → InitSnowflake → BackfillOrderNo → start cleanup jobs → init WebSocket Hub → Gin router
- **Structure**: `internal/config/`, `internal/database/`, `internal/models/`, `internal/handlers/`, `internal/middleware/`, `internal/ws/`
- **Route registration**: `cmd/main.go` registers routes in groups — public routes directly on `router.Group("/api")`, admin routes on `router.Group("/api/admin")` with `AuthMiddleware` + `AdminOnly`. SuperAdminOnly applied per-handler for sensitive endpoints.
- **Config**: `.env` file (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET, SERVER_HOST, SERVER_PORT, UPLOAD_DIR)
- **Auth middleware** (`internal/middleware/auth.go`):
  - `AuthMiddleware()` — optional JWT parsing; supports `Bearer {token}` and `user-{id}` dev bypass. Sets `userID` in Gin context.
  - `RequireAuth()` — checks `userID` exists + DB lookup. Returns 401 if missing.
- **Admin middleware** (`internal/middleware/admin.go`):
  - `AdminOnly` — allows `admin` and `service` roles. Sets `c.Set("user", user)`.
  - `SuperAdminOnly` — allows `admin` role only. Used for banners, pages, config, settings, user roles, payment methods.
- **WebSocket Hub** (`internal/ws/hub.go`): Maintains `CustomerConns` and `AdminConns` maps. Provides `SendToUser`, `SendToAdmins`, `SendToAll`. Uses gorilla/websocket with ping/pong heartbeat (54s interval). `ServiceOnline` map tracks explicitly online admins — online status is controlled manually via `service_online`/`service_offline` WS messages, not by connection state.
- **Chat handlers**: `chat.go` (customer API), `admin_chat.go` (admin API), `ws_chat.go` (WebSocket endpoints). Customer and admin each have their own WS endpoint. Conversation assignment: `assigned_to` field with `AssignedUser` preload.
- **MinIO** (`pkg/minio/minio.go`): Upload handler stores images in MinIO bucket. Admin upload at `POST /api/admin/upload`.
- **Background jobs**: `StartChatCleanup` and `StartNotificationCleanup` run on startup for auto-cleanup of old data.
- **Response helpers** (`pkg/utils/response.go`): `SuccessResponse`, `CreatedResponse`, `ErrorResponse`, `Success`.
- **Snowflake ID** (`pkg/utils/snowflake.go`): `GenerateOrderNo()` for unique order numbers.
- **Static files**: `static/` — legacy assets at `/static/themes/`. New uploads go to MinIO.
- **Swagger docs**: Generated in `docs/`, served at `/swagger/index.html`

### Frontend (`bycigar-vue/src/`)

- **Entry**: `main.js` → installs Pinia + Vue Router
- **Layout**: `App.vue` — shows TheHeader/TheFooter on all routes except admin. Global Toast + CartDrawer + ChatWidget always present. Dark theme (#0f0f0f background).
- **Router**: `router/index.js`. Navigation guard checks `localStorage` for JWT token and user role.
  - Public: `/`, `/products`, `/category/:slug`, `/products/:id`, `/search`, `/:slug(about|services|privacy-policy|statement)`
  - Auth required: `/checkout`, `/orders`, `/favorites`, `/profile`, `/cart`
  - Admin (`/admin/*`): requires `role === 'admin'` or `role === 'service'`
- **Stores** (Pinia): `auth.js`, `cart.js`, `favorites.js`, `useSettingsStore.js`, `toast.js`, `chat.js`, `notifications.js`, `contactMethods.js`. Each store defines its own `getAuthHeaders()` reading `localStorage.getItem('token')`. New stores use `/api` base path (not hardcoded localhost).
- **Composables**: `composables/useCarousel.js` — reusable carousel logic (autoplay, touch/swipe, navigation). Pattern for shared component logic.
- **Global CSS**: `style.css` — dark theme, grid system (.col-2/3/4/6/12), utility classes, 768px responsive breakpoint.
- **Markdown**: `marked` library renders CMS page content.
- **Charts**: `chart.js` + `vue-chartjs` for admin dashboard analytics.
- **Image cropping**: `vue-advanced-cropper` for product image editing.

### Default Credentials
- Admin: `admin@admin.com` / `123456` (seeded on first run)

## Important Notes

- **DEVELOPMENT_PLAN.md is outdated**: It describes a planned Node.js/Express backend that was never built. The actual backend is Go/Gin.
- **Service online persistence**: Admin online status is stored in `localStorage` (`chat_service_online`) and restored on WebSocket reconnect. Online state is controlled via explicit `service_online`/`service_offline` WS messages, not by connection lifecycle.

## Key Conventions

- **Language**: Communicate with the user in Chinese
- **No code comments** in source files
- **Keep responses concise** (max 4 lines when possible)
- **No tests or linter** configured — vitest is set up but only a few store spec files exist; no comprehensive test suite
- **Image storage**: New uploads use local filesystem via `pkg/storage/`. Legacy images stored as relative paths in `static/`.
- **API base**: Older stores hardcode `http://localhost:3000/api/*`. Newer stores use `/api` relative path. Production uses nginx to proxy `/api` to backend (see `bycigar-vue/nginx.conf`).
- **Stock sorting**: Product listings always sort `stock > 0` items first (done in SQL via CASE expression).
- **Order numbers**: Orders use snowflake IDs as display order numbers. `GetOrder` endpoint accepts both numeric ID and string orderNo.
- **Admin vs Public API separation**: Admin list endpoints must return ALL records (no `is_active` filter). Public endpoints filter by `is_active`. Each needs its own handler.
- **Three-tier admin system**: `admin` (full access via SuperAdminOnly), `service` (limited admin via AdminOnly — can manage products, orders, chat), `customer` (storefront).
- **Admin panel design patterns**: All admin pages use consistent patterns from `AdminProducts.vue`: switch toggle for boolean fields, `.form-section` with `.section-title` (left border accent) for grouping, `.input-group` with `.input-prefix`/`.input-suffix` for decorated inputs, `.badge-success`/`.badge-danger` for status display. CSS is scoped per component.
- **Admin panel theme**: Light theme (`#fff`/`#f5f5f5` background) separate from storefront dark theme. Accent color `#d4a574`.

## Data Models

Core models with relationships:
- **User** — email, password, name, role (customer/admin/service)
- **Product** — slug, price, stock, images, belongs to Category, soft delete
- **Category** — self-referential parent (ParentID), has many Products, soft delete
- **Order** → **OrderItem** — Order has snowflake-generated `OrderNo` (unique, user-facing), auto-increment `ID` is internal. Belongs to User and Address.
- **CartItem** / **Favorite** — User ↔ Product relationships
- **Address** — belongs to User, has IsDefault
- **Banner** — Go field `Image` has JSON tag `"imageUrl"`. Frontend must send `image` in request body but reads `imageUrl` from response.
- **Page** / **Setting** / **Config** — site content management
- **Conversation** → **Message** — Chat system. Conversation belongs to User, has status (open/closed) and `assigned_to` (nullable, FK to User). Message has SenderType (customer/service/system).
- **QuickReply** — 客服快捷回复，关联创建者 User，支持 sort_order 排序。
- **Notification** — belongs to User, has type (order_status/back_in_stock/price_drop), is_read, composite index (is_read, created_at) for cleanup queries.
- **PaymentMethod** — name, QR code URL, instructions, is_active, sort_order
- **PaymentProof** — belongs to Order and User and PaymentMethod, has status (pending/approved/rejected), reviewer info
