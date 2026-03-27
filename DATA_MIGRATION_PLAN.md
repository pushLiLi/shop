# 数据化迁移计划

> 目标：将静态网站内容迁移到 MySQL + MinIO，实现数据驱动的动态网站

---

## 一、项目现状分析

| 类型 | 数量 | 说明 |
|------|------|------|
| 图片文件 | ~2089 张 | 产品图、品牌图、Banner、静态资源 |
| 产品 | ~66 个 | 包含 SKU、价格、描述、属性等 |
| 品牌 | ~110 个 | 古巴品牌、非古巴品牌 |
| 分类 | ~81 个 | 多级分类结构 |
| 轮播图 | 3 张 | 首页 Banner |
| 静态页面 | 5 个 | 关于我们、服务条款、隐私政策、退换政策、品牌列表 |
| 导航菜单 | 3 个一级 + 多个二级 | 古巴品牌、世界雪茄、烟草 |

---

## 二、技术架构

```
┌─────────────────────────────────────────────────────────┐
│                    Docker Compose                        │
├─────────────┬─────────────┬─────────────────────────────┤
│   MySQL     │   MinIO     │   Node.js API Server        │
│   :3306     │   :9000     │   :3000                     │
│   (数据)     │   (图片)    │   (接口)                    │
└─────────────┴─────────────┴─────────────────────────────┘
                          │
                          ▼
              ┌─────────────────────┐
              │   Vue 3 Frontend    │
              │   :5173 (dev)       │
              └─────────────────────┘
```

---

## 三、数据库设计（MySQL + Prisma）

### 3.1 表结构概览

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    brands    │     │  categories  │     │   products   │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id           │     │ id           │     │ id           │
│ name         │     │ name         │     │ name         │
│ name_en      │     │ name_en      │     │ name_en      │
│ slug         │     │ slug         │     │ slug         │
│ description  │     │ parent_id    │◄────│ category_id  │
│ image_url    │     │ brand_id     │◄────│ brand_id     │
│ is_active    │     │ image_url    │     │ sku          │
│ sort_order   │     │ is_active    │     │ price        │
└──────────────┘     └──────────────┘     │ image_url    │
                                          │ attributes   │
                                          │ is_featured  │
                                          └──────────────┘

┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   banners    │     │  menu_items  │     │ static_pages │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id           │     │ id           │     │ id           │
│ title        │     │ name         │     │ slug         │
│ subtitle     │     │ name_en      │     │ title        │
│ image_url    │     │ link_url     │     │ content      │
│ link_url     │     │ parent_id    │     │ is_active    │
│ position     │     │ is_active    │     └──────────────┘
│ is_active    │     │ sort_order   │
│ sort_order   │     └──────────────┘
└──────────────┘

┌──────────────┐     ┌──────────────┐
│ site_config  │     │product_images│
├──────────────┤     ├──────────────┤
│ id           │     │ id           │
│ site_name    │     │ product_id   │
│ logo_url     │     │ image_url    │
│ description  │     │ sort_order   │
│ keywords     │     └──────────────┘
│ home_featured│
│ home_cuban   │
│ phone        │
│ email        │
└──────────────┘
```

### 3.2 Prisma Schema 详细定义

```prisma
// server/prisma/schema.prisma

generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "mysql"
  url      = env("DATABASE_URL")
}

// 品牌
model Brand {
  id          Int       @id @default(autoincrement())
  name        String    @db.VarChar(100)
  nameEn      String?   @map("name_en") @db.VarChar(100)
  slug        String    @unique @db.VarChar(100)
  description String?   @db.Text
  imageUrl    String?   @map("image_url") @db.VarChar(500)
  isActive    Boolean   @default(true) @map("is_active")
  sortOrder   Int       @default(0) @map("sort_order")
  products    Product[]
  categories  Category[]
  createdAt   DateTime  @default(now()) @map("created_at")
  updatedAt   DateTime  @updatedAt @map("updated_at")
  
  @@index([slug])
  @@index([isActive, sortOrder])
  @@map("brands")
}

