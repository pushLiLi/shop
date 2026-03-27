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

# Backend (server-go/)
go run cmd/main.go              # Dev at http://localhost:3000
go build -o server cmd/main.go  # Build executable
go fmt ./...                    # Format

# Docker
docker-compose up -d mysql
docker exec -it bycigar-mysql mysql -uroot -p123456 bycigar --default-character-set=utf8mb4

# Swagger
cd server-go && swag init
```

**Note**: No test framework configured yet.

## Project Structure

```
shop/
├── bycigar-vue/src/{components,views,stores}/
└── server-go/{cmd,internal/{handlers,models},pkg/utils}/
```

## Vue Code Style

```vue
<script setup>
import { ref } from 'vue'
import { useCartStore } from '../stores/cart'
import { useToastStore } from '../stores/toast'

const emit = defineEmits(['update'])
const props = defineProps({ product: { type: Object, required: true } })
const cartStore = useCartStore()
const toast = useToastStore()
const loading = ref(false)

async function handleSubmit() {
  try {
    loading.value = true
    await cartStore.addItem(props.product)
    toast.success('已添加到购物车')
  } catch (e) {
    toast.error('操作失败')
  } finally {
    loading.value = false
  }
}
</script>
```

- Use `<script setup>` Composition API
- Order: imports → emit → props → stores → refs → computed → functions
- Use `useToastStore` for notifications, **never** `alert()`
- API base: `http://localhost:3000/api`

## Go Code Style

```go
package handlers

import (
    "net/http"
    "bycigar-server/internal/database"
    "bycigar-server/internal/models"
    "bycigar-server/pkg/utils"
    "github.com/gin-gonic/gin"
)

func GetProduct(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var product models.Product
    if err := database.DB.First(&product, id).Error; err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
        return
    }
    c.JSON(http.StatusOK, product)
}
```

- Import order: standard → external → internal
- Use `utils.ErrorResponse()` for all errors
- Map camelCase to snake_case (`categoryId` → `category_id`)

## Pinia Store

```javascript
const API_BASE = 'http://localhost:3000/api'
function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return { 'Authorization': token ? `Bearer ${token}` : '' }
}

export const useCartStore = defineStore('cart', () => {
  const items = ref([])
  const total = computed(() => items.value.reduce(...))
  return { items, total, fetchCart }
})
```

## GORM Model

```go
type Product struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"not null"`
    CategoryID uint          `json:"categoryId"`
    CreatedAt time.Time      `json:"createdAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
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

## Authentication

- Token: `localStorage.getItem('token')`
- Protected routes: `meta: { requiresAuth: true }`
- Admin routes: `meta: { requiresAuth: true, requiresAdmin: true }`
- Default admin: `admin@admin.com` / `123456`

## API Endpoints

| Method | Endpoint          | Auth | Description   |
|--------|-------------------|------|---------------|
| GET    | `/api/products`   | No   | List products |
| GET    | `/api/cart`       | Yes  | Get cart      |
| POST   | `/api/cart`       | Yes  | Add to cart   |
| POST   | `/api/auth/login` | No   | Login         |

Query params: `page`, `limit`, `search`, `category`, `sortBy`, `sortOrder`

## Git Commits

```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
```

## Important

- **Language**: 中文沟通
- **Simplicity**: 避免过度设计
- **Database**: 使用 `utf8mb4`
