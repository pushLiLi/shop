# BYCIGAR 生产环境部署文档

> 部署日期：2026-04-01
> 服务器：134.209.165.205
> 域名：huauge.com

---

## 一、环境概况

| 项目 | 详情 |
|------|------|
| 服务器 | Ubuntu 20.04, 2核4GB |
| Docker | 容器化部署，bridge 网络 |
| 数据库 | MySQL 8.4 (Docker) |
| 对象存储 | MinIO (Docker) |
| 后端 | Go 1.25 (Docker) |
| 前端 | Vue 3 + Nginx (Docker) |
| SSL | Let's Encrypt 自动签发 |

---

## 二、服务状态

```
bycigar-mysql     MySQL 8.4        3306     healthy
bycigar-minio     MinIO latest     9000/9001 healthy
bycigar-backend   Go 1.25          3000     running
bycigar-frontend  Nginx + Vue     80/443   running
```

---

## 三、访问地址

| 入口 | 地址 |
|------|------|
| 前台首页 | https://huauge.com |
| 管理后台 | https://huauge.com/admin |
| MinIO 控制台 | http://134.209.165.205:9001 |

---

## 四、账号密码

### 管理员 / 客服

| 账号 | 密码 | 角色 |
|------|------|------|
| admin@bycigar.com | 123456 | admin |
| service1@bycigar.com | 123456 | service |
| service2@bycigar.com | 123456 | service |

### 演示客户（25个）

| 账号 | 密码 |
|------|------|
| user01@test.com ~ user25@test.com | 123456 |

---

## 五、测试数据

### 分类（12个）
- 一级（4个）：精品雪茄、雪茄配件、生活方式、礼盒套装
- 二级（8个）：古巴经典、多米尼加、尼加拉瓜、切割工具、保湿存储、点火设备、酒水搭配、入门礼盒

### 商品（64个）
- 使用 https://picsum.photos 随机图片
- 分布在 8 个二级分类下
- 部分商品 stock=0（缺货）、is_active=0（下架）作为展示

### Banner（6个）
- 古巴经典 · 传承百年
- 多米尼加风情 · 细腻优雅
- 尼加拉瓜激情 · 浓郁澎湃
- 配件专区 · 点亮品茄时刻
- 生活方式 · 雪茄与美酒
- 送礼佳选 · 礼盒套装

### 支付方式（8个）
- 微信支付、支付宝、银行转账、货到付款、PayPal 等

### 联系方式（6个）
- 电话、邮箱、微信、WhatsApp、Telegram、QQ

---

## 六、关键文件路径

```
/opt/bycigar/
├── docker-compose.yml       # 开发环境配置
├── docker-compose.prod.yml  # 生产环境配置
├── .env.production          # 生产环境变量
├── server-go/
│   ├── Dockerfile           # 多阶段构建（Go 1.25）
│   └── migrations/
│       ├── 2_base.sql      # 基础种子数据
│       └── 3_demo.sql      # 演示数据
├── mysql/
│   └── my.cnf               # MySQL 8.4 配置
├── bycigar-vue/
│   └── nginx.conf           # SSL 反向代理配置
└── deploy.sh                # 部署脚本
```

---

## 七、运维命令

### 查看所有容器状态
```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps
```

### 查看日志
```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f frontend
```

### 重启服务
```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml restart backend
```

### 重新导入测试数据（如需）
```bash
docker exec -i bycigar-mysql mysql -uroot -p'123456' bycigar < /opt/bycigar/2_base.sql
docker exec -i bycigar-mysql mysql -uroot -p'123456' bycigar < /opt/bycigar/3_demo.sql
```

### 备份数据库
```bash
docker exec bycigar-mysql mysqldump -uroot -p'123456' bycigar > backup_$(date +%Y%m%d).sql
```

### SSL 证书续期（Let's Encrypt 90天有效期）
```bash
certbot certonly --standalone -d huauge.com -d www.huauge.com --email admin@huauge.com --agree-tos
```

---

## 八、修复记录

### 2026-04-02：修复重启后数据丢失问题

**问题**：后端每次重启，数据库中所有订单数据被清空重建。

**根因**：`server-go/internal/database/database.go` 中的 `SeedBulkOrders()` 函数在每次启动时都会找到所有 customer 用户，删除他们的全部订单（order_items、payment_proofs），然后重新随机生成 3000 条测试订单。

**修复**：删除整个 `SeedBulkOrders()` 函数及其在 `Seed()` 中的调用。生产数据通过 SQL 脚本导入，不需要代码生成测试数据。

---

### 2026-04-02：修复 MySQL 8.4 启动失败

**问题**：MySQL 容器启动后立即崩溃，报 `unknown option '--skip-character-set-client-handshake'`。

