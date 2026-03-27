# BYCIGAR 电商项目开发蓝图

> 文档创建日期: 2026-03-26
> 最后更新: 2026-03-26

---

## 项目概览

### 当前状态

- **前端骨架**: 已完成 - Vue 3 + Vite 项目结构搭建完毕
- **路由系统**: 已完成 - 9个页面路由已配置
- **公共组件**: 已完成 - Header/Footer 组件已实现
- **首页**: 部分完成 - UI 已实现，但数据为硬编码
- **后端**: 未开始 - server 目录不存在

### 技术栈

| 层级     | 技术                    |
| -------- | ----------------------- |
| 前端框架 | Vue 3 (Composition API) |
| 构建工具 | Vite                    |
| 路由     | Vue Router 4            |
| 状态管理 | 待定 (Pinia / Local)    |
| 后端     | Node.js + Express       |
| 数据库   | SQLite (Prisma ORM)     |
| 样式     | Scoped CSS              |

---

## 页面功能清单

### 已实现页面

| 路由路径 | 页面名称  | 当前状态  | 说明                        |
| -------- | --------- | --------- | --------------------------- |
| `/`      | HomeView  | ✅ UI完成 | 轮播图+产品展示，数据硬编码 |
| -        | TheHeader | ✅ 完成   | 导航、搜索、用户图标        |
| -        | TheFooter | ✅ 完成   | 链接、版权信息              |

### 待开发页面

| 路由路径          | 页面名称              | 核心功能点                             | 所需后端 API                       | 数据模型             | 优先级 |
| ----------------- | --------------------- | -------------------------------------- | ---------------------------------- | -------------------- | ------ |
| `/category/:slug` | CategoryView          | 分类产品列表、筛选、排序、分页         | `GET /api/products?category=`      | Product, Category    | P0     |
| `/products/:id`   | ProductDetailView     | 产品详情、加入购物车、收藏             | `GET /api/products/:id`            | Product, Review      | P0     |
| `/login`          | LoginView             | 登录表单、JWT认证、错误提示            | `POST /api/auth/login`             | User                 | P0     |
| `/register`       | RegisterView (待建)   | 注册表单、邮箱验证、密码加密           | `POST /api/auth/register`          | User                 | P1     |
| `/cart`           | CartView              | 购物车列表、数量修改、删除、小计       | `GET/POST/DELETE /api/cart`        | Cart, CartItem       | P0     |
| `/favorites`      | FavoritesView         | 收藏列表、添加/移除收藏                | `GET/POST/DELETE /api/favorites`   | Favorite             | P1     |
| `/user`           | UserCenterView (待建) | 个人信息、订单历史、地址管理           | `GET /api/user`, `GET /api/orders` | User, Order, Address | P1     |
| `/checkout`       | CheckoutView (待建)   | 结算流程、地址选择、支付方式           | `POST /api/orders`                 | Order, OrderItem     | P0     |
| `/about`          | AboutView             | 静态内容展示                           | 无 / `GET /api/pages/about`        | Page (可选)          | P2     |
| `/services`       | ServicesView          | 静态内容展示                           | 无 / `GET /api/pages/services`     | Page (可选)          | P2     |
| `/privacy-policy` | PrivacyPolicyView     | 静态内容展示                           | 无                                 | 无                   | P2     |
| `/returns-policy` | ReturnsPolicyView     | 静态内容展示                           | 无                                 | 无                   | P2     |
| `/search`         | SearchView (待建)     | 搜索结果、关键词高亮                   | `GET /api/products?search=`        | Product              | P1     |
| `/admin/*`        | AdminViews (待建)     | 管理后台：产品管理、订单管理、用户管理 | `GET/POST/PUT/DELETE /api/admin/*` | 全部模型             | P1     |

---

## 数据库模型设计 (预估)

### 核心模型

