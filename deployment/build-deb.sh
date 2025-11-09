#!/bin/bash

set -e

VERSION=$(date +%Y%m%d.%H%M%S)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
ARCH="amd64"
PKG_NAME="cmeter"
BUILD_DIR="/tmp/${PKG_NAME}_${VERSION}_${ARCH}"

echo "=== Building c-meter v${VERSION} ==="

cd "$REPO_ROOT/backend"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmeter-bot cmd/server/main.go

rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR/DEBIAN"
mkdir -p "$BUILD_DIR/opt/cmeter"
mkdir -p "$BUILD_DIR/etc/supervisor/conf.d"
mkdir -p "$BUILD_DIR/var/log/cmeter"

cp cmeter-bot "$BUILD_DIR/opt/cmeter/"
cp -r migrations "$BUILD_DIR/opt/cmeter/"
cp "$SCRIPT_DIR/cmeter.conf" "$BUILD_DIR/etc/supervisor/conf.d/"
cp "$SCRIPT_DIR/fetch-secrets.sh" "$BUILD_DIR/opt/cmeter/"
chmod +x "$BUILD_DIR/opt/cmeter/fetch-secrets.sh"

cat > "$BUILD_DIR/DEBIAN/control" <<EOF
Package: cmeter
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Depends: supervisor
Maintainer: C-Meter Team
Description: C-Meter calorie tracking Telegram bot
EOF

cat > "$BUILD_DIR/DEBIAN/postinst" <<'EOF'
#!/bin/bash
set -e

useradd -r -s /bin/false cmeter 2>/dev/null || true
chown -R cmeter:cmeter /opt/cmeter /var/log/cmeter
chmod +x /opt/cmeter/cmeter-bot
chmod +x /opt/cmeter/fetch-secrets.sh

if [ ! -f /opt/cmeter/.env ]; then
    touch /opt/cmeter/.env
    chmod 600 /opt/cmeter/.env
    chown cmeter:cmeter /opt/cmeter/.env
fi

if command -v supervisorctl >/dev/null 2>&1; then
    supervisorctl reread || true
    supervisorctl update || true
fi
EOF

chmod +x "$BUILD_DIR/DEBIAN/postinst"

cat > "$BUILD_DIR/DEBIAN/prerm" <<'EOF'
#!/bin/bash
set -e
if command -v supervisorctl >/dev/null 2>&1; then
    supervisorctl stop cmeter || true
fi
EOF

chmod +x "$BUILD_DIR/DEBIAN/prerm"

dpkg-deb --build "$BUILD_DIR"
mv "${BUILD_DIR}.deb" "$SCRIPT_DIR/${PKG_NAME}_${VERSION}_${ARCH}.deb"

echo ""
echo "=== Package built: ${PKG_NAME}_${VERSION}_${ARCH}.deb ==="
