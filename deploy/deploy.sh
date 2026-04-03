#!/bin/bash
set -e

INSTALL_DIR="/opt/bycigar"
cd $INSTALL_DIR

echo "=== Pulling latest code ==="
git pull

echo "=== Building backend ==="
cd server-go
/usr/local/go/bin/go build -o ../backend ./cmd/main.go

echo "=== Building frontend ==="
cd ../bycigar-vue
npm install --production=false
npm run build
cp -r dist/* ../dist/

echo "=== Restarting backend ==="
sudo systemctl restart bycigar-backend
sudo systemctl reload nginx

echo "=== Deploy Complete ==="
systemctl status bycigar-backend --no-pager -l
