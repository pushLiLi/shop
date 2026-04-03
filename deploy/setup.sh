#!/bin/bash
set -e

DOMAIN="huauge.com"
INSTALL_DIR="/opt/bycigar"

echo "=== BYCIGAR Server Setup ==="

if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi

# Install dependencies
apt update
apt install -y nginx certbot python3-certbot-nginx git

# Install Go if not present
if ! command -v go &> /dev/null; then
    echo "Installing Go 1.23..."
    wget -q https://go.dev/dl/go1.23.4.linux-amd64.tar.gz -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    export PATH=$PATH:/usr/local/go/bin
fi

# Create directories
mkdir -p $INSTALL_DIR/uploads $INSTALL_DIR/dist /var/www/certbot

# Create user for running backend
id -u www-data &>/dev/null || useradd -r -s /bin/false www-data
chown -R www-data:www-data $INSTALL_DIR

# Copy nginx config
cp $INSTALL_DIR/deploy/nginx.conf /etc/nginx/sites-available/bycigar
ln -sf /etc/nginx/sites-available/bycigar /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
nginx -t && systemctl enable nginx && systemctl restart nginx

# Copy systemd service
cp $INSTALL_DIR/deploy/bycigar-backend.service /etc/systemd/system/
systemctl daemon-reload

# Start MySQL
cd $INSTALL_DIR
docker compose up -d

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Next steps:"
echo "1. Create $INSTALL_DIR/.env with your config (see deploy/.env.example)"
echo "2. Build backend: cd $INSTALL_DIR/server-go && go build -o ../backend ./cmd/main.go"
echo "3. Build frontend: cd $INSTALL_DIR/bycigar-vue && npm install && npm run build"
echo "4. Start backend: systemctl start bycigar-backend"
echo "5. Get SSL: certbot --nginx -d $DOMAIN -d www.$DOMAIN"
