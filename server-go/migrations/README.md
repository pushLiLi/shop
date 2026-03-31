# BYCIGAR 数据库迁移脚本

本目录包含 BYCIGAR 项目的数据库初始化和迁移脚本，使用 SQL 文件而非代码 Seed。

## 目录结构

```
migrations/
├── 1_schema.sql          # 表结构定义
├── 2_base.sql            # 基础种子数据（管理员、分类、商品）
├── 3_demo.sql            # 演示数据（用户、地址、快捷回复）
├── 4_settings.sql         # 配置数据（设置、页面、SEO）
├── init-db.sh            # Linux/Mac 初始化脚本
├── init-db.bat           # Windows 初始化脚本
├── reset-db.sh           # Linux/Mac 重置脚本
└── README.md             # 本文档
```

## 快速开始

### Linux / Mac

```bash
cd server-go/migrations

# 方式1: 使用环境变量
DB_HOST=localhost DB_USER=root DB_PASSWORD=123456 ./init-db.sh

# 方式2: 使用默认配置（root/123456）
./init-db.sh

# 初始化测试数据库
./init-db.sh test
```

### Windows

```cmd
cd server-go\migrations

# 运行初始化脚本
init-db.bat

# 初始化测试数据库
init-db.bat
test
```

## 脚本说明

### 001_init_schema.sql

创建所有数据库表，兼容 GORM AutoMigrate 的字段定义：

- `users` - 用户表
- `categories` - 分类表
- `products` - 商品表
- `addresses` - 收货地址表
- `cart_items` - 购物车表
- `favorites` - 收藏表
- `orders` - 订单表
- `order_items` - 订单明细表
- `site_configs` - 网站配置表
- `banners` - 横幅表
- `pages` - 页面表
- `settings` - 设置表
- `payment_methods` - 支付方式表
- `payment_proofs` - 支付凭证表
- `contact_methods` - 联系方式表
- `conversations` - 客服会话表
- `messages` - 客服消息表
- `quick_replies` - 快捷回复表
- `ratings` - 满意度评价表
- `notifications` - 通知表

### 002_seed_base.sql

基础数据（幂等操作，可重复执行）：

- 管理员账号: `admin@bycigar.com` / `123456`
- 客服账号: `service1@bycigar.com` / `123456`, `service2@bycigar.com` / `123456`
- 4个一级分类: 精品雪茄、雪茄配件、生活方式、礼盒套装
- 8个二级分类: 古巴经典、多米尼加、尼加拉瓜、切割工具、保湿存储、点火设备、酒水搭配、入门礼盒
- 64个商品
- 6个Banner
- 8种支付方式
- 6种联系方式

### 003_seed_demo.sql

演示数据：

- 25个客户账号: `user01@test.com` ~ `user25@test.com` / `123456`
- 对应的收货地址
- 8个快捷回复

### 004_seed_settings.sql

配置数据：

- 10条网站设置
- 16条网站配置
- 5个CMS页面: 关于我们、服务条款、隐私政策、免责说明、配送说明

## Docker 环境使用

**重要**: MySQL Docker 镜像会在**首次启动**时自动执行 `/docker-entrypoint-initdb.d/` 目录下的 `.sql` 文件。

```bash
# 1. 删除旧数据卷（如果需要重新初始化）
docker-compose down -v

# 2. 启动服务（会自动执行 migrations/*.sql）
docker-compose up -d mysql

# 3. 等待初始化完成
docker logs bycigar-mysql
```

> **注意**: 重新初始化需要删除数据卷 `docker-compose down -v`，否则 SQL 脚本不会重新执行。

### 手动执行 SQL

```bash
# 进入容器执行
docker exec -i bycigar-mysql mysql -uroot -p123456 bycigar < migrations/1_schema.sql
docker exec -i bycigar-mysql mysql -uroot -p123456 bycigar < migrations/2_base.sql
docker exec -i bycigar-mysql mysql -uroot -p123456 bycigar < migrations/3_demo.sql
docker exec -i bycigar-mysql mysql -uroot -p123456 bycigar < migrations/4_settings.sql
```

## 手动执行

```bash
# 完整初始化
mysql -h localhost -u root -p bycigar < migrations/1_schema.sql
mysql -h localhost -u root -p bycigar < migrations/2_base.sql
mysql -h localhost -u root -p bycigar < migrations/3_demo.sql
mysql -h localhost -u root -p bycigar < migrations/4_settings.sql
```

## 注意事项

1. **幂等性**: 所有脚本使用 `ON DUPLICATE KEY UPDATE`，可安全重复执行
2. **执行顺序**: 必须按编号顺序执行（001 -> 002 -> 003 -> 004）
3. **密码哈希**: 所有用户密码统一为 `123456`，使用 bcrypt 哈希存储
4. **ID固定**: 管理员 ID 为 1-3，客户 ID 从 4 开始，避免与现有数据冲突
5. **生产环境**: 请修改默认密码和邮箱地址

## 故障排除

### 连接被拒绝

```bash
# 检查 MySQL 是否运行
docker ps | grep mysql

# 检查端口
netstat -an | grep 3306  # Windows
lsof -i :3306            # Mac/Linux
```

### 表已存在错误

脚本使用 `DROP TABLE IF EXISTS` 和 `ON DUPLICATE KEY UPDATE`，理论上不会报错。如果遇到问题，先手动删除数据库：

```sql
DROP DATABASE bycigar;
CREATE DATABASE bycigar CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 账号汇总

| 角色 | Email | 密码 |
|------|-------|------|
| 管理员 | admin@bycigar.com | 123456 |
| 管理员 | service1@bycigar.com | 123456 |
| 管理员 | service2@bycigar.com | 123456 |
| 客户 | user01@test.com | 123456 |
| ... | user02@test.com ~ user25@test.com | 123456 |
