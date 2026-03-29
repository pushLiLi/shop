# AGENTS.md - AI Coding Agent Guidelines

## Tech Stack

Frontend: Vue 3.5 + Vite 8 + Vue Router 4 + Pinia 3 + marked + vue-advanced-cropper (JavaScript only, no TypeScript, no axios — uses native `fetch`)
Backend: Go 1.25 + Gin 1.9 + GORM 1.25 + JWT + Swagger (Dockerfile: golang:1.23-alpine)
Database: MySQL 8.4 (Docker, utf8mb4) | Object Storage: MinIO (Docker, API:9000, Console:9001)
Module: `bycigar-server`

## Build Commands

```bash
# Frontend (run from bycigar-vue/)
npm run dev                     # Dev server at :5173 (proxies /api -> :3000, /media -> :9000)
npm run build                   # Production build (use this to verify frontend changes)

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
No ESLint, Prettier, or Go linter. Use `go build` and `npm run build` to verify compilation.

## Test Commands

```bash
# Frontend tests (run from bycigar-vue/) — Vitest + @vue/test-utils + happy-dom
npm test                        # Run all tests once (vitest run)
npm run test:watch              # Run tests in watch mode
npx vitest run src/stores/cart.spec.js          # Run a single test file
npx vitest run --reporter=verbose               # Verbose output

# Backend tests (run from server-go/) — Go testing + testify/suite
# Requires MySQL + MinIO running (docker-compose up -d mysql minio)
# Uses database: bycigar_test (separate from dev db: bycigar)
go test ./... -v -count=1                       # Run all backend tests
go test ./test/... -v -count=1                  # Run only integration tests
go test ./test/... -v -count=1 -run TestProductSuite  # Run a single test suite
go test ./test/... -v -count=1 -run "TestProductSuite/TestGetProductsDefaultPagination"  # Single test
```
Frontend tests are unit tests (no backend needed). Backend tests are integration tests hitting real MySQL + MinIO.

## Project Structure

```
bycigar-vue/src/
├── main.js, App.vue            # Entry + layout (dark theme, Toast, CartDrawer, page transitions)
├── style.css                   # Global dark theme, grid (.col-2/3/4/6/12), utils, 768px breakpoint
├── components/                 # ProductCard, CartDrawer, Toast, TheHeader, TheFooter,
│                               #   AdminImageUpload, AddressForm, CategorySidebar,
│                               #   EditableImage, EditableText
├── composables/                # useCarousel (autoplay, touch, pause)
├── views/                      # Home, Category, ProductDetail, Search, Cart,
│                               #   Checkout, Orders, Favorites, Login, Profile, Page
├── views/admin/                # AdminLayout + Dashboard, Products, Orders, Users,
│                               #   Banners, Categories, Pages, Settings
├── stores/                     # Pinia: auth.js, cart.js, favorites.js, toast.js, useSettingsStore.js
│                               #   (+ colocated .spec.js test files)
└── router/index.js             # Routes with auth/admin/superAdmin guards (createWebHistory)

server-go/
├── cmd/main.go                 # Entry: LoadConfig -> Connect -> Migrate -> Seed -> routes -> Run
├── internal/
│   ├── config/config.go        # Loads .env via godotenv into AppConfig global
│   ├── database/database.go    # Connect, Migrate, Seed, BackfillOrderNo
│   ├── database/seed_data.go   # SeedTestData: clears + re-inserts all test data on startup
│   ├── handlers/               # product, banner, page, auth, captcha, cart, favorite,
│   │                           #   order, address, config, setting, upload, testing,
│   │                           #   admin_product, admin_category, admin_order,
│   │                           #   admin_dashboard, admin_user
│   ├── models/                 # User, Category, Product, CartItem, Favorite, Address,
│   │                           #   Order, OrderItem, SiteConfig, Banner, Page, Setting
│   └── middleware/             # cors.go, auth.go, admin.go (AdminOnly + SuperAdminOnly)
├── test/                       # Integration tests: helpers.go + 12 *_test.go files
└── pkg/
    ├── image/thumbnail.go      # Image processing: resize (max 1200px), thumbnail (300x300), JPEG conversion
    ├── minio/minio.go          # MinIO client wrapper + EnsureBucket
    └── utils/                  # Response helpers + snowflake order number + UUID upload names
