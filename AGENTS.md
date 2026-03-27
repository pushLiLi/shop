# AGENTS.md - AI Coding Agent Guidelines

## Tech Stack

Frontend: Vue 3 + Vite + Vue Router + Pinia | Backend: Go 1.23 + Gin + GORM + JWT | Database: MySQL 8.4

## Build Commands

```bash
# Frontend (bycigar-vue/)
npm run dev       # Dev at http://localhost:5173
npm run build     # Production build

# Backend (server-go/)
go run cmd/main.go              # Dev at http://localhost:3000
go fmt ./...                    # Format code

# Docker & Swagger
docker-compose up -d mysql
cd server-go && swag init       # Generate API docs
```

**Note**: No lint or test framework configured.

## Project Structure

```
bycigar-vue/src/
├── components/   # ProductCard, CartDrawer, Toast
├── views/        # Page components + admin/
├── stores/       # cart, auth, toast, favorites
└── router/       # Vue Router

server-go/
├── cmd/main.go           # Entry point
├── internal/
│   ├── handlers/         # HTTP handlers
│   ├── models/           # GORM models
│   ├── middleware/       # Auth, CORS
│   ├── database/         # DB connection
│   └── config/           # App config
└── pkg/utils/            # Response helpers
```

## Vue Code Style

```vue
<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useToastStore } from '../stores/toast'

const emit = defineEmits(['update'])
const props = defineProps({ product: { type: Object, required: true } })

const router = useRouter()
const cartStore = useCartStore()
const toast = useToastStore()
const loading = ref(false)

async function handleSubmit() {
  try {
    loading.value = true
    await cartStore.addItem(props.product, 1)
    toast.success('已添加到购物车')
  } catch (e) {
    toast.error('操作失败')
  } finally {
    loading.value = false
  }
}
</script>
```

### Vue Rules

- `<script setup>` Composition API
- Order: imports → emit → props → stores → refs → computed → functions
- Use `useToastStore` for notifications, **never** `alert()`
- API base: `http://localhost:3000/api`
- CSS: kebab-case | No TypeScript

## Pinia Store Pattern

```javascript
const API_BASE = 'http://localhost:3000/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

export const useCartStore = defineStore('cart', {
  state: () => ({ items: [], loading: false }),
  
  getters: {
    total: (state) => state.items.reduce((sum, item) => 
      sum + (item.product?.price || 0) * item.quantity, 0)
  },
  
  actions: {
    async fetchCart() {
      try {
        this.loading = true
        const res = await fetch(`${API_BASE}/cart`, { headers: getAuthHeaders() })
        this.items = (await res.json()).items || []
      } finally {
        this.loading = false
      }
    }
  }
})
```

### Store Rules

- Options API style (`state`, `getters`, `actions`)
- Extract `getAuthHeaders()` at module level

## Go Code Style

```go
package handlers

import (
    "net/http"
    "strconv"
    "bycigar-server/internal/database"
    "bycigar-server/internal/models"
    "bycigar-server/pkg/utils"
    "github.com/gin-gonic/gin"
)

// GetProduct godoc
// @Summary 获取产品详情
// @Tags products
// @Param id path int true "产品ID"
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func GetProduct(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
        return
    }

    var product models.Product
    if err := database.DB.Preload("Category").First(&product, id).Error; err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
        return
    }

    c.JSON(http.StatusOK, product)
}
```

### Go Rules

- Import order: standard → external → internal
- Use `utils.ErrorResponse()` for errors
- Add Swagger comments for all handlers
- Early return pattern

## GORM Model

```go
type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"not null"`
    Price       float64        `json:"price"`
    CategoryID  uint           `json:"categoryId"`
    IsActive    bool           `json:"isActive" gorm:"default:true"`
    CreatedAt   time.Time      `json:"createdAt"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
```

- Use `uint` for IDs | Always include `CreatedAt` | JSON tags use camelCase

## Response Utilities

```go
utils.SuccessResponse(c, data)
utils.CreatedResponse(c, data)
utils.Success(c)
utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
utils.ErrorResponse(c, http.StatusNotFound, "Resource not found")
```

## Naming Conventions

| Type           | Convention        | Example           |
|----------------|-------------------|-------------------|
| Vue Components | PascalCase        | `ProductCard.vue` |
| Views          | PascalCase + View | `HomeView.vue`    |
| Stores         | camelCase + Store | `useCartStore`    |
| CSS Classes    | kebab-case        | `product-card`    |
| Go Handlers    | PascalCase        | `GetProducts`     |
| Go Functions   | camelCase         | `getAuthHeaders`  |

## Authentication

- Token: `localStorage.getItem('token')`
- User: `localStorage.getItem('user')`
- Protected routes: `meta: { requiresAuth: true }`
- Admin routes: `meta: { requiresAuth: true, requiresAdmin: true }`
- Default admin: `admin@admin.com` / `123456`
- JWT middleware sets `c.Get("userID")`

## API Endpoints

| Method | Endpoint              | Auth | Description        |
|--------|-----------------------|------|--------------------|
| GET    | `/api/products`       | No   | List products      |
| GET    | `/api/products/:id`   | No   | Product detail     |
| GET    | `/api/categories`     | No   | List categories    |
| POST   | `/api/auth/login`     | No   | Login              |
| GET    | `/api/auth/me`        | Yes  | Get profile        |
| GET    | `/api/cart`           | Yes  | Get cart           |
| POST   | `/api/cart`           | Yes  | Add to cart        |
| DELETE | `/api/cart/:id`       | Yes  | Remove from cart   |
| GET    | `/api/favorites`      | Yes  | List favorites     |
| POST   | `/api/favorites`      | Yes  | Add favorite       |
| DELETE | `/api/favorites/:id`  | Yes  | Remove favorite    |
| GET    | `/api/orders`         | Yes  | List orders        |
| POST   | `/api/orders`         | Yes  | Create order       |
| GET    | `/api/addresses`      | Yes  | List addresses     |
| POST   | `/api/addresses`      | Yes  | Create address     |

Query params: `page`, `limit`, `search`, `category`, `categoryId`, `sortBy`, `sortOrder`

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