// 分类
model Category {
  id          Int        @id @default(autoincrement())
  name        String     @db.VarChar(100)
  nameEn      String?    @map("name_en") @db.VarChar(100)
  slug        String     @unique @db.VarChar(100)
  parentId    Int?       @map("parent_id")
  parent      Category?  @relation("CategoryTree", fields: [parentId], references: [id], onDelete: Cascade)
  children    Category[] @relation("CategoryTree")
  brandId     Int?       @map("brand_id")
  brand       Brand?     @relation(fields: [brandId], references: [id])
  imageUrl    String?    @map("image_url") @db.VarChar(500)
  isActive    Boolean    @default(true) @map("is_active")
  sortOrder   Int        @default(0) @map("sort_order")
  products    Product[]
  createdAt   DateTime   @default(now()) @map("created_at")
  updatedAt   DateTime   @updatedAt @map("updated_at")
  
  @@index([slug])
  @@index([parentId])
  @@index([brandId])
  @@index([isActive, sortOrder])
  @@map("categories")
}

// 产品
model Product {
  id            Int            @id @default(autoincrement())
  name          String         @db.VarChar(200)
  nameEn        String?        @map("name_en") @db.VarChar(200)
  slug          String         @unique @db.VarChar(200)
  sku           String?        @unique @db.VarChar(50)
  price         Float
  originalPrice Float?         @map("original_price")
  description   String?        @db.Text
  stock         Int            @default(0)
  isActive      Boolean        @default(true) @map("is_active")
  isFeatured    Boolean        @default(false) @map("is_featured")
  
  brandId       Int?           @map("brand_id")
  brand         Brand?         @relation(fields: [brandId], references: [id])
  categoryId    Int?           @map("category_id")
  category      Category?      @relation(fields: [categoryId], references: [id])
  
  imageUrl      String?        @map("image_url") @db.VarChar(500)
  attributes    Json?
  
  images        ProductImage[]
  createdAt     DateTime       @default(now()) @map("created_at")
  updatedAt     DateTime       @updatedAt @map("updated_at")
  
  @@index([slug])
  @@index([sku])
  @@index([brandId])
  @@index([categoryId])
  @@index([isActive, isFeatured])
  @@map("products")
}

// 产品图片
model ProductImage {
  id        Int     @id @default(autoincrement())
  productId Int     @map("product_id")
  product   Product @relation(fields: [productId], references: [id], onDelete: Cascade)
  imageUrl  String  @map("image_url") @db.VarChar(500)
  sortOrder Int     @default(0) @map("sort_order")
  
  @@index([productId])
  @@map("product_images")
}

// 轮播图
model Banner {
  id        Int      @id @default(autoincrement())
  title     String?  @db.VarChar(100)
  subtitle  String?  @db.VarChar(200)
  imageUrl  String   @map("image_url") @db.VarChar(500)
  linkUrl   String?  @map("link_url") @db.VarChar(500)
  position  String   @default("home") @db.VarChar(50)
  isActive  Boolean  @default(true) @map("is_active")
  sortOrder Int      @default(0) @map("sort_order")
  createdAt DateTime @default(now()) @map("created_at")
  updatedAt DateTime @updatedAt @map("updated_at")
  
  @@index([position, isActive, sortOrder])
  @@map("banners")
}

// 导航菜单
model MenuItem {
  id        Int        @id @default(autoincrement())
  name      String     @db.VarChar(50)
  nameEn    String?    @map("name_en") @db.VarChar(50)
  linkUrl   String?    @map("link_url") @db.VarChar(200)
  parentId  Int?       @map("parent_id")
  parent    MenuItem?  @relation("MenuTree", fields: [parentId], references: [id], onDelete: Cascade)
  children  MenuItem[] @relation("MenuTree")
  isActive  Boolean    @default(true) @map("is_active")
  sortOrder Int        @default(0) @map("sort_order")
  createdAt DateTime   @default(now()) @map("created_at")
  updatedAt DateTime   @updatedAt @map("updated_at")
  
  @@index([parentId])
  @@index([isActive, sortOrder])
  @@map("menu_items")
}

// 静态页面
model StaticPage {
  id        Int      @id @default(autoincrement())
  slug      String   @unique @db.VarChar(50)
  title     String   @db.VarChar(100)
  content   String   @db.Text
  isActive  Boolean  @default(true) @map("is_active")
  createdAt DateTime @default(now()) @map("created_at")
  updatedAt DateTime @updatedAt @map("updated_at")
  
  @@index([slug])
  @@map("static_pages")
}

