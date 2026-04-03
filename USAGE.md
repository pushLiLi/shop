# BYCIGAR 运维手册

> 服务器：167.99.115.47
> 域名：huauge.com

## 访问地址

| 入口 | 地址 |
|------|------|
| 前台 | https://huauge.com |
| 管理后台 | https://huauge.com/admin |
| MinIO 控制台 | http://167.99.115.47:9001 |

## 账号

| 账号 | 密码 | 角色 |
|------|------|------|
| admin@bycigar.com | 123456 | 管理员 |
| service1@bycigar.com | 123456 | 客服 |
| user01@test.com | 123456 | 客户 |

## 常用命令

```bash
cd /opt/bycigar

# 查看状态
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps

# 查看日志
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f backend

# 重启服务
docker compose -f docker-compose.yml -f docker-compose.prod.yml restart

# 备份数据库
docker exec bycigar-mysql mysqldump -uroot -p'PASSWORD' bycigar | gzip > /opt/backups/backup_$(date +%Y%m%d).sql.gz

# 健康检查
curl https://huauge.com/api/health

# SSL 续期
certbot renew
```

## 服务状态

- MySQL: 3306（仅容器内部）
- MinIO: 9000/9001（仅容器内部）
- Backend: 3000（仅容器内部）
- Frontend: 80/443（对外）

## 关键注意事项

### 字符集
导入数据**必须**指定字符集：
```bash
docker exec -i bycigar-mysql mysql --default-character-set=utf8mb4 -uroot -p'PASSWORD' bycigar < data.sql
```

### 安全
- 数据库端口 3306 **严禁**暴露公网
- 必须使用强密码
- 定期检查 `/opt/backups/` 确保备份正常

### 端口
```bash
ufw status  # 应只开放 22, 80, 443
```
