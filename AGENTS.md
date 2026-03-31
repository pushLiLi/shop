# AGENTS.md - AI Coding Agent Guidelines

## Tech Stack

Frontend: Vue 3.5 + Vite 8 + Vue Router 4 + Pinia 3 + marked + vue-advanced-cropper + vue-chartjs + chart.js (JavaScript only, no TypeScript, no axios — uses native `fetch`)
Backend: Go 1.23 + Gin 1.9 + GORM 1.25 + JWT + gorilla/websocket + Swagger (Dockerfile: golang:1.23-alpine)
Database: MySQL 8.4 (Docker, utf8mb4) | Object Storage: MinIO (Docker, API:9000, Console:9001)
Module: `bycigar-server`

## Build Commands

```bash
# Frontend (run from bycigar-vue/)
npm run dev                     # Dev server at :5173 (proxies /api -> :3000, /media -> :9000)
npm run build                   # Production build (verifies frontend changes)

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
# Frontend unit tests (run from bycigar-vue/) — Vitest 4 + @vue/test-utils + happy-dom
npm test                        # Run all tests once
npx vitest run src/stores/cart.spec.js          # Run single test file
npx vitest run --reporter=verbose               # Verbose output

# Backend tests (run from server-go/) — Go testing + testify/suite
# Requires MySQL + MinIO running (docker-compose up -d mysql minio)
go test ./test/... -v -count=1                  # All integration tests
go test ./test/... -v -count=1 -run TestProductSuite  # Single test suite
go test ./test/... -v -count=1 -run "TestProductSuite/TestGetProductsDefaultPagination"  # Single test

# E2E tests (run from bycigar-vue/) — Playwright
npx playwright test                    # All e2e tests
npx playwright test --project=chromium # Specific browser
```
Vitest config has `globals: true` — no need to import `describe`/`it`/`expect`. Frontend tests are isolated; backend tests hit real MySQL + MinIO.

## Roles & Authorization

Three roles: `"admin"` (超级管理员), `"service"` (管理员), `"customer"` (客户). Priority: admin > service > customer.

| Middleware | Access |
|---|---|
| `AdminOnly` | admin + service — products, categories, orders, dashboard, users, upload, chat admin |
| `SuperAdminOnly` | admin only — banners, pages, config, settings, user role changes |
| `RequireAuth` | any logged-in user — cart, favorites, addresses, orders, notifications, chat |

Frontend auth: `isAdmin` = `role === 'admin' || role === 'service'`; `isSuperAdmin` = `role === 'admin'`; `isService` = `role === 'service'`.

## Code Style

### Go

**Imports** (strict order, blank lines between groups):
```
stdlib
(blank line)
external (github.com, etc.)
(blank line)
internal (bycigar-server/...)
```
**Formatting**: Early return for all error handling. Handlers PascalCase, helpers camelCase. Admin handlers in `admin_*.go` files. SuperAdmin handlers in domain files alongside public handlers.
**Comments**: Swagger `godoc` on all handlers only. No other comments.
**Errors**: `utils.ErrorResponse(c, statusCode, "msg")` or `c.JSON(status, gin.H{"error": "..."})`.
**Models**: Define base fields inline (no shared Model struct). Use `uint` for IDs. JSON tags camelCase, GORM tags snake_case. Always include `CreatedAt`/`UpdatedAt` (except Notification, Message which have only `CreatedAt`). Use `DeletedAt` for soft-delete. Input structs in model files (exception: `ProductInput` in `admin_product.go`).
**Notifications**: Create `models.Notification` records via `database.DB.Create(&notifications)` when admin actions affect users.

### Vue

