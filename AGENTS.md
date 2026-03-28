# AGENTS.md - AI Coding Agent Guidelines

## Tech Stack

Frontend: Vue 3.5 + Vite 8 + Vue Router 4 + Pinia 3 + marked + vue-advanced-cropper (JavaScript only, no TypeScript)
Backend: Go 1.25 + Gin 1.9 + GORM 1.25 + JWT + Swagger (Dockerfile uses golang:1.23-alpine)
Database: MySQL 8.4 (Docker, utf8mb4) | Object Storage: MinIO (Docker, API:9000, Console:9001)
Module: `bycigar-server`

## Build Commands

```bash
# Frontend (run from bycigar-vue/)
npm run dev                     # Dev server at :5173 (proxies /api -> :3000, /media -> :9000)
npm run build                   # Production build (use this to verify frontend changes)
npm run preview                 # Preview production build locally

# Backend (run from server-go/)
go run ./cmd/main.go            # Dev server at :3000
go build -o nul ./...           # Check Go compilation (Windows; use -o /dev/null on Linux/Mac)
go fmt ./...                    # Format Go code
swag init                       # Generate Swagger docs (run from server-go/)

# Docker
docker-compose up -d mysql      # MySQL only (root/123456, db: bycigar)
docker-compose up -d minio      # MinIO only (minioadmin/minioadmin123, bucket: bycigar)
docker-compose up --build       # Full stack: MySQL + MinIO + Backend + Frontend (:80)
```
**No test or lint framework.** No ESLint, Prettier, or Go linter. Use `go build` and `npm run build` to verify compilation.

## Project Structure

```
bycigar-vue/src/
├── main.js, App.vue            # Entry + layout (dark theme, Toast, CartDrawer)
├── style.css                   # Global dark theme, grid (.col-2/3/4/6/12), 768px breakpoint
├── components/                 # ProductCard, CartDrawer, Toast, TheHeader, TheFooter,
│                               #   AdminImageUpload, AddressForm, CategorySidebar,
│                               #   EditableImage, EditableText
├── composables/                # useCarousel (autoplay, touch, pause)
├── views/                      # 11 public: Home, Category, ProductDetail, Search, Cart,
│                               #   Checkout, Orders, Favorites, Login, Profile, Page
├── views/admin/                # 6 admin: AdminLayout, AdminProducts, AdminBanners,
│                               #   AdminCategories, AdminPages, AdminSettings
├── stores/                     # Pinia: auth.js, cart.js, favorites.js, toast.js, useSettingsStore.js
└── router/index.js             # Routes with auth/admin guards

server-go/
├── cmd/main.go                 # Entry: wires all routes and starts server
├── internal/
│   ├── config/config.go        # Loads .env via godotenv into AppConfig global
│   ├── database/database.go    # Connect, Migrate, Seed, BackfillOrderNo
│   ├── handlers/               # product, banner, page, auth, captcha, cart, favorite,
│   │                           #   order, address, config, setting, upload,
│   │                           #   admin_product, admin_category
│   ├── models/                 # User, Category, Product, CartItem, Favorite, Address,
│   │                           #   Order, OrderItem, SiteConfig, Banner, Page, Setting
│   └── middleware/             # cors.go, auth.go, admin.go
└── pkg/
    ├── minio/minio.go          # MinIO client wrapper + EnsureBucket
    └── utils/                  # Response helpers + snowflake order number + UUID upload names
```

## Vue Code Style

