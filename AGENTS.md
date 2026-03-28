# AGENTS.md - AI Coding Agent Guidelines

## Tech Stack

Frontend: Vue 3 + Vite + Vue Router + Pinia + marked (JavaScript only, no TypeScript)
Backend: Go 1.25 + Gin + GORM + JWT + Swagger
Database: MySQL 8.4 (Docker, utf8mb4 charset)
Object Storage: MinIO (Docker, API:9000, Console:9001)
Module: `bycigar-server` (Go module name)

## Build Commands

```bash
# Frontend (bycigar-vue/)
npm run dev       # Dev at http://localhost:5173 (proxies /api -> localhost:3000, /media -> localhost:9000)
npm run build     # Production build to dist/

# Backend (server-go/)
go run ./cmd/main.go              # Dev at http://localhost:3000
go fmt ./...                      # Format code
go build -o /dev/null ./...       # Check compilation (use -o nul on Windows)

# Docker & Swagger
docker-compose up -d mysql        # Start MySQL only (root/123456, database: bycigar)
docker-compose up -d minio        # Start MinIO only (minioadmin/minioadmin123, bucket: bycigar)
swag init                         # Generate Swagger docs (run from server-go/)

# Full stack
docker-compose up --build         # MySQL:3306 + MinIO:9000 + Backend:3000 + Frontend:80
```

**No test or lint framework is configured.** Use `go build` and `npm run build` to verify compilation.

## Project Structure

```
bycigar-vue/src/
├── main.js, App.vue        # Entry + layout (dark theme, Toast, CartDrawer)
├── style.css               # Global dark theme, grid (.col-2/3/4/6/12), 768px breakpoint
├── components/              # ProductCard, CartDrawer, Toast, AdminImageUpload, TheHeader/TheFooter
├── composables/             # Reusable logic: useCarousel (carousel with autoplay, touch, pause)
├── utils/                   # Pure helpers: states.js (US state lookup)
├── views/                   # Page views + admin/ (light theme, accent #d4a574)
├── stores/                  # Pinia: auth, cart, favorites, toast, useSettingsStore
└── router/index.js          # Routes with auth/admin guards

server-go/
├── cmd/main.go              # Entry: config -> DB -> migrate -> seed -> MinIO -> snowflake -> Gin
├── internal/
│   ├── config/config.go     # Loads .env via godotenv into AppConfig global
│   ├── database/database.go # Connect, Migrate, Seed (admin, pages, settings, siteconfig), BackfillOrderNo
│   ├── handlers/            # Public handlers + admin_*.go for admin endpoints
│   ├── models/              # GORM models with json/gorm struct tags
│   └── middleware/          # CORS, AuthMiddleware (optional), RequireAuth, AdminOnly
├── pkg/
│   ├── minio/minio.go       # MinIO client wrapper (InitMinio, EnsureBucket, exposes Client+Bucket)
│   └── utils/               # Response helpers + snowflake order number generation
```

## Vue Code Style

- `<script setup>` Composition API only
- Declaration order: imports → `defineEmits` → `defineProps` → stores → refs/computed → functions
- Use `useToastStore` for notifications, **never** `alert()`
- Composables go in `src/composables/`, named `use` + PascalCase, return reactive state and handlers
- API base: `http://localhost:3000/api` (hardcoded in stores, Vite proxy also configured)
- CSS: `<style scoped>` with kebab-case class names
- No TypeScript, no code comments
- Inline event handlers use named `async function` declarations, not arrow functions
- When sending JSON to backend, **match the Go model's JSON tags exactly** (e.g. `"imageUrl"` not `"image"`)

## Pinia Store Patterns

Two styles exist in the codebase:
- **Composition API** (`auth.js`, `toast.js`): `defineStore('name', () => { ... })`
- **Options API** (`cart.js`, `favorites.js`, `useSettingsStore.js`): `defineStore('name', { state, getters, actions })`

Both define `getAuthHeaders()` at module level:
```javascript
const API_BASE = 'http://localhost:3000/api'
function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}
```

## Go Code Style

- Import order: standard library → blank line → external → blank line → internal
- Early return pattern for all error handling
- Handler functions are PascalCase (`GetProducts`, `CreateProduct`)
- Helper/utility functions are camelCase (`generateSlug`, `getEnv`)
- Admin handlers live in separate `admin_*.go` files with separate handler functions
- Add Swagger `godoc` comments on all handler functions
- No code comments except Swagger annotations
- Use `utils.ErrorResponse(c, statusCode, "message")` for error responses
- Some handlers use `c.JSON()` directly with `gin.H{}` for inline responses

## GORM Models

