# MinIO 图片存储迁移方案

## 一、总体方案

| 项目 | 方案 |
|------|------|
| **存储后端** | MinIO（Docker 部署，API:9000, Console:9001） |
| **读取方式** | 前端通过 nginx/Vite 反代 `/media` 路径访问 MinIO，Bucket 设为公开读取 |
| **迁移范围** | 仅 `AdminImageUpload` 上传的图片（product/banner/config），favicon 等静态资源保留本地 |
| **URL 格式** | `/media/bycigar/{timestamp}_{uuid}{ext}`（相对路径，环境无关） |
| **反代策略** | 本地开发：Vite proxy `/media` → `localhost:9000`；Docker 生产：nginx proxy `/media` → `minio:9000` |

### 数据流

```
上传流程：
管理后台 → Go 后端 → 存文件到 MinIO → 返回相对路径 /media/bycigar/xxx.jpg → 存入数据库

展示流程：
前端 <img src="/media/bycigar/xxx.jpg"> → nginx/Vite 反代 → MinIO 返回图片
```

Go 后端只负责"写入"（上传），"读取"由前端通过反代直接访问 MinIO，Go 不做中转。

---

## 二、MinIO 配置（docker-compose.yml）

在 `docker-compose.yml` 中新增 MinIO 服务：

```yaml
minio:
  image: minio/minio:latest
  container_name: bycigar-minio
  ports:
    - "9000:9000"
    - "9001:9001"
  environment:
    MINIO_ROOT_USER: minioadmin
    MINIO_ROOT_PASSWORD: minioadmin123
  volumes:
    - minio_data:/data
  command: server /data --console-address ":9001"
  healthcheck:
    test: ["CMD", "mc", "ready", "local"]
    interval: 5s
    timeout: 5s
    retries: 10
```

在 `volumes` 部分新增 `minio_data`。

> 端口 9000/9001 暴露到宿主机，方便本地开发时直接访问 MinIO Console。
> Docker 生产环境中，前端不直连 9000 端口，全部通过 nginx 反代 `/media` 访问。

---

## 三、后端改动（server-go/）

### 3.1 新增依赖

```bash
cd server-go
go get github.com/minio/minio-go/v7
```

### 3.2 配置项扩展（internal/config/config.go）

`Config` 结构体新增字段：

```go
MinioEndpoint  string  // 默认 "localhost:9000"（Docker 中为 "minio:9000"）
MinioAccessKey string  // 默认 "minioadmin"
MinioSecretKey string  // 默认 "minioadmin123"
MinioBucket    string  // 默认 "bycigar"
MinioUseSSL    bool    // 默认 false
```

对应 `.env` 环境变量和 `docker-compose.yml` backend 的 environment。

> 不需要 `MinioPublicURL` 配置项。上传接口返回相对路径 `/media/bycigar/xxx`，
> 前端通过 nginx/Vite 反代访问，无需按环境切换 URL 前缀。

### 3.3 新建 MinIO 客户端包（pkg/minio/minio.go）

职责：
- 初始化 MinIO client 连接
- `InitMinio()` — 创建连接，暴露全局 `Client` 和 `Bucket`
- `EnsureBucket()` — 启动时检查 bucket 是否存在，不存在则创建并设为公开读取策略（anonymous read）

### 3.4 改写上传接口（internal/handlers/upload.go）

**当前流程：** `os.Create` → `io.Copy` → 返回 `/static/media/xxx`

**新流程：**
1. 从 FormFile 读取文件（不变）
2. 校验格式和大小（不变）
3. 生成文件名（不变：`{timestamp}_{uuid}{ext}`）
4. 调用 `minio.Client.PutObject` 上传到 MinIO bucket
5. 返回相对路径：`/media/{bucket}/{filename}`

返回 JSON：
```json
{"url": "/media/bycigar/1774680065_ec565b5c-xxx.jpg", "success": true}
```

### 3.5 修改 main.go

- 在启动流程中加入 `minio.InitMinio()` 和 `minio.EnsureBucket()`（在 DB migrate 之后）
- 保留 `r.Static("/static", "./static")`（favicon 等本地静态资源仍需服务）

### 3.6 更新种子数据（internal/database/database.go）

