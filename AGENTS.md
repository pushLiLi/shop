# AGENTS.md - AI Coding Agent Guidelines

Guidelines for AI coding agents (Claude, GPT-4, Copilot, Cursor, etc.) working on this codebase.

## Project Overview

BYCIGAR e-commerce platform with Vue 3 frontend and Go backend.

| Layer       | Technology                                    |
|-------------|-----------------------------------------------|
| Frontend    | Vue 3 + Vite + Vue Router + Pinia             |
| Backend     | Go 1.23 + Gin + GORM + JWT                    |
| Database    | MySQL 8.4 (Docker)                            |
| API Docs    | Swagger (swaggo)                              |

## Project Structure

```
shop/
├── bycigar-vue/           # Vue frontend
│   ├── src/components/    # Reusable components (ProductCard, Toast)
│   ├── src/views/         # Page components + admin/
│   ├── src/stores/        # Pinia stores (auth, cart, favorites, toast)
│   └── vite.config.js     # Dev server proxy: /api -> localhost:3000
├── server-go/             # Go backend
│   ├── cmd/main.go        # Entry point
│   ├── internal/handlers/ # HTTP handlers (REST API)
│   ├── internal/middleware/ # CORS, Auth, Admin middleware
│   ├── internal/models/   # GORM models
│   └── pkg/utils/         # response.go (ErrorResponse, SuccessResponse)
└── bycigar_site/          # Legacy static site (reference only)
```

## Build Commands

```bash
# Frontend (bycigar-vue/)
npm run dev       # Dev server at http://localhost:5173
npm run build     # Production build -> dist/

# Backend (server-go/)
go run cmd/main.go           # Dev server at http://localhost:3000
go build -o server cmd/main.go  # Build executable
go fmt ./...                 # Format code

# Docker
docker-compose up -d         # Start all services
docker-compose up -d mysql   # Start MySQL only
docker exec -it bycigar-mysql mysql -uroot -p123456 bycigar
```

## Code Style - Vue Components

```vue
<script setup>
import { ref, computed, onMounted } from 'vue'
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
- Always use `<script setup>` (Composition API)
- Order: imports -> emit/props -> refs -> computed -> functions
- Use `defineEmits()` before `defineProps()` (emit first)
- All styles in `<style scoped>` blocks
- Component files: PascalCase (`ProductCard.vue`)

## Code Style - Go Backend

### Handler Pattern
```go
func GetProducts(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    query := database.DB.Model(&models.Product{}).Where("is_active = ?", true)
    
    var products []models.Product
    if err := query.Find(&products).Error; err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch")
        return
    }
    c.JSON(http.StatusOK, gin.H{"products": products})
}
```

### Go Rules
- Use `gofmt` for formatting
- Import order: standard -> external -> internal
- Always use `utils.ErrorResponse()` for error responses
- Use `gin.H{}` for JSON responses
- Map frontend camelCase to DB snake_case (e.g., `createdAt` -> `created_at`)

## Naming Conventions

| Type          | Convention              | Example                    |
|---------------|------------------------|----------------------------|
| Vue Components| PascalCase             | `ProductCard.vue`          |
| Views         | PascalCase + "View"    | `HomeView.vue`             |
| Stores        | camelCase + "Store"    | `useCartStore`             |
| Props         | camelCase              | `productId`                |
| CSS Classes   | kebab-case             | `product-card`             |
| API Endpoints | kebab-case             | `/api/cart-items`          |

## Error Handling

### Frontend
```javascript
async function fetchData() {
  try {
    loading.value = true
    const res = await fetch(`${API_BASE}/products`)
    if (!res.ok) throw new Error('Request failed')
    products.value = (await res.json()).products || []
  } catch (e) {
    error.value = e.message
    console.error('Fetch failed:', e)
    toast.error('加载失败')
  } finally {
    loading.value = false
  }
}
```

### Backend
```go
if err := database.DB.First(&user, id).Error; err != nil {
    utils.ErrorResponse(c, http.StatusNotFound, "User not found")
    return
}
```

## User Notifications

Use `useToastStore` for non-blocking notifications. **DO NOT** use `alert()`.

```javascript
import { useToastStore } from '../stores/toast'
const toast = useToastStore()
toast.success('操作成功')
toast.error('操作失败')
```

## API Endpoints

| Method | Endpoint              | Auth | Description               |
|--------|----------------------|------|---------------------------|
| GET    | `/api/products`      | No   | List products (paginated) |
| GET    | `/api/products/:id`  | No   | Get product detail        |
| GET    | `/api/categories`    | No   | List categories           |
| GET    | `/api/cart`          | Yes  | Get user's cart           |
| POST   | `/api/cart`          | Yes  | Add item to cart          |
| PUT    | `/api/cart/:id`      | Yes  | Update quantity           |
| DELETE | `/api/cart/:id`      | Yes  | Remove item               |
| POST   | `/api/auth/login`    | No   | Login                     |
| POST   | `/api/auth/register` | No   | Register                  |
| GET    | `/api/auth/me`       | Yes  | Get current user          |

Query params for `/api/products`: `page`, `limit`, `search`, `category`, `sortBy`, `sortOrder`

## Authentication

- Token: `localStorage.getItem('token')`
- User: `localStorage.getItem('user')` (JSON)
- Header: `Authorization: Bearer <token>`
- Protected routes: `meta: { requiresAuth: true }`
- Admin routes: `meta: { requiresAuth: true, requiresAdmin: true }`

## Git Commit Convention

Use conventional commits with Chinese descriptions:
```
feat: 添加商品收藏功能
fix: 修复购物车数量计算错误
refactor: 重构用户认证逻辑
```

## Important Notes

- **No tests yet** - Testing framework not configured
- **Language**: Communicate with user in Chinese
- **Simplicity**: Keep code clean, avoid over-engineering
- **Comments**: Only for complex logic explanations
- **Step-by-step**: Execute one phase at a time, report success, then proceed
