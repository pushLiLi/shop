# AGENTS.md - AI Coding Agent Guidelines

This document provides guidelines for AI coding agents (such as Claude, GPT-4, Copilot, Cursor) etc.) working on this codebase.

## Project Overview

This is a **Vue 3 + Vite** e-commerce project migrated from a legacy static website. The project consists of:
- **Frontend**: `bycigar-vue/` (Vue 3 + Vite + Vue Router + Pinia)
- **Backend**: `server-go/` (Go + Gin + GORM + MySQL)

## Tech Stack

| Layer       | Technology                                    |
|-------------|-----------------------------------------------|
| Framework   | Vue 3 (Composition API, `<script setup>`)     |
| Build Tool  | Vite                                          |
| Language    | JavaScript (Frontend), Go (Backend)           |
| State       | Pinia (only if necessary, otherwise local)    |
| Routing     | Vue Router 4                                  |
| Backend     | Go + Gin                                      |
| Database    | MySQL (GORM)                                  |
| Styling     | Scoped CSS within `.vue` files                |

## Project Structure

```
shop/
├── bycigar-vue/           # Vue frontend
│   ├── src/
│   │   ├── components/    # Reusable Vue components
│   │   ├── views/         # Page-level components
│   │   ├── stores/        # Pinia stores (auth, cart, favorites)
│   │   ├── router/        # Vue Router configuration
│   │   └── main.js        # App entry point
│   └── vite.config.js
├── server-go/             # Go backend
│   ├── cmd/main.go        # Application entry point
│   ├── internal/
│   │   ├── handlers/      # HTTP handlers
│   │   ├── middleware/    # Middleware (CORS, Auth)
│   │   ├── models/        # Data models
│   │   ├── database/      # Database connection
│   │   └── config/        # Configuration
│   └── pkg/utils/         # Utility functions
└── bycigar_site/          # Legacy static site (source reference)
```

## Build Commands

### Frontend (bycigar-vue/)
```bash
cd bycigar-vue
npm run dev      # Start dev server at http://localhost:5173
npm run build    # Build for production
npm run preview  # Preview production build
```

### Backend (server-go/)
```bash
cd server-go
go run cmd/main.go    # Start API server at http://localhost:3000
go build -o server cmd/main.go  # Build executable
./server              # Run built executable
```

### Docker
```bash
docker-compose up -d  # Start all services (MySQL + Backend + Frontend)
```

## Code Style Guidelines

### JavaScript (Frontend)

- **ES Modules**: Use ES modules (`import`/`export`) for all frontend code.
- **Async/Await**: Use `async/await` for asynchronous operations.

### Go (Backend)

- **Standard Go formatting**: Use `gofmt` for formatting.
- **Project structure**: Follow the existing internal/pkg structure.
- **Error handling**: Always handle errors explicitly.

### Vue Components

- **Script Setup**: Always use `<script setup>` syntax (Composition API).
- **Imports at top**: Place all imports at the beginning of `<script setup>`.
- **Props definition**: Use `defineProps()` before other code.
- **Emits definition**: Use `defineEmits()` before other code.

```vue
<script setup>
import { ref, computed, onMounted } from 'vue'

const props = defineProps({
  product: { type: Object, required: true }
})

const emit = defineEmits(['update', 'delete'])
</script>
```

- **Naming**: Use PascalCase for component files (`ProductCard.vue`, `TheHeader.vue`).
- **Scoped CSS**: All component styles should be in `<style scoped>` blocks.

### State Management (Pinia)

```javascript
export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [],
    loading: false
  }),
  getters: {
    total: (state) => state.items.reduce((sum, item) => sum + item.price, 0)
  },
  actions: {
    async fetchCart() { /* ... */ }
  }
})
```

### API Calls

```javascript
const API_BASE = 'http://localhost:3000/api'

async function fetchData() {
  try {
    loading.value = true
    const res = await fetch(`${API_BASE}/products`)
    const data = await res.json()
  } catch (e) {
    error.value = e.message
    console.error('Error:', e)
  } finally {
    loading.value = false
  }
}
```

### Error Handling

- Always wrap API calls in try/catch blocks.
- Log errors to console with `console.error()`.
- Show user-friendly error messages in UI.
- Use `finally` for cleanup.

### Naming Conventions

| Type          | Convention              | Example                    |
|---------------|------------------------|----------------------------|
| Components    | PascalCase             | `ProductCard.vue`          |
| Views         | PascalCase + "View"    | `HomeView.vue`             |
| Stores        | camelCase              | `useCartStore`             |
| Props         | camelCase              | `productId`                |
| CSS Classes   | kebab-case             | `product-card`             |
| API Endpoints | kebab-case             | `/api/cart-items`          |

### Backend API Structure

| Method | Endpoint              | Description               |
|--------|----------------------|---------------------------|
| GET    | `/api/products`      | List products (with filters) |
| GET    | `/api/products/:id`  | Get single product        |
| GET    | `/api/categories`    | List categories           |
| GET    | `/api/cart`          | Get user's cart           |
| POST   | `/api/cart`          | Add item to cart          |
| PUT    | `/api/cart/:id`      | Update cart item quantity |
| DELETE | `/api/cart/:id`      | Remove cart item          |
| POST   | `/api/auth/login`    | User login                |
| POST   | `/api/auth/register` | User registration         |

### Authentication

- Token stored in `localStorage` as `token`.
- User info stored in `localStorage` as `user`.
- Authorization header: `Bearer <token>`.
- Protected routes use `meta: { requiresAuth: true }` in router.

### Development Workflow

1. Start MySQL: `docker-compose up -d mysql`
2. Start backend: `cd server-go && go run cmd/main.go`
3. Start frontend: `cd bycigar-vue && npm run dev`

### Important Notes

- **Step-by-Step**: Execute one phase at a time, report success, then proceed.
- **Error Handling**: Retry failed commands up to 3 times before asking.
- **Simplicity**: Keep code clean. Avoid over-engineering.
- **Comments**: Add comments only for complex logic explanations.
- **Language**: Primary communication with user in Chinese.
