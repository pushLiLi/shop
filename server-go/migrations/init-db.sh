#!/bin/bash
# BYCIGAR 数据库初始化脚本
# 用法: ./init-db.sh [dev|test]
#
# dev: 初始化开发数据库 (bycigar)
# test: 初始化测试数据库 (bycigar_test)

set -e

# 配置
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-123456}"
DB_NAME="${DB_NAME:-bycigar}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 解析参数
ENV=${1:-dev}
if [ "$ENV" = "test" ]; then
    DB_NAME="bycigar_test"
    log_warn "使用测试数据库: $DB_NAME"
fi

# 获取脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

log_info "开始初始化数据库: $DB_NAME"
log_info "数据库地址: $DB_HOST:$DB_PORT"

# 检查MySQL连接
log_info "检查MySQL连接..."
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" &>/dev/null; then
    log_error "无法连接到MySQL，请检查配置"
    exit 1
fi
log_info "MySQL连接正常"

# 创建数据库（如果不存在）
log_info "创建数据库（如果不存在）..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null

# 导入表结构
log_info "导入表结构 (1_schema.sql)..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "1_schema.sql"
log_info "表结构创建完成"

# 导入基础数据
log_info "导入基础数据 (2_base.sql)..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "2_base.sql"
log_info "基础数据导入完成"

# 导入演示数据
log_info "导入演示数据 (3_demo.sql)..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "3_demo.sql"
log_info "演示数据导入完成"

# 导入配置数据
log_info "导入配置数据 (4_settings.sql)..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "4_settings.sql"
log_info "配置数据导入完成"

log_info ""
log_info "=========================================="
log_info "数据库初始化完成！"
log_info "=========================================="
log_info ""
log_info "管理员账号:"
log_info "  Email: admin@bycigar.com"
log_info "  Password: 123456"
log_info ""
log_info "客服账号:"
log_info "  Email: service1@bycigar.com"
log_info "  Password: 123456"
log_info ""
