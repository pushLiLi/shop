#!/bin/bash
# BYCIGAR 服务器初始化脚本
# 执行方式: ./server-setup.sh 或 ./server-setup.sh --security-only

set -e

SECURITY_ONLY=false
if [ "$1" == "--security-only" ]; then
    SECURITY_ONLY=true
fi

security_hardening() {
    echo "===== 安全加固开始 ====="

    echo "--- 1. 配置防火墙 (ufw) ---"
    apt-get update
    apt-get install -y ufw
    ufw --force reset
    ufw default deny incoming
    ufw default allow outgoing
    ufw allow 22/tcp
    ufw allow 80/tcp
    ufw allow 443/tcp
    ufw --force enable
    echo "防火墙已启用，只开放 22, 80, 443 端口"

    echo "--- 2. 检查可疑数据库 ---"
    if docker ps | grep -q bycigar-mysql; then
        SUSPICIOUS_DB=$(docker exec bycigar-mysql mysql -u root -pLi88613583v+. -N -e "SHOW DATABASES LIKE 'RECOVER_YOUR_DATA';" 2>/dev/null || echo "")
        if [ -n "$SUSPICIOUS_DB" ]; then
            echo "发现可疑数据库 RECOVER_YOUR_DATA，正在删除..."
            docker exec bycigar-mysql mysql -u root -pLi88613583v+. -e "DROP DATABASE RECOVER_YOUR_DATA;" 2>/dev/null
            echo "可疑数据库已删除"
        else
            echo "未发现可疑数据库"
        fi
    else
        echo "MySQL 容器未运行，跳过数据库检查"
    fi

    echo "--- 3. 检查 MySQL 端口暴露 ---"
    if ss -tlnp | grep -q ":3306.*0.0.0.0"; then
        echo "警告: MySQL 3306 端口对外暴露，建议检查 docker-compose 配置"
    else
        echo "MySQL 端口安全（未对外暴露）"
    fi

    echo "--- 4. 修复 my.cnf 权限 ---"
    if [ -f /opt/bycigar/mysql/my.cnf ]; then
        chmod 644 /opt/bycigar/mysql/my.cnf
        echo "my.cnf 权限已修复为 644"
    fi

    echo "--- 5. 设置备份脚本定时任务 ---"
    if [ -f /opt/bycigar/backup.sh ]; then
        chmod +x /opt/bycigar/backup.sh
        if ! crontab -l 2>/dev/null | grep -q "backup.sh"; then
            (crontab -l 2>/dev/null; echo "0 3 * * * /opt/bycigar/backup.sh >> /var/log/bycigar-backup.log 2>&1") | crontab -
            echo "每日 3:00 自动备份已设置"
        else
            echo "备份定时任务已存在"
        fi
    fi

    echo "===== 安全加固完成 ====="
}

if [ "$SECURITY_ONLY" = true ]; then
    security_hardening
    exit 0
fi

echo "===== 1. 安装 Docker ====="
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com | sh
    usermod -aG docker root
    echo "Docker 安装完成"
else
    echo "Docker 已安装"
fi

echo "===== 2. 安装 Docker Compose ====="
if ! docker compose version &> /dev/null; then
    apt-get update
    apt-get install -y docker-compose-plugin
    echo "Docker Compose 安装完成"
else
    echo "Docker Compose 已安装"
fi

echo "===== 3. Certbot (已改用 Docker 容器管理) ====="
echo "Certbot 不再需要安装在宿主机上"
echo "SSL 证书的获取和自动续费均通过 docker-compose.prod.yml 中的 certbot 容器处理"

echo "===== 4. 克隆代码 ====="
mkdir -p /opt/bycigar
cd /opt/bycigar

if [ ! -d ".git" ]; then
    git clone https://github.com/pushLiLi/shop.git .
    echo "代码克隆完成"
else
    git pull origin main
    echo "代码更新完成"
fi

echo "===== 5. 创建环境配置 ====="
cat > /opt/bycigar/.env.production << 'ENVEOF'
DOMAIN=huauge.com

MYSQL_ROOT_PASSWORD=Byc1gar@Pr0d!2024
DB_PASSWORD=Byc1gar@Pr0d!2024

JWT_SECRET=hUaUgE_s3cr3t_k3y_2024_r@nd0m_str1ng

MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=M1n10@Pr0d!2024
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=M1n10@Pr0d!2024
MINIO_BUCKET=bycigar
ENVEOF

echo ".env.production 创建完成"

echo "===== 6. 运行部署 ====="
chmod +x /opt/bycigar/deploy.sh
cd /opt/bycigar
./deploy.sh

echo "===== 部署完成 ====="
echo "下一步: 配置 SSL 证书"
echo "执行: docker compose -f docker-compose.yml -f docker-compose.prod.yml run --rm certbot certonly --webroot --webroot-path /var/www/certbot -d huauge.com -d www.huauge.com --email admin@huauge.com --agree-tos"
echo "然后: docker compose -f docker-compose.yml -f docker-compose.prod.yml restart frontend"