- `<script setup>` Composition API only. No TypeScript.
- Declaration order: imports -> `defineEmits` -> `defineProps` -> stores -> refs/computed -> functions
- Use `useToastStore` for notifications, **never** `alert()`
- CSS: `<style scoped>` with kebab-case class names. Storefront dark theme (#0f0f0f), admin light theme (#fff).
- No code comments. Inline event handlers use named `async function` declarations, not arrow functions.
- **Match Go model JSON tags exactly** (e.g. `"imageUrl"` not `"image"`, `"isFeatured"` not `"featured"`)
- **All API calls use relative paths** (`'/api/...'`), never hardcoded `http://localhost:3000/api/...`. The Vite proxy handles routing in dev; nginx handles it in production.

## Pinia Store Patterns

Two styles exist:
- **Composition API** (`auth.js`, `toast.js`): `defineStore('name', () => { ... })`
- **Options API** (`cart.js`, `favorites.js`, `useSettingsStore.js`): `defineStore('name', { state, getters, actions })`

Store IDs are simple lowercase: `'auth'`, `'cart'`, `'favorites'`, `'toast'`, `'settings'`.
Options API stores define `getAuthHeaders()` at module level reading `localStorage.getItem('token')`.
Auth store exposes `getAuthHeaders` as a returned method using the reactive `token` ref.

## Go Code Style

- Import order: stdlib -> blank line -> external -> blank line -> internal
- Early return for all error handling. Handlers: PascalCase. Helpers: camelCase.
- Admin handlers in separate `admin_*.go` files. Swagger `godoc` comments on all handlers.
- No code comments except Swagger annotations.
- Error responses: `utils.ErrorResponse(c, statusCode, "msg")` or `c.JSON(status, gin.H{"error": "..."})`.
- Input structs live in model files; exception: `ProductInput` in `admin_product.go`.

## GORM Models

- Each model defines base fields inline (no shared `Model` struct). `uint` for all IDs.
- JSON tags: camelCase. GORM tags: snake_case. Always include `CreatedAt`/`UpdatedAt`. Use `DeletedAt` for soft-delete.
- Response and input structs live alongside their models (e.g. `CreateOrderInput` in `models/order.go`).
- Preload: `database.DB.Preload("Category").Find(&products)`

## Response Helpers (pkg/utils/response.go)

```go
utils.SuccessResponse(c, data)                            // 200 + data
utils.CreatedResponse(c, data)                            // 201 + data
utils.Success(c)                                          // 200 {"success": true}
utils.ErrorResponse(c, http.StatusBadRequest, "Invalid")  // 4xx {"error": "msg"}
```

## Naming Conventions

| Type | Convention | Example |
|---|---|---|
| Vue Components | PascalCase | `ProductCard.vue` |
| Views | PascalCase + View | `HomeView.vue` |
| Admin Views | Admin + Pascal | `AdminProducts.vue` |
| Stores | use + camel + Store | `useCartStore` |
| Composables | use + PascalCase | `useCarousel` |
| CSS Classes | kebab-case | `.product-card` |
| Go Handlers | PascalCase | `GetProducts` |
| Go Helpers | camelCase | `generateSlug` |
| Go Files | snake_case | `admin_product.go` |

## Authentication

- Token in `localStorage.getItem('token')`, user JSON in `localStorage.getItem('user')` (has `role`).
- `AuthMiddleware()`: global, optional JWT parsing, sets `c.Set("userID", ...)`.
- `RequireAuth()`: per-route, only on `ChangePassword`. Checks userID + DB lookup.
- `AdminOnly`: plain `func(c *gin.Context)` used as group middleware for `/api/admin`.
- Cart/favorites/addresses/orders check `c.Get("userID")` inline in handlers.
- Dev bypass: `Authorization: user-{id}` header skips JWT.

## API Endpoints

Public: products, categories, banners, pages, auth/login, auth/register, config, settings, auth/captcha
Auth-dependent: auth/me, auth/profile, auth/change-password, cart, favorites, orders, addresses
Admin (AdminOnly): admin/products, admin/categories, admin/banners, admin/pages, admin/upload, admin/settings/:key
**Security note**: `PUT /api/admin/config/:key` is outside admin group -- no role check.
Query params: `page`, `limit`, `search`, `category` (slug), `categoryId`, `sortBy` (alias: `sort`), `sortOrder` (alias: `order`), `featured`, `active`

## Architecture Notes

- Admin endpoints return ALL records; public endpoints filter `WHERE is_active = true`. Never reuse public handlers on admin routes.
- Product listings always sort `stock > 0` first via SQL CASE expression.
- Orders use snowflake `OrderNo` (user-facing) + auto-increment `ID` (internal). `GetOrder` accepts both.
- `Product.Image` and `Banner.Image` have JSON tag `"imageUrl"`. Frontend must send/read `"imageUrl"`.
- Config API: `GET /api/config` returns `{"key": "value"}` map. Settings API: `GET /api/settings` returns `{"success": true, "data": {...}}`.
- Image upload: `POST /api/admin/upload` multipart `file` -> `{"success": true, "url": "/media/{bucket}/{timestamp}_{uuid}{ext}"}`. Max 10MB, jpg/png/gif/webp.
- `AdminImageUpload.vue` supports optional image cropping via `vue-advanced-cropper`. Pass `:aspect-ratio` prop to enforce ratio (e.g. `1` for products, `7/3` for banners). `null` (default) skips cropping.
- CMS pages: slugs `about`, `services`, `privacy-policy`, `statement`. Markdown via `marked`.
- Cart debounce: 300ms. Address limit: 5 per user. Captcha TTL: 5 min.
- **ProductCard horizontal mode**: `ProductCard` accepts `horizontal` prop (Boolean). CategoryView/SearchView pass `:horizontal="isCompact"` (where `isCompact = window.innerWidth <= 992`) to show left-image/right-text layout on tablet/mobile. HomeView and ProductDetailView do not pass it (vertical cards). Horizontal styles use `.product-card.horizontal` class with `@media (max-width: 768px)` for smaller image size (120px vs 180px).
- **Responsive grids**: CategoryView/SearchView use `grid-template-columns: 1fr` at ≤992px for horizontal cards; HomeView uses `repeat(2, 1fr)` at ≤768px for vertical cards. When using `flex-direction: column` on a parent flex container, always set `align-items: stretch` so children fill the full width.
- **CategorySidebar mobile drawer**: On mobile (≤768px), CategorySidebar shows a button that opens a bottom drawer via `<Teleport to="body">`. Uses `drawerOpen` ref + `openDrawer()`/`closeDrawer()` methods. Drawer has overlay (`rgba(0,0,0,0.6)`), slides up with `transform: translateY(100%)` animation. Categories with children expand inline via `expandedCategories` Set. Always set `document.body.style.overflow = 'hidden'` when drawer opens to prevent background scroll.

## Git Commits

```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
refactor: 简化购物车侧边栏逻辑
```

## Important Rules

- **Language**: Communicate in Chinese
- **No comments**: Do not add code comments in source files
- **No tests**: No test framework, use build to verify
- **Simplicity**: Avoid over-engineering
- **Error handling**: Frontend uses toast, backend uses `utils.ErrorResponse()`
- **Database**: Always `utf8mb4`, `parseTime=True&loc=Local` in DSN
- **API paths**: Frontend always uses relative `/api/...` paths, never hardcoded localhost URLs
- **Concise responses**: Keep responses under 4 lines
