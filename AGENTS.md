# AGENTS.md - AI Coding Agent Guidelines

BYCIGAR e-commerce codebase guidelines.

## Tech Stack

| Layer    | Technology                        |
|----------|-----------------------------------|
| Frontend | Vue 3 + Vite + Vue Router + Pinia |
| Backend  | Go 1.23 + Gin + GORM + JWT        |
| Database | MySQL 8.4 (Docker)                |
| API Docs | Swagger (swaggo)                  |

## Build Commands

```bash
# Frontend (bycigar-vue/)
npm run dev       # Dev server at http://localhost:5173
npm run build     # Production build
npm run preview   # Preview production build

# Backend (server-go/)
go run cmd/main.go              # Dev at http://localhost:3000
go build -o server cmd/main.go  # Build executable
go fmt ./...                    # Format code

# Docker
docker-compose up -d mysql
docker exec -it bycigar-mysql mysql -uroot -p123456 bycigar --default-character-set=utf8mb4

# Swagger
cd server-go && swag init       # Generate API docs
```

**Note**: No test framework configured yet.

## Project Structure

```
shop/
├── bycigar-vue/src/
│   ├── components/      # Reusable components (ProductCard, CartDrawer, Toast)
│   ├── views/           # Page components (HomeView, CartView, CheckoutView)
│   ├── views/admin/     # Admin panel pages
│   ├── stores/          # Pinia stores (cart, auth, toast, favorites)
│   └── router/          # Vue Router configuration
└── server-go/
    ├── cmd/main.go              # Entry point
    ├── internal/
    │   ├── handlers/            # HTTP handlers
    │   ├── models/              # GORM models
    │   ├── middleware/          # Auth, CORS middleware
    │   ├── database/            # DB connection
    │   └── config/              # App configuration
    └── pkg/utils/               # Response helpers
```

## Vue Code Style

```vue
<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useToastStore } from '../stores/toast'

const emit = defineEmits(['update', 'remove'])
const props = defineProps({
  product: { type: Object, required: true },
  quantity: { type: Number, default: 1 }
})

const router = useRouter()
const cartStore = useCartStore()
const toast = useToastStore()
const loading = ref(false)

const totalPrice = computed(() => props.product.price * props.quantity)

async function handleSubmit() {
  try {
    loading.value = true
    await cartStore.addItem(props.product, props.quantity)
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

- Use `<script setup>` Composition API
- Order: imports → emit → props → stores → refs → computed → functions
- Use `useToastStore` for notifications, **never** `alert()`
- API base: `http://localhost:3000/api`
- Use kebab-case for CSS classes

## Pinia Store Pattern

```javascript
import { defineStore } from 'pinia'
import { useToastStore } from './toast'

const API_BASE = 'http://localhost:3000/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [],
    loading: false
  }),
  
  getters: {
    total: (state) => state.items.reduce((sum, item) => 
      sum + (item.product?.price || 0) * item.quantity, 0
    ),
    count: (state) => state.items.reduce((sum, item) => sum + item.quantity, 0)
  },
  
  actions: {
    async fetchCart() {
      try {
        this.loading = true
        const res = await fetch(`${API_BASE}/cart`, { headers: getAuthHeaders() })
        const data = await res.json()
        this.items = data.items || []
      } catch (e) {
        console.error('Fetch cart failed:', e)
      } finally {
        this.loading = false
      }
    }
  }
})
```

### Store Rules

- Use Options API style (`state`, `getters`, `actions`)
- Handle errors with try/catch and `useToastStore().error()`
- For optimistic updates: update local state first, sync API in background

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
// @Description 根据ID获取产品详情
// @Tags products
// @Param id path int true "产品ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]interface{}
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
- Use `utils.ErrorResponse()` for all error responses
- Add Swagger comments for all handlers
- Use `Preload()` for eager loading relations
- Map camelCase (JSON) to snake_case (DB): `categoryId` → `category_id`

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

## Naming Conventions

| Type           | Convention        | Example           |
|----------------|-------------------|-------------------|
| Vue Components | PascalCase        | `ProductCard.vue` |
| Views          | PascalCase + View | `HomeView.vue`    |
| Stores         | camelCase + Store | `useCartStore`    |
| CSS Classes    | kebab-case        | `product-card`    |
| Go Handlers    | PascalCase        | `GetProducts`     |
| Go Models      | PascalCase        | `CartItem`        |

## Authentication

- Token stored in: `localStorage.getItem('token')`
- User info in: `localStorage.getItem('user')`
- Protected routes: `meta: { requiresAuth: true }`
- Admin routes: `meta: { requiresAuth: true, requiresAdmin: true }`
- Default admin: `admin@admin.com` / `123456`

## API Endpoints

| Method | Endpoint              | Auth | Description        |
|--------|-----------------------|------|--------------------|
| GET    | `/api/products`       | No   | List products      |
| GET    | `/api/products/:id`   | No   | Product detail     |
| GET    | `/api/categories`     | No   | List categories    |
| POST   | `/api/auth/login`     | No   | Login              |
| POST   | `/api/auth/register`  | No   | Register           |
| GET    | `/api/auth/me`        | Yes  | Get profile        |
| GET    | `/api/cart`           | Yes  | Get cart           |
| POST   | `/api/cart`           | Yes  | Add to cart        |
| PUT    | `/api/cart/:id`       | Yes  | Update quantity    |
| DELETE | `/api/cart/:id`       | Yes  | Remove from cart   |
| GET    | `/api/orders`         | Yes  | List orders        |
| POST   | `/api/orders`         | Yes  | Create order       |
| GET    | `/api/favorites`      | Yes  | List favorites     |
| POST   | `/api/favorites`      | Yes  | Add favorite       |

Query params: `page`, `limit`, `search`, `category`, `categoryId`, `sortBy`, `sortOrder`

## Git Commits

```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
refactor: 简化购物车侧边栏逻辑
docs: 更新 AGENTS.md 开发指南
```

## Important

- **Language**: 中文沟通
- **Simplicity**: 避免过度设计
- **Database**: 使用 `utf8mb4` 字符集
- **Error handling**: 前端用 toast，后端用 `utils.ErrorResponse()`
