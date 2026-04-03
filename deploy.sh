#!/bin/bash

# BYCIGAR 部署脚本
# 域名: huauge.com / www.huauge.com
# 服务器: 134.209.165.205

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查环境
check_env() {
    if [ ! -f .env.production ]; then
        log_error ".env.production 不存在"
        log_info "请复制 .env.example 并填入生产环境配置"
        exit 1
    fi

    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装"
        exit 1
    fi

    if ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装"
        exit 1
    fi
}

# 拉取最新代码
pull_code() {
    log_info "拉取最新代码..."
    git pull origin main
}

# 构建并启动服务
deploy_services() {
    log_info "构建并启动 Docker 服务..."
    docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
    log_info "等待服务启动..."
    sleep 10
}

# 检查服务状态
check_status() {
    log_info "服务状态:"
    docker compose ps
}

# 执行种子数据（仅首次）
seed_data() {
    SEED_MARKER=".seeded"

    if [ -f "$SEED_MARKER" ]; then
        log_warn "种子数据已执行过，跳过"
        return
    fi

    log_info "执行种子数据..."

    # 等待 MySQL 就绪
    log_info "等待 MySQL 就绪..."
    for i in {1..30}; do
        if docker exec bycigar-mysql mysqladmin ping -h localhost -u root -p"${MYSQL_ROOT_PASSWORD:-123456}" --silent 2>/dev/null; then
            log_info "MySQL 已就绪"
            break
        fi
        sleep 2
    done

    # 执行迁移脚本
    log_info "执行数据库迁移..."
    docker exec -i bycigar-mysql mysql -u root -p"${MYSQL_ROOT_PASSWORD:-123456}" bycigar < server-go/migrations/1_schema.sql 2>/dev/null || true
    docker exec -i bycigar-mysql mysql -u root -p"${MYSQL_ROOT_PASSWORD:-123456}" bycigar < server-go/migrations/2_base.sql 2>/dev/null || true
    docker exec -i bycigar-mysql mysql -u root -p"${MYSQL_ROOT_PASSWORD:-123456}" bycigar < server-go/migrations/3_demo.sql 2>/dev/null || true
    docker exec -i bycigar-mysql mysql -u root -p"${MYSQL_ROOT_PASSWORD:-123456}" bycigar < server-go/migrations/4_settings.sql 2>/dev/null || true

    # 执行种子数据
    log_info "执行演示数据..."
    docker exec -i bycigar-mysql mysql -u root -p"${MYSQL_ROOT_PASSWORD:-123456}" bycigar < server-go/scripts/seed.sql 2>/dev/null || true

    # 创建标记文件
    touch "$SEED_MARKER"
    log_info "种子数据执行完成"
}

# 显示 SSL 配置提示
show_ssl_hint() {
    echo ""
    log_info "=========================================="
    log_info "部署完成!"
    log_info "=========================================="
    echo ""
    log_warn "SSL 证书配置步骤:"
    echo ""
    echo "1. 确保 DNS 已解析:"
    echo "   huauge.com      -> 134.209.165.205"
    echo "   www.huauge.com  -> 134.209.165.205"
    echo ""
    echo "2. 获取 SSL 证书 (首次):"
    echo "   docker compose -f docker-compose.yml -f docker-compose.prod.yml run --rm certbot certonly --webroot --webroot-path /var/www/certbot -d huauge.com -d www.huauge.com --email admin@huauge.com --agree-tos"
    echo ""
    echo "3. 获取证书后重启前端容器加载证书:"
    echo "   docker compose -f docker-compose.yml -f docker-compose.prod.yml restart frontend"
    echo ""
    echo "4. 测试续费:"
    echo "   docker compose -f docker-compose.yml -f docker-compose.prod.yml run --rm certbot renew --dry-run"
    echo ""
    echo "   自动续费已通过 certbot 容器配置，每 12 小时检查一次，无需手动干预"
    echo ""
    log_info "访问地址: https://huauge.com"
    log_info "管理后台: https://huauge.com/admin"
    echo ""
}

main() {
    log_info "开始部署 BYCIGAR..."
    check_env
    pull_code
    deploy_services
    check_status
    seed_data
    show_ssl_hint
}

main "$@"
