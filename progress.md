# BYCIGAR 项目开发进度

## 已完成

### 基础架构 ✅
- [x] Go 后端服务器 (Gin + GORM + MySQL)
- [x] Vue 3 前端 (Vite + Vue Router + Pinia)
- [x] JWT 认证系统
- [x] 管理后台框架

### 核心功能 ✅
- [x] 商品管理 (CRUD)
- [x] 分类管理 (CRUD)
- [x] 轮播图管理
- [x] CMS 页面管理 (Markdown + 数据库)
- [x] 购物车功能
- [x] 订单功能
- [x] 收藏功能
- [x] 收货地址管理
- [x] 用户认证与资料

### 站点设置 ✅ (新增)
- [x] Setting 模型 (key-value 存储)
- [x] 设置 API (GET /api/settings, PUT /api/admin/settings/:key)
- [x] AdminSettings.vue 管理页面
- [x] TheFooter.vue 动态页脚内容
- [x] useSettingsStore.js 状态管理

## 服务器状态

- 后端: http://localhost:3000
- 前端: http://localhost:5173

## 技术栈

- 前端: Vue 3 + Vite + Vue Router + Pinia + marked
- 后端: Go 1.23 + Gin + GORM + JWT + Swagger
- 数据库: MySQL 8.4 (Docker)
