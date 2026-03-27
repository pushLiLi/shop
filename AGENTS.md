# AGENTS.md - AI Coding Agent Guidelines

Guidelines for AI coding agents working on the BYCIGAR e-commerce codebase.

## Tech Stack

| Layer    | Technology                          |
|----------|-------------------------------------|
| Frontend | Vue 3 + Vite + Vue Router + Pinia   |
| Backend  | Go 1.23 + Gin + GORM + JWT          |
| Database | MySQL 8.4 (Docker)                  |
| API Docs | Swagger (swaggo)                    |

## Project Structure

```
shop/
├── bycigar-vue/
│   ├── src/components/     # Reusable components
│   ├── src/views/          # Page components + admin/
│   ├── src/stores/         # Pinia stores
│   └── vite.config.js      # Dev proxy: /api -> localhost:3000
├── server-go/
│   ├── cmd/main.go         # Entry point
│   ├── internal/handlers/  # HTTP handlers
│   ├── internal/models/    # GORM models + input structs
│   └── pkg/utils/          # ErrorResponse, SuccessResponse
└── bycigar_site/           # Legacy static site (reference only)
```

## Build Commands

```bash
# Frontend (bycigar-vue/)
npm run dev       # Dev server at http://localhost:5173
npm run build     # Production build -> dist/

# Backend (server-go/)
go run cmd/main.go              # Dev server at http://localhost:3000
go build -o server cmd/main.go  # Build executable
go fmt ./...                    # Format code

# Docker
docker-compose up -d mysql      # Start MySQL only
docker exec -it bycigar-mysql mysql -uroot -p123456 bycigar --default-character-set=utf8mb4

# Swagger docs (after adding new endpoints)
cd server-go && swag init
```

**Note**: No test framework configured yet.

## Code Style - Vue Components

```vue
<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useToastStore } from '../stores/toast'

const emit = defineEmits(['update', 'delete'])

const props = defineProps({
  product: { type: Object, required: true }
})

const router = useRouter()
const cartStore = useCartStore()
const toast = useToastStore()
const loading = ref(false)

async function handleSubmit() {
  try {
    loading.value = true
    await cartStore.addItem(props.product)
    toast.success('已添加到购物车')
  } catch (e) {
    console.error('Error:', e)
    toast.error('操作失败')
  } finally {
    loading.value = false
  }
}
</script>
```

### Vue Rules
- Use `<script setup>` (Composition API)
- Order: imports -> emit -> props -> router/stores -> refs -> computed -> functions
- All styles in `<style scoped>` blocks
- Use `useToastStore` for notifications, **never** use `alert()`
- API base: `http://localhost:3000/api`
- Auth header: `Authorization: Bearer ${token}`

## Code Style - Go Backend

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

func GetProduct(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
        return
    }

    var product models.Product
    if err := database.DB.First(&product, id).Error; err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
        return
    }

    c.JSON(http.StatusOK, product)
}
```

### Go Rules
- Import order: standard -> external -> internal
- Use `gofmt` for formatting
- Always use `utils.ErrorResponse()` for error responses
- Use `gin.H{}` for JSON responses
- Add Swagger comments for new endpoints
- Map frontend camelCase to DB snake_case (e.g., `categoryId` -> `category_id`)

## Pinia Store Pattern

```javascript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const API_BASE = 'http://localhost:3000/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

export const useCartStore = defineStore('cart', () => {
  const items = ref([])
  const total = computed(() => items.value.reduce(...))
  async function fetchCart() { }
  return { items, total, fetchCart }
})
```

## GORM Model Pattern

```go
type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"not null"`
    CategoryID  uint           `json:"categoryId"`
    Category    Category       `json:"category" gorm:"foreignKey:CategoryID"`
    CreatedAt   time.Time      `json:"createdAt"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
```

## Naming Conventions

| Type          | Convention           | Example                    |
|---------------|----------------------|----------------------------|
| Vue Components| PascalCase           | `ProductCard.vue`          |
| Views         | PascalCase + View    | `HomeView.vue`             |
| Stores        | camelCase + Store    | `useCartStore`             |
| Props         | camelCase            | `productId`                |
| CSS Classes   | kebab-case           | `product-card`             |
| API Endpoints | kebab-case           | `/api/cart-items`          |
| Go Handlers   | PascalCase           | `GetProducts`              |
| Go Models     | PascalCase           | `Product`                  |

## Authentication

- Token: `localStorage.getItem('token')`
- User: `localStorage.getItem('user')` (JSON)
- Protected routes: `meta: { requiresAuth: true }`
- Admin routes: `meta: { requiresAuth: true, requiresAdmin: true }`
- Default admin: `admin@admin.com` / `123456`

## API Endpoints

| Method | Endpoint              | Auth | Description               |
|--------|-----------------------|------|---------------------------|
| GET    | `/api/products`       | No   | List products (paginated) |
| GET    | `/api/products/:id`   | No   | Get product detail        |
| GET    | `/api/categories`     | No   | List categories           |
| GET    | `/api/cart`           | Yes  | Get user's cart           |
| POST   | `/api/cart`           | Yes  | Add item to cart          |
| POST   | `/api/auth/login`     | No   | Login                     |
| GET    | `/api/auth/me`        | Yes  | Get current user          |

Query params for `/api/products`: `page`, `limit`, `search`, `category`, `sortBy`, `sortOrder`

## Git Commit Convention

Use conventional commits with Chinese descriptions:
```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
refactor: 重构用户认证逻辑
```

## Important Notes

- **Language**: Communicate with user in Chinese
- **Simplicity**: Keep code clean, avoid over-engineering
- **Comments**: Only for complex logic explanations
- **Database charset**: Always use `utf8mb4` for Chinese characters