```go
// 旧：
ConfigValue: "/static/media/banner-1.png"
// 新：
ConfigValue: "/media/bycigar/banner-1.png"
```

### 3.7 数据迁移工具（cmd/migrate_images/main.go）

一次性迁移工具，执行流程：
1. 连接 MinIO 和 MySQL
2. 扫描 `static/media/` 目录中的图片文件
3. 逐个上传到 MinIO bucket
4. 查询数据库中所有含 `/static/media/` 路径的记录：
   - `products` 表的 `image` 字段
   - `banners` 表的 `image` 字段
   - `site_configs` 表的 `config_value` 字段（`home_banner_*` 相关）
5. 将 URL 从 `/static/media/xxx` 更新为 `/media/bycigar/xxx`

使用方式：
```bash
cd server-go
go run ./cmd/migrate_images/main.go
```

---

## 四、前端改动（bycigar-vue/）

### 4.1 AdminImageUpload.vue

**无需改动核心逻辑。** 组件只负责将 `data.url` 赋值给 `imageUrl`，后端返回的相对路径 `/media/bycigar/xxx` 通过 `<img :src="imageUrl">` 可直接加载。

### 4.2 HomeView.vue 硬编码 fallback

```javascript
// 旧：
const fallbackBanners = [
  '/static/media/微信图片_20260303152810_1_341(2).jpg',
  '/static/media/banner-4.png',
  '/static/media/banner-5.jpg'
]
// 新：
const fallbackBanners = [
  '/media/bycigar/微信图片_20260303152810_1_341(2).jpg',
  '/media/bycigar/banner-4.png',
  '/media/bycigar/banner-5.jpg'
]
```

### 4.3 index.html favicon

将 favicon 从 `/static/media/favicon.png` 改为 `/favicon.png`，并将文件复制到 `bycigar-vue/public/favicon.png`。

### 4.4 Vite proxy 配置（vite.config.js）

将 `/static` 代理改为 `/media` 代理，指向本地 MinIO：

```javascript
proxy: {
  '/api': {
    target: 'http://localhost:3000',
    changeOrigin: true
  },
  '/media': {
    target: 'http://localhost:9000',
    changeOrigin: true
  }
}
```

### 4.5 nginx.conf（生产环境）

新增 `/media` 反代到 MinIO：

```nginx
server {
    listen 80;
    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    location /api {
        proxy_pass http://backend:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
    location /media/ {
        proxy_pass http://minio:9000/;
        proxy_set_header Host $host;
    }
}
```

> Docker 环境中 nginx 和 minio 在同一网络，使用 Docker 服务名 `minio` 连接。

---

## 五、需要注意的问题

### 5.1 跨域问题

MinIO bucket 设为公开读取后，图片 GET 请求通过 nginx 同域反代访问，不存在跨域问题。

### 5.2 环境差异（本地开发 vs Docker 生产）

| 环境 | 前端访问图片 | 后端连接 MinIO | MinIO 内部地址 |
|------|-------------|---------------|---------------|
| 本地开发 | `/media/xxx` → Vite proxy → `localhost:9000` | `localhost:9000` | - |
| Docker 生产 | `/media/xxx` → nginx proxy → `minio:9000` | `minio:9000` | Docker 网络 |

后端 `MINIO_ENDPOINT` 通过环境变量区分：
- `.env` 本地开发：`MINIO_ENDPOINT=localhost:9000`
- `docker-compose.yml` backend environment：`MINIO_ENDPOINT=minio:9000`

### 5.3 删除图片时的清理（可选增强）

当前删除 product/banner 不删除文件。迁移到 MinIO 后可在删除 handler 中调用 `minio.RemoveObject` 清理文件，此为可选优化，不在本次迁移必须范围内。

### 5.4 保留本地静态文件服务

`r.Static("/static", "./static")` 保留，用于 favicon 和其他非上传静态资源。上传的图片完全走 MinIO。

---