```
┌─────────────┐     ┌─────────────────┐     ┌─────────────┐
│    User     │────<│   CartItem      │>────│   Product   │
├─────────────┤     ├─────────────────┤     ├─────────────┤
│ id          │     │ id              │     │ id          │
│ email       │     │ userId          │     │ name        │
│ password    │     │ productId       │     │ price       │
│ name        │     │ quantity        │     │ description │
│ role        │     │ createdAt       │     │ imageUrl    │
│ createdAt   │     └─────────────────┘     │ categoryId  │
└─────────────┘                             │ brandId     │
      │                                     │ stock       │
      │         ┌─────────────────┐         │ isActive    │
      │         │    Favorite     │         └─────────────┘
      │         ├─────────────────┤               │
      └────────>│ id              │<──────────────┘
                │ userId          │
                │ productId       │     ┌─────────────┐
                └─────────────────┘     │  Category   │
                                        ├─────────────┤
      ┌─────────────┐                   │ id          │
      │    Order    │                   │ name        │
      ├─────────────┤                   │ slug        │
      │ id          │                   │ parentId    │
      │ userId      │                   └─────────────┘
      │ status      │
      │ total       │     ┌─────────────────┐
      │ address     │     │   OrderItem     │
      │ createdAt   │────>│ id              │
      └─────────────┘     │ orderId         │
                          │ productId       │
                          │ quantity        │
                          │ price           │
                          └─────────────────┘
```

### 模型字段详情

#### User (用户)

| 字段      | 类型     | 说明            |
| --------- | -------- | --------------- |
| id        | Int      | 主键，自增      |
| email     | String   | 邮箱，唯一      |
| password  | String   | 加密密码        |
| name      | String?  | 用户名          |
| phone     | String?  | 电话            |
| role      | Enum     | customer, admin |
| createdAt | DateTime | 创建时间        |

#### Product (产品)

| 字段        | 类型     | 说明     |
| ----------- | -------- | -------- |
| id          | Int      | 主键     |
| name        | String   | 产品名称 |
| price       | Float    | 价格     |
| description | String?  | 描述     |
| imageUrl    | String   | 图片路径 |
| categoryId  | Int?     | 分类ID   |
| brandId     | Int?     | 品牌ID   |
| stock       | Int      | 库存     |
| isActive    | Boolean  | 是否上架 |
| createdAt   | DateTime | 创建时间 |

#### Category (分类)

| 字段     | 类型   | 说明     |
| -------- | ------ | -------- |
| id       | Int    | 主键     |
| name     | String | 分类名   |
| slug     | String | URL别名  |
| parentId | Int?   | 父分类ID |

#### SiteConfig (站点动态配置 - 新增)

| 字段         | 类型     | 说明                                      |
| ------------ | -------- | ----------------------------------------- |
| id           | Int      | 主键                                      |
| config_key   | String   | 配置键名 (唯一，如 home_banner_text)      |
| config_value | String   | 配置内容 (支持长文本或序列化的 JSON)      |
| type         | String   | 数据类型 (text, image_url, boolean, json) |
| updatedAt    | DateTime | 最后更新时间                              |

#### Order (订单)

| 字段      | 类型     | 说明                                         |
| --------- | -------- | -------------------------------------------- |
| id        | Int      | 主键                                         |
| userId    | Int      | 用户ID                                       |
| status    | Enum     | pending, paid, shipped, completed, cancelled |
| total     | Float    | 总金额                                       |
| address   | String   | 收货地址                                     |
| phone     | String   | 联系电话                                     |
| createdAt | DateTime | 创建时间                                     |

---

## 依赖关系图

