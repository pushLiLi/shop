#!/bin/bash
set -euo pipefail

BASE_DIR="/opt/bycigar"
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

need_frontend=false
need_backend=false

check_changes() {
    git fetch
    local_head=$(git rev-parse HEAD)
    remote_head=$(git rev-parse @{u})

    if [ "$local_head" = "$remote_head" ]; then
        log_info "代码无变动，退出"
        exit 0
    fi

    log_info "检测到新提交，拉取..."
    git pull --ff-only

    # 检测前后端目录变动
    if ! git diff --quiet "$local_head" "$remote_head" -- bycigar-vue/; then
        need_frontend=true
        log_info "前端有变动"
    fi

    if ! git diff --quiet "$local_head" "$remote_head" -- server-go/; then
        need_backend=true
        log_info "后端有变动"
    fi
}

deploy_frontend() {
    if [ "$need_frontend" = false ]; then
        return
    fi
    log_info "构建前端..."
    (cd "$BASE_DIR/bycigar-vue" && npm ci && npm run build) || {
        log_error "前端构建失败"
        exit 1
    }
    cp -rf "$BASE_DIR/bycigar-vue/dist/"* "$BASE_DIR/dist/"
    log_info "重载 nginx"
    cp "$BASE_DIR/deploy/nginx.conf" /etc/nginx/sites-available/bycigar
    nginx -t && systemctl reload nginx
}

deploy_backend() {
    if [ "$need_backend" = false ]; then
        return
    fi
    log_info "构建后端..."
    (cd "$BASE_DIR/server-go" && go mod download && CGO_ENABLED=0 go build -o "$BASE_DIR/backend" ./cmd/main.go) || {
        log_error "后端构建失败"
        exit 1
    }
    systemctl restart bycigar-backend
    sleep 2
    if ! systemctl is-active --quiet bycigar-backend; then
        log_error "后端启动失败，查看日志"
        journalctl -u bycigar-backend -n 20 --no-pager
        exit 1
    fi
}

main() {
    cd "$BASE_DIR"
    check_changes
    deploy_frontend
    deploy_backend
    log_info "部署完成！"
}

main "$@"