// 站点配置
model SiteConfig {
  id                Int     @id @default(1)
  siteName          String? @map("site_name") @db.VarChar(100)
  logoUrl           String? @map("logo_url") @db.VarChar(500)
  description       String? @db.Text
  keywords          String? @db.VarChar(500)
  homeFeaturedTitle String? @map("home_featured_title") @db.VarChar(50)
  homeCubanTitle    String? @map("home_cuban_title") @db.VarChar(50)
  phone             String? @db.VarChar(50)
  email             String? @db.VarChar(100)
  address           String? @db.VarChar(200)
  wechat            String? @db.VarChar(100)
  weibo             String? @db.VarChar(100)
  
  @@map("site_config")
}
```

---

## 四、MinIO 存储结构

```
minio-data/
└── bycigar/                    # Bucket 名称
    ├── products/               # 产品图片
    │   ├── original/           # 原图 (最大尺寸)
    │   ├── large/              # 大图 (600x600)
    │   └── thumbnail/          # 缩略图 (100x100)
    ├── brands/                 # 品牌图片
    │   └── original/
    ├── categories/             # 分类图片
    │   └── original/
    ├── banners/                # 轮播图
    │   └── original/
    ├── static/                 # 静态资源
    │   ├── logo.png
    │   └── favicon.png
    └── misc/                   # 其他图片
```

### 图片 URL 格式

```
开发环境: http://localhost:9000/bycigar/products/original/abc123.jpg
生产环境: https://cdn.bycigar.com/products/original/abc123.jpg
```

---

## 五、执行阶段

### 阶段 1：基础设施搭建（1-2 天）

- [ ] 1.1 创建 `docker-compose.yml`
- [ ] 1.2 初始化后端项目 `server/`
- [ ] 1.3 配置 Prisma
- [ ] 1.4 执行数据库迁移
- [ ] 1.5 配置 MinIO 客户端
- [ ] 1.6 创建 MinIO Bucket

### 阶段 2：图片迁移到 MinIO（1 天）

- [ ] 2.1 编写图片扫描脚本
- [ ] 2.2 上传图片到 MinIO
- [ ] 2.3 生成路径映射表
- [ ] 2.4 验证图片完整性

### 阶段 3：数据提取与导入（2-3 天）

- [ ] 3.1 解析品牌页面
- [ ] 3.2 解析分类页面
- [ ] 3.3 解析产品页面
- [ ] 3.4 解析轮播图
- [ ] 3.5 解析导航菜单
- [ ] 3.6 解析静态页面
- [ ] 3.7 导入数据库
- [ ] 3.8 数据验证

### 阶段 4：API 开发（2 天）

- [ ] 4.1 产品相关 API
- [ ] 4.2 品牌相关 API
- [ ] 4.3 分类相关 API
- [ ] 4.4 轮播图 API
- [ ] 4.5 导航菜单 API
- [ ] 4.6 静态页面 API
- [ ] 4.7 站点配置 API
- [ ] 4.8 图片代理 API

### 阶段 5：前端重构（2-3 天）

- [ ] 5.1 轮播图数据动态化
- [ ] 5.2 产品列表数据动态化
- [ ] 5.3 导航菜单数据动态化
- [ ] 5.4 静态页面数据动态化
- [ ] 5.5 图片路径更新
- [ ] 5.6 状态管理优化

### 阶段 6：测试与清理（1 天）

- [ ] 6.1 功能测试
- [ ] 6.2 性能测试
- [ ] 6.3 清理旧资源

---

## 六、API 设计

### 6.1 产品 API

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/products` | GET | 产品列表（支持分页、筛选） |
| `/api/products/:id` | GET | 产品详情（ID） |
| `/api/products/slug/:slug` | GET | 产品详情（Slug） |
| `/api/products/featured` | GET | 推荐产品 |

**查询参数：**
```
?page=1&limit=20
&brand=cohiba
&category=cuban-cigars
&is_featured=true
&sort=price_asc
```

### 6.2 品牌 API

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/brands` | GET | 品牌列表 |
| `/api/brands/:slug` | GET | 品牌详情 |
| `/api/brands/:slug/products` | GET | 品牌下的产品 |

### 6.3 分类 API

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/categories` | GET | 分类树 |
| `/api/categories/:slug` | GET | 分类详情 |
| `/api/categories/:slug/products` | GET | 分类下的产品 |