```go
type Model struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

- `uint` for all IDs | JSON tags use camelCase | GORM tags use snake_case column names
- Always include `CreatedAt` and `UpdatedAt` | Use `DeletedAt` for soft-delete models
- Input structs (`ProductInput`, `CreateOrderInput`) defined in handler files, not model files
- Preload associations: `database.DB.Preload("Category").Find(&products)`

## Response Utilities (pkg/utils/response.go)

```go
utils.SuccessResponse(c, data)                           // 200 + data
utils.CreatedResponse(c, data)                           // 201 + data
utils.Success(c)                                         // 200 {"success": true}
utils.ErrorResponse(c, http.StatusBadRequest, "Invalid") // 4xx {"error": "msg"}
```

## Naming Conventions

| Type            | Convention          | Example              |
|-----------------|---------------------|----------------------|
| Vue Components  | PascalCase          | `ProductCard.vue`    |
| Views           | PascalCase + View   | `HomeView.vue`       |
| Admin Views     | Admin + Pascal      | `AdminProducts.vue`  |
| Stores          | use + camel + Store | `useCartStore`       |
| Composables     | use + PascalCase    | `useCarousel`        |
| CSS Classes     | kebab-case          | `.product-card`      |
| Go Handlers     | PascalCase          | `GetProducts`        |
| Go Helpers      | camelCase           | `generateSlug`       |
| Go Files        | snake_case          | `admin_product.go`   |

## Authentication Flow

- Token stored in `localStorage.getItem('token')`
- User JSON in `localStorage.getItem('user')` (contains `role` field)
- `AuthMiddleware()`: optional JWT parsing, sets `c.Set("userID", ...)` in Gin context
- `RequireAuth()`: wraps routes requiring any authenticated user (returns 401)
- `AdminOnly`: standalone middleware, checks `userID` + DB lookup + `role === "admin"`, sets `c.Set("user", user)`
- Dev bypass: `Authorization: user-{id}` header skips JWT validation
- Default admin: `admin@admin.com` / `123456`

## API Endpoints

Public: `/api/products`, `/api/products/:id`, `/api/categories`, `/api/banners`, `/api/pages/:slug`, `/api/auth/login`, `/api/auth/register`, `/api/config`, `/api/settings`

Authenticated: `/api/auth/me`, `/api/auth/profile`, `/api/auth/change-password`, `/api/cart`, `/api/favorites`, `/api/orders`, `/api/addresses`

Admin: `/api/admin/products`, `/api/admin/categories`, `/api/admin/banners`, `/api/admin/pages`, `/api/admin/upload`, `/api/admin/config/:key`, `/api/admin/settings/:key`

Query params: `page`, `limit`, `search`, `category` (slug), `categoryId`, `sortBy`, `sortOrder`, `featured`, `active`

## Vite Proxy Configuration

`vite.config.js` proxies `/api` to `http://localhost:3000` and `/media` to `http://localhost:9000` (MinIO). The `/media` proxy serves uploaded images during development. Both Go backend and MinIO must be running alongside Vite dev server.

## Architecture Notes

- **Admin vs Public endpoints**: Admin list endpoints return ALL records (no `is_active` filter); public endpoints filter `WHERE is_active = true`. Never reuse public handlers on admin routes.
- **Stock sorting**: Public product listings always sort `stock > 0` items first via SQL CASE expression, regardless of user-selected sort.
- **Order numbers**: Snowflake-generated `OrderNo` (user-facing) separate from auto-increment `ID` (internal). `GetOrder` accepts both.
- **Banner JSON tag**: Go field `Image` has JSON tag `"imageUrl"`. Frontend must send `"imageUrl"` in request bodies and reads `"imageUrl"` from responses.
- **Config API**: `GET /api/config` returns flat `map[string]string`. Update via `PUT /api/admin/config/:key` with body `{"value": "..."}`.
- **Settings API**: `GET /api/settings` returns `{"success": true, "data": {...}}`. Update via `PUT /api/admin/settings/:key`.
- **Image upload**: `POST /api/admin/upload` returns `{"url": "/media/bycigar/xxx.jpg", "success": true}`. Images stored in MinIO, served via nginx/Vite reverse proxy `/media/` → MinIO. Go backend handles writes only; reads go directly to MinIO through the proxy.
- **CMS pages**: Slugs `about`, `services`, `privacy-policy`, `statement`. Content rendered as Markdown via `marked` library.
- **Auto-migration**: All tables auto-migrate on backend startup. Admin user and default data seeded on first run.

## Git Commits

```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
refactor: 简化购物车侧边栏逻辑
```

## Important Rules

- **Language**: Communicate in Chinese (中文沟通)
- **No comments**: Do not add code comments in source files
- **No tests**: No test framework configured, use build to verify correctness
- **Simplicity**: Avoid over-engineering, keep implementations simple
- **Error handling**: Frontend uses toast notifications, backend uses `utils.ErrorResponse()`
- **Database**: Always use `utf8mb4` charset, `parseTime=True&loc=Local` in DSN
- **Concise responses**: Keep responses under 4 lines