**根因**：MySQL 8.4 移除了 `--skip-character-set-client-handshake` 选项，而 `mysql/my.cnf` 中包含了这个配置。

**修复**：从 `mysql/my.cnf` 中删除 `skip-character-set-client-handshake` 行。

---

### 2026-04-02：修复 Nginx 启动失败（host not found in upstream "backend"）

**问题**：Frontend 容器中 Nginx 启动报错 `host not found in upstream "backend"`，导致整个站点无法访问。

**根因**：Nginx 在启动时静态解析 upstream 主机名，而 Docker 容器可能还未完全就绪。`nginx.conf` 中直接写 `proxy_pass http://backend:3000` 会在启动时解析失败。

**修复**：在 `bycigar-vue/nginx.conf` 中改用动态 DNS 解析：
```nginx
resolver 127.0.0.11 valid=30s;
set $backend_upstream http://backend:3000;
proxy_pass $backend_upstream;
```
同理处理 minio 的 upstream。

---

### 2026-04-02：修复 AutoMigrate 外键类型不兼容

**问题**：后端启动时 AutoMigrate 报错 `Error 3780: Referencing column 'parent_id' and referenced column 'id' in foreign key constraint are incompatible`，导致后端反复崩溃重启。

**根因**：SQL schema 脚本（`1_schema.sql`）中所有表的 `id` 和外键列使用 `INT UNSIGNED`，但 GORM 的 `uint` 类型映射为 `BIGINT UNSIGNED`。AutoMigrate 尝试将列从 `INT UNSIGNED` 改为 `BIGINT UNSIGNED` 时，被数据库中已有的 20 个外键约束阻止。本地开发环境不存在此问题，因为本地数据库是 AutoMigrate 直接创建的（已经是 `BIGINT UNSIGNED`）。

**修复**：在 `database.go` 的 `Migrate()` 函数开头，先自动删除所有旧外键约束，再执行 AutoMigrate。同时设置 `DisableForeignKeyConstraintWhenMigrating: true` 防止 GORM 创建新的外键约束。

---

## 九、已知问题

1. **商品图片使用外部占位图**：商品图片为 https://picsum.photos 随机图片，非实际产品图
2. **Banner 图片 404**：`/media/bycigar/banner-1.png` 等文件不存在于 MinIO，需通过管理后台上传真实图片
3. **MinIO 桶策略**：需确认 `bycigar` bucket 是否已创建且策略为公开读

---

## 九、后续建议

1. 上传真实商品图片到 MinIO
2. 上传 Banner 图片并配置到管理后台
3. 配置微信/支付宝商户号以接收真实支付
4. 配置邮件 SMTP 以发送通知邮件
5. 启用自动备份策略
6. 配置 Nginx 限流防止 DDoS

---

### 2026-04-02：代码稳定性与健壮性修复

**问题**：事故分析报告（`incident-2026-04-02.md`）中识别出 11 个可能导致服务崩溃的代码风险。

**修复内容**：

| 风险等级 | 问题 | 修复方案 |
|---------|------|---------|
| 严重 | WebSocket goroutine 无 panic recovery | `ws/hub.go` 添加 defer recover |
| 严重 | 6个后台清理 goroutine 无 recover | 所有 `*_cleanup.go` 添加 panic recovery |
| 严重 | 邮件发送 goroutine 无 recover | `admin_order.go` 添加 defer recover |
| 严重 | 无 graceful shutdown | `main.go` 添加 SIGINT/SIGTERM 处理 |
| 高危 | 17处不安全的类型断言 | 改为 comma-ok 模式 |
| 高危 | loginFailures map 无上限 | 添加 10000 上限 + 后台清理 |
| 高危 | 导出订单无记录限制 | 添加 10000 条上限 |
| 中等 | 无健康检查端点 | 添加 `GET /health` |

**验证**：
```bash
curl https://huauge.com/api/health
# 应返回 {"status":"ok","timestamp":...}
```

---

### 2026-04-02：服务器迁移与新服务器部署

**新服务器**：167.99.115.47 (ubuntu-s-2vcpu-4gb-120gb-intel-nyc3-01)
**域名**：huauge.com

**关键配置**：

| 配置项 | 值 |
|--------|-----|
| 数据库密码 | 见 `/opt/bycigar/.env` |
| JWT Secret | 见 `/opt/bycigar/.env` |
| MinIO 密钥 | minioadmin / minioadmin123 |

**重要经验**：
1. **MySQL 8.4 已移除 `skip-character-set-client-handshake`**，需从 my.cnf 删除
2. **导入数据必须指定字符集**：`mysql --default-character-set=utf8mb4`
3. **数据库端口严禁暴露公网**，必须使用 UFW 限制
4. **定期备份数据库**，并监控容器健康状态
5. **Docker Compose 多文件合并时 volumes 会被覆盖**，生产文件需完整定义
