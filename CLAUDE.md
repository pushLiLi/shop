# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

BYCIGAR e-commerce platform — a cigar shop migrated from a static HTML site to Vue 3 + Go. The project is feature-complete with products, categories, cart, orders, favorites, user auth, admin panel, CMS pages, and site settings.

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Vue 3 (Composition API, `<script setup>`) + Vite + Vue Router + Pinia |
| Backend | Go 1.23 + Gin + GORM + JWT + Swagger |
| Database | MySQL 8.4 (Docker) |
| Deployment | Docker Compose (MySQL + Go backend + Vue frontend) |
| Language | JavaScript only (no TypeScript) |

## Common Commands

### Frontend (`bycigar-vue/`)
```bash
cd bycigar-vue
npm run dev      # Dev server at localhost:5173, proxies /api -> localhost:3000
npm run build    # Production build to dist/
npm run preview  # Preview production build
```

### Backend (`server-go/`)
```bash
cd server-go
go run ./cmd/main.go    # Start API server at localhost:3000
go build -o server ./cmd/main.go   # Build binary
```

### Docker (full stack)
```bash
docker-compose up --build    # MySQL:3306 + Backend:3000 + Frontend:80
```

### Database
```bash
docker-compose up mysql      # Start MySQL only (root/123456, database: bycigar)
```
Auto-migrates all tables on backend startup. Seeds admin user and default data on first run. MySQL uses healthcheck (`mysqladmin ping`) with `my.cnf` forcing utf8mb4 charset.

## Architecture

### Frontend (`bycigar-vue/src/`)

- **Entry**: `main.js` → installs Pinia + Vue Router
- **Layout**: `App.vue` — shows TheHeader/TheFooter on all routes except admin. Global Toast + CartDrawer always present. Dark theme (#0f0f0f background).
- **Router**: `router/index.js` — 17 routes. Navigation guard checks `localStorage` for JWT token and user role.
  - Public: `/`, `/products`, `/category/:slug`, `/products/:id`, `/search`, `/:slug(about|services|privacy-policy|statement)`
  - Auth required: `/checkout`, `/orders`, `/favorites`, `/profile`
  - Admin (`/admin/*`): requires `role === 'admin'`
- **Stores** (Pinia): `auth.js` (Composition API), `cart.js` (Options API), `favorites.js` (Options API), `useSettingsStore.js` (Options API), `toast.js` (Composition API). Each store defines its own `getAuthHeaders()` reading `localStorage.getItem('token')` directly. API calls hardcode `http://localhost:3000/api/*` (Vite proxy at `/api` exists in config but stores bypass it).
- **Category sidebar**: `components/CategorySidebar.vue` — dynamic sidebar on product listing page (`CategoryView.vue`), fetches `GET /api/categories` (returns top-level categories with nested `children`). Mobile: horizontal pill tags.
- **Global CSS**: `style.css` — dark theme, grid system (.col-2/3/4/6/12), utility classes, 768px responsive breakpoint.
- **Markdown**: `marked` library renders CMS page content.

### Backend (`server-go/`)

- **Entry**: `cmd/main.go` — loads config → connects MySQL → migrations → seeds → InitSnowflake → BackfillOrderNo → Gin router
- **Structure**: `internal/config/`, `internal/database/`, `internal/models/`, `internal/handlers/`, `internal/middleware/`
- **Config**: `.env` file (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET, SERVER_PORT)
- **Auth middleware** (`internal/middleware/auth.go`):
  - `AuthMiddleware()` — optional JWT parsing; supports `Bearer {token}` and `user-{id}` dev bypass. Sets `userID` in Gin context.
  - `RequireAuth()` — checks `userID` exists + DB lookup to verify user still exists. Returns 401 if missing.
  - `AdminOnly` — standalone function (not wrapped). Checks `userID`, does DB lookup, verifies `role === 'admin'`, sets `c.Set("user", user)`.
- **Captcha**: `internal/handlers/captcha.go` — `base64Captcha` library generates 4-digit numeric image captcha with 5-min TTL.
- **Response helpers** (`pkg/utils/response.go`): `SuccessResponse`, `CreatedResponse`, `ErrorResponse`, `Success` — standard JSON response wrappers.
- **Snowflake ID** (`pkg/utils/snowflake.go`): `InitSnowflake(nodeID)` and `GenerateOrderNo()` — generates unique order numbers using `bwmarrin/snowflake` library.
- **Static files**: `static/` — uploaded images at `/static/uploads/`, original site assets at `/static/themes/`
- **Swagger docs**: Generated in `docs/`, served at `/swagger/index.html`

### Default Credentials
- Admin: `admin@admin.com` / `123456` (seeded on first run)

## Key Conventions

- **Language**: Communicate with the user in Chinese
- **No code comments** in source files
- **Keep responses concise** (max 4 lines when possible)
- **No tests or linter** configured
- **Pinia stores**: Mix of Composition API (`auth.js`, `toast.js`) and Options API (`cart.js`, `favorites.js`, `useSettingsStore.js`). Each defines its own `getAuthHeaders()` helper reading from `localStorage` directly.
- **Image paths**: Stored as relative paths, served from `static/` directory
- **API base**: All frontend API calls hardcode `http://localhost:3000/api/*` in stores. Production uses nginx to proxy `/api` to backend (see `bycigar-vue/nginx.conf`).
- **Stock sorting**: Product listings always sort `stock > 0` items first, regardless of user-selected sort order (done in SQL via CASE expression).
- **Order numbers**: Orders use snowflake IDs as display order numbers. `GetOrder` endpoint accepts both numeric ID and string orderNo. Existing orders backfilled on startup.
- **Admin vs Public API separation**: Admin list endpoints must return ALL records (no `is_active` filter). Public endpoints filter by `is_active`. Each needs its own handler — never reuse the public handler on admin routes.
- **Admin panel design patterns**: All admin pages use consistent patterns from `AdminProducts.vue`: switch toggle for boolean fields, `.form-section` with `.section-title` (left border accent) for grouping, `.input-group` with `.input-prefix`/`.input-suffix` for decorated inputs, `.badge-success`/`.badge-danger` for status display. CSS is scoped per component with no shared stylesheet — styles are duplicated across admin pages.
- **Admin panel theme**: Light theme (`#fff`/`#f5f5f5` background) separate from storefront dark theme. Accent color `#d4a574`.

## Data Models

Core models with relationships:
- **User** — email, password, name, role (customer/admin)
- **Product** — slug, price, stock, images, belongs to Category, soft delete
- **Category** — self-referential parent (ParentID), has many Products, soft delete
- **Order** → **OrderItem** — Order has snowflake-generated `OrderNo` (unique, user-facing), auto-increment `ID` is internal. Belongs to User and Address.
- **CartItem** / **Favorite** — User ↔ Product relationships
- **Address** — belongs to User, has IsDefault
- **Banner** — Go field `Image` has JSON tag `"imageUrl"`. Frontend must send `image` in request body but reads `imageUrl` from response.
- **Page** / **Setting** — site content management