```

## Roles & Authorization

Three user roles in `User.Role` field: `"admin"` (super admin), `"service"` (customer service), `"customer"` (regular user).

### Backend Middleware
- **`AdminOnly`**: Allows `admin` + `service` — shared admin routes (products, categories, orders, dashboard, users, upload)
- **`SuperAdminOnly`**: Allows `admin` only — restricted routes (banners, pages, config, settings, user role changes)

### Frontend Auth Store (`useAuthStore`)
- `isAdmin`: `role === 'admin' || role === 'service'` — can access admin panel
- `isSuperAdmin`: `role === 'admin'` — can see revenue data, banners, pages, settings menus
- `isService`: `role === 'service'` — for conditional UI in shared views

### Frontend Route Guards
- `meta.requiresAdmin`: allows admin + service (applied to parent `/admin` route)
- `meta.requiresSuperAdmin`: allows admin only (applied to banners, pages, settings child routes)

### Permission Matrix
| Module | admin | service | customer |
|---|---|---|---|
| Dashboard (full) | ✅ | ✅ (no revenue) | ❌ |
| Products CRUD | ✅ | ✅ | ❌ |
| Categories CRUD | ✅ | ✅ | ❌ |
| Orders (view + status) | ✅ | ✅ | ❌ |
| Users (view + detail) | ✅ | ✅ | ❌ |
| Reset user password | ✅ | ✅ | ❌ |
| Change user role | ✅ | ❌ | ❌ |
| Banners | ✅ | ❌ | ❌ |
| Pages | ✅ | ❌ | ❌ |
| Config / Settings | ✅ | ❌ | ❌ |

## Vue Code Style

- `<script setup>` Composition API only. No TypeScript. No import alias (`@/`) — use relative paths.
- Declaration order: imports -> `defineEmits` -> `defineProps` -> stores -> refs/computed -> functions
- Use `useToastStore` for notifications, **never** `alert()`. Icons are inline SVG, no icon library.
- CSS: `<style scoped>` with kebab-case classes. Storefront dark theme (#0f0f0f bg, #d4a574 accent). Admin light theme (#fff).
- No code comments. Inline event handlers use named `async function` declarations, not arrow functions.
- **Match Go model JSON tags exactly** (e.g. `"imageUrl"` not `"image"`, `"isFeatured"` not `"featured"`)
- **All API calls use relative paths** (`'/api/...'`), never hardcoded localhost URLs.
- **Auth guard pattern**: Before cart/favorites actions, check `authStore.isLoggedIn`. If false, `router.push('/login')`.
- **Images in list views**: Use `product.thumbnailUrl || product.imageUrl` with `loading="lazy"`. Detail pages use full `imageUrl`.

## Go Code Style

- Import order: stdlib -> blank line -> external -> blank line -> internal
- Early return for all error handling. Handlers: PascalCase. Helpers: camelCase.
- Admin handlers in separate `admin_*.go` files. Swagger `godoc` comments on all handlers.
- No code comments except Swagger annotations.
- Error responses: `utils.ErrorResponse(c, statusCode, "msg")` or `c.JSON(status, gin.H{"error": "..."})`.
- Input structs live in model files; exception: `ProductInput` in `admin_product.go`.

## Pinia Stores

- **Composition API** (`auth.js`, `toast.js`): `defineStore('name', () => { ... })`
- **Options API** (`cart.js`, `favorites.js`, `useSettingsStore.js`): `defineStore('name', { state, getters, actions })`
- Store IDs: `'auth'`, `'cart'`, `'favorites'`, `'toast'`, `'settings'`.
- Options API stores define `getAuthHeaders()` at module level reading `localStorage.getItem('token')`.
- Auth store exposes `getAuthHeaders` as a returned method using the reactive `token` ref.
- Cart store uses 300ms debounce for quantity updates via `pendingUpdates` Map.

## GORM Models

- Each model defines base fields inline (no shared `Model` struct). `uint` for all IDs.
- JSON tags: camelCase. GORM tags: snake_case. Always include `CreatedAt`/`UpdatedAt`. Use `DeletedAt` for soft-delete.
- Response and input structs live alongside their models (e.g. `CreateOrderInput` in `models/order.go`).
- Preload: `database.DB.Preload("Category").Find(&products)`

## Test Patterns

### Frontend (Vitest)
- Test files colocated: `Component.spec.js` next to `Component.vue` in `src/`, or alongside stores
- Setup: `setActivePinia(createPinia())` + `localStorage.clear()` in `beforeEach`
- API calls mocked with `global.fetch = vi.fn().mockResolvedValue(...)`
- Component tests use `@vue/test-utils` `mount` with stubbed `vue-router` (`createMemoryHistory`)

### Backend (testify/suite)
- All tests in `server-go/test/` package, using `github.com/stretchr/testify/suite`
- Each domain has `XxxTestSuite` struct with `SetupSuite` (calls `SetupTestConfig` -> `SetupTestDB` -> `SetupRouter`)
- `MakeRequest(router, method, path, body, headers)` returns `*httptest.ResponseRecorder`
- Auth via dev bypass: `GetAdminAuthHeader()` / `GetCustomerAuthHeader()` returns `{"Authorization": "user-{id}"}`
- Uses separate `bycigar_test` database; `CleanDB()` truncates all tables between suite runs

## Key Architecture

- Admin endpoints return ALL records; public endpoints filter `WHERE is_active = true`. Never reuse public handlers on admin routes.
- Product listings sort `stock > 0` first via SQL CASE expression. Batch operations: `PUT /api/admin/products/batch/status`, `DELETE /api/admin/products/batch`.
- Orders use snowflake `OrderNo` (user-facing) + auto-increment `ID` (internal). `GetOrder` accepts both.
- Order status flow: `pending -> processing -> shipped -> completed`, with `cancelled` from pending/processing.
- Image upload: `POST /api/admin/upload` multipart `file` -> `{"success": true, "url": "/media/...", "thumbnailUrl": "/media/..."}`. Max 10MB, jpg/png/gif/webp.
- Server-side image processing: original resized to max 1200px wide, 300x300 JPEG thumbnail generated. GIF originals preserved.
- **Vite proxy**: `/api` -> `localhost:3000`, `/media` -> `localhost:9000` (strips `/media` prefix). Nginx mirrors this with proxy_cache (500MB, 30d).
- **App.vue layout**: TheHeader + TheFooter on all routes except `/admin/*`. Toast and CartDrawer always mounted.
- **Captcha**: Register always requires captcha. Login uses progressive captcha (3 failures -> required). Password change requires captcha.
- **Admin route groups**: `cmd/main.go` has two groups — `admin` (AdminOnly middleware) for shared routes, `superAdmin` (SuperAdminOnly) for restricted routes.

## Naming Conventions

| Type | Convention | Example |
|---|---|---|
| Vue Components | PascalCase | `ProductCard.vue` |
| Vue Test Files | PascalCase + .spec.js | `ProductCard.spec.js` |
| Views | PascalCase + View | `HomeView.vue` |
| Admin Views | Admin + Pascal | `AdminProducts.vue` |
| Stores | use + camel + Store | `useCartStore` |
| CSS Classes | kebab-case | `.product-card` |
| Go Handlers | PascalCase | `GetProducts` |
| Go Helpers | camelCase | `generateSlug` |
| Go Files | snake_case | `admin_product.go` |
| Go Test Suites | PascalCase + TestSuite | `ProductTestSuite` |

## Important Rules

- **Language**: Communicate in Chinese
- **No comments**: Do not add code comments in source files
- **No linter**: Use `go build` and `npm run build` to verify compilation
- **Run tests after changes**: `npm test` (frontend) and `go test ./... -v -count=1` (backend)
- **Simplicity**: Avoid over-engineering
- **Error handling**: Frontend uses toast, backend uses `utils.ErrorResponse()`
- **Database**: Always `utf8mb4`, `parseTime=True&loc=Local` in DSN
- **API paths**: Frontend always uses relative `/api/...` paths, never hardcoded localhost URLs

## Git Commits

```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
refactor: 简化购物车侧边栏逻辑
```