```
┌────────────────────────────────────────────────────────────────┐
│                      Phase 1: 用户认证                         │
│  ┌──────────┐    ┌──────────┐    ┌──────────────┐             │
│  │  Login   │───>│  JWT     │───>│ Route Guard  │             │
│  │ Register │    │  Store   │    │ (权限控制)    │             │
│  └──────────┘    └──────────┘    └──────────────┘             │
└────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────────┐
│                    Phase 2: 核心业务数据                       │
│  ┌──────────────┐    ┌────────────┐    ┌────────────┐         │
│  │  Category    │───>│  Product   │───>│  Search    │         │
│  │  (分类列表)   │    │  (产品详情) │    │  (搜索)    │         │
│  └──────────────┘    └────────────┘    └────────────┘         │
└────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────────┐
│                     Phase 3: 交易流程                          │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌─────────┐  │
│  │  Cart    │───>│ Checkout │───>│  Order   │───>│ Payment │  │
│  │ (购物车)  │    │ (结算)    │    │ (订单)    │    │(模拟支付)│  │
│  └──────────┘    └──────────┘    └──────────┘    └─────────┘  │
│                                                                 │
│  依赖: Phase 1 (用户登录)                                       │
└────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────────┐
│                  Phase 4: 管理后台与编辑                        │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐     │
│  │ Admin Auth   │───>│ Product CRUD │───>│ Order Mgmt   │     │
│  │ (管理员登录)  │    │ (产品管理)    │    │ (订单管理)    │     │
│  └──────────────┘    └──────────────┘    └──────────────┘     │
│                                                                 │
│  依赖: Phase 1 + Phase 2 + Phase 3                             │
└────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────────┐
│                    Phase 5: 优化与部署                          │
│  ┌────────────┐    ┌────────────┐    ┌────────────┐           │
│  │    SEO     │───>│  MinIO     │───>│  Deploy    │           │
│  │ (优化)      │    │ (图片存储)  │    │ (部署)      │           │
│  └────────────┘    └────────────┘    └────────────┘           │
└────────────────────────────────────────────────────────────────┘
```

---

## 分阶段执行计划

### Phase 1: 动态架构与编辑模式 (当前最优先)

**目标**: 建立完全动态的内容管理机制，优先实现首页的数据解耦，并加入管理员“编辑模式”基础设施（遵循 CLAUDE 规范暂缓普通用户Auth，引入简易 Admin/Mock Auth 进行编辑模式的通信）。

**涉及文件**:

- 后端 (新建)
  - `server/index.js` - Express 服务入口，配置 `public/uploads` 静态目录
  - `server/prisma/schema.prisma` - 初始化 Product, SiteConfig 模型
  - `server/routes/api.js` - `GET /api/config`（全局配置） 与 `GET /api/products` (替代合并的 `/api/public/home`)
  - `server/routes/admin.js` - 管理接口 (`POST /api/auth/login`, `PUT /api/admin/config/:key`, `POST /api/admin/upload`)
  - `server/seed.js` - 预设初始站点配置（如 "欢迎来到雪茄之家" 等文案）及硬编码产品数据
- 前端 (修改/新建)
  - `src/stores/auth.js` - Pinia 状态，存储 `userToken` 与 `isAdmin`
  - `src/components/EditableText.vue` - 可在 span 和 input 间动态切换的复用组件，失去焦点时自动调用 API
  - `src/components/EditableImage.vue` - 支持点击上传及预览的图片编辑组件
  - `src/views/HomeView.vue` - 将原硬编码文案、Banner 和产品替换为 `SiteConfig` 取值及 Editable 组件

**验收标准**:

- [ ] 仅运行 `npm run dev` 前端即可读取后端 SQLite (Product, SiteConfig) 数据。
- [ ] `EditableText` 在管理员状态(isAdmin)下能渲染为输入框，失去焦点触发自动保存。
- [ ] 上传图片通过 `POST /api/admin/upload` 保存到 `public/uploads` 并可实时渲染。

---

### Phase 2: 核心业务数据与用户展示

**目标**: 实现产品分类、产品列表、产品详情、搜索功能

**涉及文件**:

- 后端
  - `server/routes/products.js` - 产品 API
  - `server/routes/categories.js` - 分类 API
  - `server/prisma/schema.prisma` - 添加 Product, Category 模型
- 前端
  - `src/views/CategoryView.vue` - 分类页面
  - `src/views/ProductDetailView.vue` - 产品详情 (新建)
  - `src/views/SearchView.vue` - 搜索结果 (新建)
  - `src/components/ProductCard.vue` - 产品卡片组件 (新建)
  - `src/components/ProductGrid.vue` - 产品网格组件 (新建)
  - `src/router/index.js` - 添加新路由

**验收标准**:

- [ ] 分类页面显示该分类下的所有产品
- [ ] 支持按价格排序 (升序/降序)
- [ ] 支持分页 (每页12条)
- [ ] 产品详情页显示完整信息
- [ ] 搜索功能正常工作

---

### Phase 3: 交易流程

**目标**: 实现购物车、结算、订单生成的完整交易流程

**依赖**: Phase 1 (用户登录)

