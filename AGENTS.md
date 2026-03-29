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
# Frontend tests (run from bycigar-vue/) — Vitest 4 + @vue/test-utils + happy-dom
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
Vitest config has `globals: true` — no need to import `describe`/`it`/`expect`.
Frontend tests are unit tests (no backend needed). Backend tests are integration tests hitting real MySQL + MinIO.

## Roles & Authorization

Three user roles in `User.Role` field: `"admin"` (超级管理员), `"service"` (管理员), `"customer"` (客户). Role priority: admin > service > customer.

### Backend Middleware
- **`AdminOnly`**: Allows `admin` + `service` — shared admin routes (products, categories, orders, dashboard, users, upload)
- **`SuperAdminOnly`**: Allows `admin` only — restricted routes (banners, pages, config, settings, user role changes)
- **`RequireAuth`**: Requires any logged-in user — used for cart, favorites, addresses, orders, notifications
- **Role-based action guards**: `ResetUserPassword` checks current user role — `service` cannot reset `admin` password

### Frontend Auth Store (`useAuthStore`)
- `isAdmin`: `role === 'admin' || role === 'service'` — can access admin panel
- `isSuperAdmin`: `role === 'admin'` — can see revenue data, banners, pages, settings menus
- `isService`: `role === 'service'` — for conditional UI in shared views

### Frontend Route Guards
- `meta.requiresAuth`: requires logged-in user (cart, checkout, orders, favorites, profile, notifications)
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
| Reset user password | ✅ | ✅ (not admin) | ❌ |
| Change user role | ✅ | ❌ | ❌ |
| Banners / Pages / Config / Settings | ✅ | ❌ | ❌ |

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
- Admin handlers in separate `admin_*.go` files. SuperAdmin handlers (banners, pages, config, settings) live in their domain files (`banner.go`, `page.go`, etc.) alongside public handlers.
- Swagger `godoc` comments on all handlers. No other code comments.
- Error responses: `utils.ErrorResponse(c, statusCode, "msg")` or `c.JSON(status, gin.H{"error": "..."})`.
- Input structs live in model files; exception: `ProductInput` in `admin_product.go`.
- **Notification creation**: When admin actions affect users (order status change, product stock/price change), create `models.Notification` records. Use batch `database.DB.Create(&notifications)` for multiple recipients.

## Pinia Stores

- **Composition API** (`auth.js`, `toast.js`): `defineStore('name', () => { ... })`
- **Options API** (`cart.js`, `favorites.js`, `notifications.js`, `useSettingsStore.js`): `defineStore('name', { state, getters, actions })`
- Store IDs: `'auth'`, `'cart'`, `'favorites'`, `'toast'`, `'settings'`, `'notifications'`.
- Options API stores define `getAuthHeaders()` at module level reading `localStorage.getItem('token')`.
- Auth store exposes `getAuthHeaders` as a returned method using the reactive `token` ref.
- Cart store uses 300ms debounce for quantity updates via `pendingUpdates` Map.

## GORM Models

- Each model defines base fields inline (no shared `Model` struct). `uint` for all IDs.
- JSON tags: camelCase. GORM tags: snake_case. Always include `CreatedAt`/`UpdatedAt`. Exception: `Notification` has only `CreatedAt`.
- Use `DeletedAt` for soft-delete (products, categories). Response and input structs live alongside their models.
- Notification model uses `*uint` (pointer) for optional foreign keys (`ProductID`, `OrderID`).
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
- Test users: admin@test.com (admin), user1@test.com (customer), user2@test.com (customer)

## Key Architecture

- Admin endpoints return ALL records; public endpoints filter `WHERE is_active = true`. Never reuse public handlers on admin routes.
- Product listings sort `stock > 0` first via SQL CASE expression. Batch operations: `PUT /api/admin/products/batch/status`, `DELETE /api/admin/products/batch`.
- Orders use snowflake `OrderNo` (user-facing) + auto-increment `ID` (internal). `GetOrder` accepts both.
- Order status flow: `pending -> processing -> shipped -> completed`, with `cancelled` from pending/processing.
- **Notifications**: System-generated, read-only. Types: `order_status`, `back_in_stock`, `price_drop`. Triggered in admin handlers. Frontend: bell icon in header + right-side panel. Clicking a notification navigates to `/notifications/:id` detail page (auto marks read), which has a link to the related order or product.
- Image upload: `POST /api/admin/upload` multipart `file` -> `{"success": true, "url": "/media/...", "thumbnailUrl": "/media/..."}`. Max 10MB, jpg/png/gif/webp.
- **Vite proxy**: `/api` -> `localhost:3000`, `/media` -> `localhost:9000` (strips `/media` prefix). Nginx mirrors this with proxy_cache (500MB, 30d).
- **App.vue layout**: TheHeader + TheFooter on all routes except `/admin/*`. Toast, CartDrawer always mounted. NotificationPanel inside TheHeader via Teleport.
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
| Notification Types | snake_case constant | `NotificationTypeOrderStatus` |

## Important Rules

- **Language**: Communicate in Chinese
- **No comments**: Do not add code comments in source files
- **No linter**: Use `go build` and `npm run build` to verify compilation
- **Run tests after changes**: `npm test` (frontend) and `go test ./... -v -count=1` (backend)
- **Simplicity**: Avoid over-engineering
- **Error handling**: Frontend uses toast, backend uses `utils.ErrorResponse()`
- **Database**: Always `utf8mb4`, `parseTime=True&loc=Local` in DSN
- **API paths**: Frontend always uses relative `/api/...` paths, never hardcoded localhost URLs
- **New models**: Must be added to `database.Migrate()` in `database/database.go`

## Git Commits

```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
refactor: 简化购物车侧边栏逻辑
```