### 6.4 其他 API

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/banners` | GET | 轮播图列表 |
| `/api/menu` | GET | 导航菜单 |
| `/api/pages/:slug` | GET | 静态页面 |
| `/api/config` | GET | 站点配置 |
| `/api/images/*` | GET | MinIO 图片代理 |

---

## 七、Docker Compose 配置

```yaml
# docker-compose.yml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: bycigar-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD:-root123}
      MYSQL_DATABASE: bycigar
      MYSQL_USER: ${MYSQL_USER:-bycigar}
      MYSQL_PASSWORD: ${MYSQL_PASSWORD:-bycigar123}
    ports:
      - "${MYSQL_PORT:-3306}:3306"
    volumes:
      - ./server/mysql-data:/var/lib/mysql
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  minio:
    image: minio/minio:latest
    container_name: bycigar-minio
    restart: unless-stopped
    environment:
      MINIO_ROOT_USER: ${MINIO_USER:-minioadmin}
      MINIO_ROOT_PASSWORD: ${MINIO_PASSWORD:-minioadmin123}
    ports:
      - "${MINIO_API_PORT:-9000}:9000"
      - "${MINIO_CONSOLE_PORT:-9001}:9001"
    volumes:
      - ./server/minio-data:/data
    command: server /data --console-address ":9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

networks:
  default:
    name: bycigar-network
```

### 环境变量文件 `.env`

```env
# MySQL
MYSQL_ROOT_PASSWORD=root123
MYSQL_USER=bycigar
MYSQL_PASSWORD=bycigar123
MYSQL_PORT=3306

# MinIO
MINIO_USER=minioadmin
MINIO_PASSWORD=minioadmin123
MINIO_API_PORT=9000
MINIO_CONSOLE_PORT=9001

# Database URL (Prisma)
DATABASE_URL=mysql://bycigar:bycigar123@localhost:3306/bycigar

# MinIO Config
MINIO_ENDPOINT=localhost
MINIO_PORT=9000
MINIO_USE_SSL=false
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin123
MINIO_BUCKET=bycigar
```

---

## 八、目录结构

```
shop/
├── docker-compose.yml
├── .env
├── DATA_MIGRATION_PLAN.md        # 本文档
├── CLAUDE.md                     # 项目指南
│
├── bycigar_site/                 # 原始静态站点（迁移后可删除）
│
├── bycigar-vue/                  # Vue 前端
│   ├── src/
│   │   ├── views/
│   │   ├── components/
│   │   ├── stores/
│   │   ├── router/
│   │   └── api/                  # API 请求封装
│   └── ...
│
└── server/                       # Node.js 后端
    ├── package.json
    ├── prisma/
    │   └── schema.prisma
    ├── src/
    │   ├── index.js              # 入口
    │   ├── routes/               # 路由
    │   │   ├── products.js
    │   │   ├── brands.js
    │   │   ├── categories.js
    │   │   ├── banners.js
    │   │   ├── menu.js
    │   │   ├── pages.js
    │   │   └── config.js
    │   ├── services/             # 业务逻辑
    │   ├── utils/                # 工具函数
    │   └── middleware/           # 中间件
    ├── scripts/                  # 脚本
    │   ├── migrate-images.js     # 图片迁移
    │   ├── parse-products.js     # 产品解析
    │   └── seed.js               # 数据导入
    ├── mysql-data/               # MySQL 数据（挂载）
    └── minio-data/               # MinIO 数据（挂载）
```

---

## 九、后续扩展

### 管理后台（Phase 2）

后期开发管理后台时需要：

1. **管理 API**
   - CRUD 操作
   - 图片上传接口
   - 批量操作

2. **管理界面**
   - Vue Admin 模板
   - 产品管理
   - 品牌管理
   - 分类管理
   - 轮播图管理
   - 订单管理

3. **权限系统**
   - 用户认证
   - 角色权限
   - 操作日志

---

## 十、风险与注意事项

1. **图片迁移**
   - 保持原图质量
   - 注意文件名编码（中文文件名）
   - 备份原始文件

2. **数据完整性**
   - 外键关联正确
   - 路径映射准确
   - 字符编码 UTF-8

3. **性能考虑**
   - API 分页
   - 图片懒加载
   - 数据库索引

4. **安全考虑**
   - 环境变量管理
   - API 访问控制
   - SQL 注入防护

---

## 十一、时间估算

| 阶段 | 预计时间 | 说明 |
|------|----------|------|
| 阶段 1 | 1-2 天 | 基础设施 |
| 阶段 2 | 1 天 | 图片迁移 |
| 阶段 3 | 2-3 天 | 数据导入 |
| 阶段 4 | 2 天 | API 开发 |
| 阶段 5 | 2-3 天 | 前端重构 |
| 阶段 6 | 1 天 | 测试清理 |
| **总计** | **9-12 天** | - |

---

*文档创建时间：2026-03-26*
*最后更新：2026-03-26*