**涉及文件**:

- 后端
  - `server/routes/cart.js` - 购物车 API
  - `server/routes/orders.js` - 订单 API
  - `server/prisma/schema.prisma` - 添加 Cart, Order 模型
- 前端
  - `src/views/CartView.vue` - 购物车页面
  - `src/views/CheckoutView.vue` - 结算页面 (新建)
  - `src/views/OrdersView.vue` - 订单列表 (新建)
  - `src/stores/cart.js` - 购物车状态 (新建)
  - `src/components/CartItem.vue` - 购物车项组件 (新建)

**验收标准**:

- [ ] 可以添加产品到购物车
- [ ] 购物车支持修改数量、删除商品
- [ ] 购物车显示总价
- [ ] 结算流程完整 (选地址 -> 确认 -> 提交)
- [ ] 订单生成后显示订单号
- [ ] 可在个人中心查看历史订单

---

### Phase 4: 管理后台与编辑模式

**目标**: 为管理员提供产品管理、订单管理界面

**依赖**: Phase 1 + Phase 2 + Phase 3

**涉及文件**:

- 后端
  - `server/routes/admin/products.js` - 产品管理 API
  - `server/routes/admin/orders.js` - 订单管理 API
  - `server/routes/admin/users.js` - 用户管理 API
  - `server/middleware/admin.js` - 管理员权限中间件
  - `server/utils/upload.js` - 图片上传工具
- 前端
  - `src/views/admin/DashboardView.vue` - 管理首页 (新建)
  - `src/views/admin/ProductListView.vue` - 产品列表 (新建)
  - `src/views/admin/ProductEditView.vue` - 产品编辑 (新建)
  - `src/views/admin/OrderListView.vue` - 订单列表 (新建)
  - `src/layouts/AdminLayout.vue` - 管理后台布局 (新建)
  - `src/router/admin.js` - 管理路由模块 (新建)

**验收标准**:

- [ ] 管理员可以登录后台
- [ ] 可以新增/编辑/删除产品
- [ ] 可以上传产品图片
- [ ] 可以查看和更新订单状态
- [ ] 普通用户无法访问管理后台

---

### Phase 5: 优化与部署

**目标**: SEO 优化、性能提升、生产环境配置

**涉及文件**:

- 前端
  - `index.html` - Meta 标签优化
  - `src/router/index.js` - 路由懒加载
  - `vite.config.js` - 生产配置
  - `public/robots.txt` - 爬虫配置
  - `public/sitemap.xml` - 站点地图
- 后端
  - `server/routes/sitemap.js` - 动态站点地图
  - `server/config/minio.js` - MinIO 配置 (可选)
  - `.env.production` - 生产环境变量
- 部署
  - `docker-compose.yml` - Docker 配置
  - `nginx.conf` - Nginx 配置

**验收标准**:

- [ ] Lighthouse 性能分数 > 80
- [ ] 首屏加载时间 < 3秒
- [ ] SEO Meta 标签完整
- [ ] 图片使用懒加载
- [ ] 生产环境正常运行

---

## 风险与注意事项

### 技术风险

1. **图片存储**: 当前使用本地存储，需规划 MinIO/CDN 迁移路径
2. **认证安全**: JWT 需设置合理过期时间，考虑 Refresh Token
3. **并发处理**: SQLite 在高并发下可能成为瓶颈，后期可迁移至 PostgreSQL

### 业务风险

1. **支付集成**: 当前为模拟支付，真实支付需对接第三方
2. **数据迁移**: 硬编码的产品数据需要迁移到数据库

---

## 下一步行动

**当前建议**: 从 **Phase 1: 动态架构与编辑模式** 开始执行

1. 创建 `server` 目录和 Express 项目。
2. 配置 Prisma + SQLite，先从 `SiteConfig` 和 `Product` 的 Schema 设计写起。
3. 编写 `seed.js` 初始化默认文案及产品数据。
4. 跑通基础的公开数据查询接口 (`GET /api/config`, `GET /api/products`) 供前端首屏渲染。
5. 前后端逐步完善对应的 Editable 组件架构（可先 Mock `isAdmin=true` 以专注核心组件逻辑）。

---

_文档结束 - 等待审核_
