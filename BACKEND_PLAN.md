# Go 后端项目计划

> 为 bycigar-vue 前端创建 Go 语言后端 API

## 技术栈

| 层级 | 技术 |
|------|------|
| 框架 | Gin |
| 数据库 | MySQL |
| ORM | GORM |
| 认证 | JWT (golang-jwt) |
| 配置管理 | Viper + .env |

---

## 项目结构

```
server-go/
├── cmd/
│   └── main.go              # 程序入口
├── internal/
│   ├── config/              # 配置加载
│   │   └── config.go
│   ├── middleware/          # 中间件
│   │   ├── auth.go          # JWT认证中间件
│   │   └── cors.go          # CORS中间件
│   ├── models/              # 数据模型
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── category.go
│   │   ├── cart.go
│   │   ├── favorite.go
│   │   ├── address.go
│   │   ├── order.go
│   │   └── config.go
│   ├── handlers/            # HTTP处理器
│   │   ├── auth.go
│   │   ├── product.go
│   │   ├── cart.go
│   │   ├── favorite.go
│   │   ├── address.go
│   │   ├── order.go
│   │   └── config.go
│   └── database/
│       └── database.go      # 数据库连接
├── pkg/
│   └── utils/
│       └── response.go      # 统一响应格式
├── .env                     # 环境变量
├── .env.example
├── go.mod
└── go.sum
```

---

## 数据模型

### User (用户)
```go
type User struct {
    gorm.Model
    Name     string `gorm:"size:100"`
    Email    string `gorm:"uniqueIndex;size:100"`
    Password string `gorm:"size:255"`
    Role     string `gorm:"default:'user'"` // user/admin
}
```

### Category (分类)
```go
type Category struct {
    gorm.Model
    Name string `gorm:"size:100"`
    Slug string `gorm:"uniqueIndex;size:100"`
}
```

### Product (产品)
```go
type Product struct {
    gorm.Model
    Name        string  `gorm:"size:255"`
    Price       float64
    Description string
    ImageURL    string `gorm:"size:500"`
    Brand       string `gorm:"size:100"`
    IsFeatured  bool   `gorm:"default:false"`
    CategoryID  uint
    Category    Category
}
```

### CartItem (购物车)
```go
type CartItem struct {
    gorm.Model
    UserID    uint
    ProductID uint
    Product   Product
    Quantity  int
}
```

### Favorite (收藏)
```go
type Favorite struct {
    gorm.Model
    UserID    uint
    ProductID uint
    Product   Product
}
```

### Address (地址)
```go
type Address struct {
    gorm.Model
    UserID       uint
    FullName     string `gorm:"size:100"`
    AddressLine1 string `gorm:"size:255"`
    AddressLine2 string `gorm:"size:255"`
    City         string `gorm:"size:100"`
    State        string `gorm:"size:50"`
    ZipCode      string `gorm:"size:20"`
    Phone        string `gorm:"size:20"`
    IsDefault    bool   `gorm:"default:false"`
}
```

### Order (订单)
```go
type Order struct {
    gorm.Model
    UserID     uint
    AddressID  uint
    Address    Address
    Status     string `gorm:"default:'pending'"` // pending/paid/shipped/completed/cancelled
    Total      float64
    Remark     string
    Items      []OrderItem
}
```

### OrderItem (订单项)
```go
type OrderItem struct {
    gorm.Model
    OrderID   uint
    ProductID uint
    Product   Product
    Quantity  int
    Price     float64
}
```

### SiteConfig (网站配置)
```go
type SiteConfig struct {
    gorm.Model
    Key   string `gorm:"uniqueIndex;size:100"`
    Value string
}
```

---

## API 路由

### 认证模块
| 方法 | 路由 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/auth/login` | 登录 | 否 |
| POST | `/api/auth/register` | 注册 | 否 |
| GET | `/api/auth/me` | 获取当前用户 | 是 |
| PUT | `/api/auth/profile` | 更新资料 | 是 |

### 产品模块
| 方法 | 路由 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/products` | 产品列表 (分页、搜索、分类) | 否 |
| GET | `/api/products/:id` | 产品详情 | 否 |
| GET | `/api/config` | 网站配置 | 否 |

### 购物车模块
| 方法 | 路由 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/cart` | 购物车列表 | 是 |
| POST | `/api/cart` | 添加到购物车 | 是 |
| PUT | `/api/cart/:id` | 更新数量 | 是 |
| DELETE | `/api/cart/:id` | 删除商品 | 是 |

### 收藏模块
| 方法 | 路由 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/favorites` | 收藏列表 | 是 |
| POST | `/api/favorites` | 添加收藏 | 是 |
| DELETE | `/api/favorites/:id` | 删除收藏 | 是 |

### 地址模块
| 方法 | 路由 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/addresses` | 地址列表 | 是 |
| POST | `/api/addresses` | 添加地址 | 是 |
| PUT | `/api/addresses/:id` | 更新地址 | 是 |
| DELETE | `/api/addresses/:id` | 删除地址 | 是 |
| PUT | `/api/addresses/:id/default` | 设为默认 | 是 |

### 订单模块
| 方法 | 路由 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/orders` | 订单列表 | 是 |
| POST | `/api/orders` | 创建订单 | 是 |

---

## 查询参数

### GET /api/products

| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码 (默认 1) |
| limit | int | 每页数量 (默认 12) |
| search | string | 搜索关键词 |
| categoryId | int | 分类ID |
| categorySlug | string | 分类Slug |
| sortBy | string | 排序字段 (createdAt/price/name) |
| sortOrder | string | 排序方向 (asc/desc) |

### 响应格式

```json
{
  "products": [...],
  "total": 100,
  "page": 1,
  "limit": 12,
  "category": { "id": 1, "name": "...", "slug": "..." }
}
```

---

## 环境变量

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=bycigar

# JWT配置
JWT_SECRET=your_jwt_secret_key

# 服务器配置
PORT=3000
```

---

## 执行步骤

| 阶段 | 任务 | 状态 |
|------|------|------|
| 1. 初始化 | 创建项目目录、初始化Go模块 | 待开始 |
| 2. 安装依赖 | 安装Gin、GORM、JWT、Viper等 | 待开始 |
| 3. 数据模型 | 创建所有model文件 | 待开始 |
| 4. 数据库连接 | 配置MySQL连接、自动迁移 | 待开始 |
| 5. 中间件 | JWT认证、CORS | 待开始 |
| 6. API处理器 | 按模块实现handler | 待开始 |
| 7. 路由配置 | 注册所有路由 | 待开始 |
| 8. 测试验证 | 运行并测试API | 待开始 |

---

## 依赖列表

```bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/spf13/viper
go get -u github.com/joho/godotenv
go get -u golang.org/x/crypto/bcrypt
```

---

## 前端配置调整

前端 `API_BASE` 当前为 `http://localhost:3000/api`，Go后端将监听相同端口，无需修改。

如需修改，编辑以下文件：
- `bycigar-vue/src/stores/auth.js`
- `bycigar-vue/src/stores/cart.js`
- `bycigar-vue/src/stores/favorites.js`
- `bycigar-vue/src/views/*.vue`