**Imports**: Relative paths only (no `@/` alias). Declaration order: imports -> `defineEmits` -> `defineProps` -> stores -> refs/computed -> functions.
**API calls**: Use relative `/api/...` paths, never hardcoded URLs. Match Go model JSON tags exactly (e.g., `"imageUrl"` not `"image"`, `"isFeatured"` not `"featured"`).
**Auth guard**: Before cart/favorites actions, check `authStore.isLoggedIn`. If false, `router.push('/login')`.
**Errors**: Use `useToastStore`, never `alert()`.
**Images**: List views: `product.thumbnailUrl || product.imageUrl` with `loading="lazy"`. Detail pages use full `imageUrl`.
**Inline handlers**: Use named `async function` declarations, not arrow functions.
**CSS**: `<style scoped>` with kebab-case classes. Storefront dark theme (#0f0f0f bg, #d4a574 accent). Admin light theme (#fff). Icons: inline SVG only.

### Pinia Stores

- **Composition API** (`auth.js`, `toast.js`): `defineStore('name', () => { ... })`
- **Options API** (`cart.js`, `favorites.js`, `notifications.js`, `chat.js`, `useSettingsStore.js`, `contactMethods.js`): `defineStore('name', { state, getters, actions })`
- Store IDs: `'auth'`, `'cart'`, `'favorites'`, `'toast'`, `'settings'`, `'notifications'`, `'chat'`, `'contactMethods'`.
- Options API stores define `getAuthHeaders()` at module level.
- Cart store uses 300ms debounce for quantity updates via `pendingUpdates` Map.

## Test Patterns

### Frontend (Vitest)
- Test files colocated: `Component.spec.js` next to `Component.vue`
- Setup: `setActivePinia(createPinia())` + `localStorage.clear()` in `beforeEach`
- Mock API: `global.fetch = vi.fn().mockResolvedValue(...)`
- Components: `@vue/test-utils` `mount` with stubbed `vue-router` (`createMemoryHistory`)

### Backend (testify/suite)
- Tests in `server-go/test/` package using `github.com/stretchr/testify/suite`
- Each domain has `XxxTestSuite` with `SetupSuite` (calls `SetupTestConfig` -> `SetupTestDB` -> `SetupRouter`)
- `MakeRequest(router, method, path, body, headers)` returns `*httptest.ResponseRecorder`
- `MakeFormRequest(...)` for multipart uploads; `ParseResponse(w)` unmarshals to `map[string]interface{}`
- Auth bypass: `GetAdminAuthHeader()` / `GetCustomerAuthHeader()` returns `{"Authorization": "user-{id}"}`
- Uses `bycigar_test` database; `CleanDB()` truncates all tables between runs
- Test users: admin@test.com (admin), user1@test.com (customer), user2@test.com (customer)

## Key Architecture

- Admin endpoints return ALL records; public endpoints filter `WHERE is_active = true`.
- Product listings sort `stock > 0` first via SQL CASE expression. Batch ops: `PUT /api/admin/products/batch/status`, `DELETE /api/admin/products/batch`.
- Orders: snowflake `OrderNo` (user-facing) + auto-increment `ID` (internal). `GetOrder` accepts both.
- Order flow: `pending -> {processing, cancelled}` -> `{shipped, cancelled}` -> `shipped -> completed`. Transitions enforced by `ValidOrderStatusTransitions` map in `models/order.go`.
- Notifications: system-generated, read-only. Types: `order_status`, `back_in_stock`, `price_drop`.
- Chat: WebSocket (gorilla/websocket) with HTTP fallback. One conversation per customer. 500 char limit. 30-day auto-cleanup via `pkgutils.StartChatCleanup()`.
- Upload: `POST /api/admin/upload` multipart `file` -> `{"success": true, "url": "/media/...", "thumbnailUrl": "/media/..."}`. Max 10MB, jpg/png/gif/webp.
- Captcha: Register always requires captcha. Login progressive (3 failures -> required). Password change requires captcha.
- New models: add to `database.Migrate()` in `database/database.go`.
- New routes: add to both `cmd/main.go` and `test/helpers.go` `SetupRouter()`.

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
- **Run tests after changes**: `npm test` (frontend) and `go test ./... -v -count=1` (backend)
- **Error handling**: Frontend uses toast, backend uses `utils.ErrorResponse()`
- **Database**: Always `utf8mb4`, `parseTime=True&loc=Local` in DSN

## Git Commits

```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
refactor: 简化购物车侧边栏逻辑
```
