# AGENTS.md - AI Coding Agent Guidelines

## Tech Stack

Frontend: Vue 3 + Vite + Vue Router + Pinia + marked
Backend: Go 1.23 + Gin + GORM + JWT + Swagger
Database: MySQL 8.4 (Docker)

## Build Commands

```bash
# Frontend (bycigar-vue/)
npm run dev       # Dev at http://localhost:5173
npm run build     # Production build

# Backend (server-go/)
go run cmd/main.go              # Dev at http://localhost:3000
go fmt ./...                    # Format code
go build -o /dev/null ./...     # Check compilation

# Docker & Swagger
docker-compose up -d mysql      # Start MySQL only
cd server-go && swag init       # Generate API docs at /swagger/index.html
```

**Note**: No lint or test framework configured.

## Project Structure

```
bycigar-vue/src/
├── components/     # ProductCard, CartDrawer, Toast, TheFooter
├── views/          # Page components + admin/
├── stores/         # cart, auth, toast, favorites
└── router/         # Vue Router

server-go/
├── cmd/main.go             # Entry point
├── internal/
│   ├── handlers/           # HTTP handlers (public + admin_*)
│   ├── models/             # GORM models
│   ├── middleware/         # Auth, CORS, AdminOnly
│   ├── database/           # DB connection + migrations + seeds
│   └── config/             # App config from .env
└── pkg/utils/              # Response helpers
```

## Vue Rules

- `<script setup>` Composition API only
- Order: imports → emit → props → stores → refs → computed → functions
- Use `useToastStore` for notifications, **never** `alert()`
- API base: `http://localhost:3000/api`
- CSS: scoped + kebab-case class names
- No TypeScript

## Pinia Store Rules

- Options API style (`state`, `getters`, `actions`)
- Extract `getAuthHeaders()` at module level:
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

## Go Rules

- Import order: standard → external → internal
- Use `utils.ErrorResponse()` for all error responses
- Add Swagger comments (`godoc`) for all handlers
- Early return pattern for error handling
- Admin handlers prefixed with `admin_*.go`

## GORM Model

- Use `uint` for IDs | Always include `CreatedAt` | JSON tags use camelCase | Soft delete with `DeletedAt`

## Response Utilities

```go
utils.SuccessResponse(c, data)                              // 200 OK with data
utils.CreatedResponse(c, data)                              // 201 Created
utils.Success(c)                                            // 200 { success: true }
utils.ErrorResponse(c, http.StatusBadRequest, "Invalid")    // 4xx with error
```

## Naming Conventions

| Type            | Convention         | Example             |
|-----------------|--------------------|---------------------|
| Vue Components  | PascalCase         | `ProductCard.vue`   |
| Views           | PascalCase + View  | `HomeView.vue`      |
| Admin Views     | Admin + Pascal     | `AdminProducts.vue` |
| Stores          | use + camel + Store| `useCartStore`      |
| CSS Classes     | kebab-case         | `.product-card`     |
| Go Handlers     | PascalCase         | `GetProducts`       |
| Go Functions    | camelCase          | `getAuthHeaders`    |

## Authentication

- Token: `localStorage.getItem('token')`
- User: `localStorage.getItem('user')` (JSON with `role` field)
- Protected routes: `meta: { requiresAuth: true }`
- Admin routes: `meta: { requiresAuth: true, requiresAdmin: true }`
- Default admin: `admin@admin.com` / `123456`
- JWT middleware sets `c.Get("userID")` in context

## API Endpoints

### Public: `/api/products`, `/api/products/:id`, `/api/categories`, `/api/banners`, `/api/pages/:slug`, `/api/auth/login`, `/api/auth/register`

### Authenticated: `/api/auth/me`, `/api/cart`, `/api/favorites`, `/api/orders`, `/api/addresses`

### Admin: `/api/admin/products`, `/api/admin/categories`, `/api/admin/banners`, `/api/admin/pages`, `/api/admin/upload`

Query params: `page`, `limit`, `search`, `category`, `categoryId`, `sortBy`, `sortOrder`

## CMS Pages

Slugs: `about`, `services`, `privacy-policy`, `statement`
Frontend uses `marked` library to render Markdown content.

## Git Commits

```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
refactor: 简化购物车侧边栏逻辑
```

## Important

- **Language**: 中文沟通
- **Simplicity**: 避免过度设计
- **Database**: `utf8mb4` 字符集
- **Error handling**: 前端 toast，后端 `utils.ErrorResponse()`
- **No comments**: 不添加代码注释
- **Concise responses**: 回复不超过4行