## 六、实施步骤顺序

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | docker-compose.yml 添加 MinIO 服务 | `docker-compose.yml` |
| 2 | 启动 MinIO 容器验证 | - |
| 3 | Go 依赖安装 minio-go | `server-go/go.mod` |
| 4 | config 新增 MinIO 配置项（无 PublicURL） | `server-go/internal/config/config.go`, `server-go/.env` |
| 5 | 新建 MinIO 客户端包 | `server-go/pkg/minio/minio.go` |
| 6 | main.go 加入 MinIO 初始化 | `server-go/cmd/main.go` |
| 7 | 改写 upload handler（返回相对路径） | `server-go/internal/handlers/upload.go` |
| 8 | 更新种子数据 URL 为相对路径 | `server-go/internal/database/database.go` |
| 9 | 后端编译验证 | `go build` |
| 10 | 启动后端验证上传功能 | - |
| 11 | 编写并执行数据迁移脚本 | `server-go/cmd/migrate_images/main.go` |
| 12 | Vite proxy 改为 `/media` → `localhost:9000` | `bycigar-vue/vite.config.js` |
| 13 | nginx.conf 新增 `/media` 反代 | `bycigar-vue/nginx.conf` |
| 14 | 前端适配（fallback URL、favicon） | `HomeView.vue`, `index.html` |
| 15 | 前端编译验证 | `npm run build` |
| 16 | 更新 AGENTS.md | `AGENTS.md` |

---

## 七、影响范围汇总

### 后端文件

| 文件 | 改动类型 | 说明 |
|------|---------|------|
| `docker-compose.yml` | 修改 | 添加 MinIO 服务 + volume |
| `server-go/.env` | 修改 | 添加 MinIO 环境变量 |
| `server-go/internal/config/config.go` | 修改 | 添加 MinIO 配置字段（无 PublicURL） |
| `server-go/pkg/minio/minio.go` | **新建** | MinIO 客户端封装 |
| `server-go/cmd/main.go` | 修改 | 初始化 MinIO |
| `server-go/internal/handlers/upload.go` | **重写** | 上传到 MinIO，返回相对路径 |
| `server-go/internal/database/database.go` | 修改 | 种子数据 URL 改为 `/media/bycigar/xxx` |
| `server-go/cmd/migrate_images/main.go` | **新建** | 一次性迁移工具 |

### 前端文件

| 文件 | 改动类型 | 说明 |
|------|---------|------|
| `bycigar-vue/vite.config.js` | 修改 | `/static` 代理改为 `/media` → `localhost:9000` |
| `bycigar-vue/nginx.conf` | 修改 | 新增 `location /media/` 反代到 MinIO |
| `bycigar-vue/src/views/HomeView.vue` | 修改 | fallback banner URL 改为 `/media/bycigar/xxx` |
| `bycigar-vue/index.html` | 修改 | favicon 路径改为 `/favicon.png` |

### 不需要改动的文件

| 文件 | 原因 |
|------|------|
| `AdminImageUpload.vue` | 只使用 `data.url` 赋值，自动适配新 URL |
| `AdminBanners.vue` | `form.imageUrl` 存储/展示相对 URL，无需改动 |
| `AdminProducts.vue` | 同上 |
| `AdminSettings.vue` | 同上 |
| `ProductCard.vue` 等 | `<img :src="product.imageUrl">` 支持相对路径 |

---

## 八、现有文件清单（需要迁移到 MinIO 的文件）

`server-go/static/media/` 目录中的上传图片：

| 文件 | 说明 |
|------|------|
| `banner-1.png` | 首页横幅种子图 |
| `banner-2.png` | 首页横幅种子图 |
| `banner-3.png` | 首页横幅种子图 |
| `banner-4.png` | fallback 横幅图 |
| `banner-5.jpg` | fallback 横幅图 |
| `favicon.png` | 网站图标（**不迁移**） |
| `微信图片_*.jpg` | fallback 横幅图 |
| `33sq063kpkcrw2l1746422237.png` | 商品图片 |
| `6lh38qdoltluioq1755850965-600x600.jpg` | 商品图片 |
| `ScreenShot_*.png` | 商品图片 |
| `1774680065_*.jpg` | 商品图片（修复后上传） |
| `1774680286_*.png` | 商品图片（修复后上传） |
| `1774*_*.jpg`（含 %!s(MISSING)） | 旧 bug 产生的损坏文件（迁移时跳过或手动处理） |
| `产品/` | 中文目录（检查内容后决定是否迁移） |
