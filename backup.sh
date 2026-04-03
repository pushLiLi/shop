#!/bin/bash
# BYCIGAR 数据库备份脚本
# 每日自动备份 MySQL 数据，保留最近 7 天

set -e

BACKUP_DIR="/opt/backups"
MYSQL_CONTAINER="bycigar-mysql"
MYSQL_USER="root"
MYSQL_PASS="Li88613583v+."
DATABASE="bycigar"
KEEP_DAYS=7

mkdir -p "$BACKUP_DIR"

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/bycigar_$DATE.sql"

echo "开始备份: $DATE"

docker exec "$MYSQL_CONTAINER" mysqldump -u "$MYSQL_USER" -p"$MYSQL_PASS" --single-transaction --routines --triggers "$DATABASE" > "$BACKUP_FILE" 2>/dev/null

gzip "$BACKUP_FILE"

echo "备份完成: ${BACKUP_FILE}.gz"
echo "文件大小: $(ls -lh ${BACKUP_FILE}.gz | awk '{print $5}')"

echo "清理 $KEEP_DAYS 天前的旧备份..."
find "$BACKUP_DIR" -name "bycigar_*.sql.gz" -type f -mtime +$KEEP_DAYS -delete

echo "当前备份列表:"
ls -lh "$BACKUP_DIR"/bycigar_*.sql.gz 2>/dev/null || echo "无备份文件"

echo "备份任务完成"
