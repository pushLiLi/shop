#!/bin/bash
# BYCIGAR 数据库重置脚本
# 警告: 此脚本会清空数据库所有数据，请谨慎使用！
# 用法: ./reset-db.sh [dev|test]

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
fi

log_warn "=========================================="
log_warn "警告: 此操作将清空数据库 $DB_NAME 的所有数据！"
log_warn "=========================================="
echo ""

read -p "确认执行? (yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
    log_info "操作已取消"
    exit 0
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

log_info "开始重置数据库: $DB_NAME"

# 导入表结构（会删除并重建所有表）
log_info "重建表结构..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "1_schema.sql"

# 导入种子数据
log_info "导入种子数据..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "2_base.sql"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "3_demo.sql"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "4_settings.sql"

log_info ""
log_info "=========================================="
log_info "数据库重置完成！"
log_info "=========================================="
