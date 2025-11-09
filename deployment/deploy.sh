#!/bin/bash

set -e

read -p "Target host (user@host): " TARGET_HOST
if [ -z "$TARGET_HOST" ]; then
    echo "Error: Target host is required"
    exit 1
fi

read -p "Yandex Cloud Lockbox ID: " LOCKBOX_ID
if [ -z "$LOCKBOX_ID" ]; then
    echo "Error: Lockbox ID is required"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

DEB_FILE=$(ls -t "$SCRIPT_DIR"/cmeter_*.deb 2>/dev/null | head -1)
if [ -z "$DEB_FILE" ]; then
    echo "Error: No DEB package found. Run ./build-deb.sh first"
    exit 1
fi

echo ""
echo "=== Fetching secrets from Lockbox ==="

if ! command -v yc &> /dev/null; then
    echo "Error: Yandex Cloud CLI not installed locally"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo "Error: jq not installed. Run: sudo apt-get install -y jq"
    exit 1
fi

PAYLOAD=$(yc lockbox payload get --id "$LOCKBOX_ID" --format json 2>/dev/null)
if [ $? -ne 0 ]; then
    echo "Error: Failed to fetch secrets from Lockbox"
    exit 1
fi

CONNSTRING=$(echo "$PAYLOAD" | jq -r '.entries[] | select(.key=="connstring") | .text_value')
TOKEN=$(echo "$PAYLOAD" | jq -r '.entries[] | select(.key=="token") | .text_value')

if [ -z "$CONNSTRING" ] || [ "$CONNSTRING" == "null" ]; then
    echo "Error: 'connstring' key not found in lockbox"
    exit 1
fi

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo "Error: 'token' key not found in lockbox"
    exit 1
fi

echo "Secrets fetched successfully"
echo ""
echo "=== Deploying $(basename "$DEB_FILE") to $TARGET_HOST ==="
echo ""

scp "$DEB_FILE" "$TARGET_HOST:/tmp/"
ssh "$TARGET_HOST" "sudo dpkg -i /tmp/$(basename "$DEB_FILE") && rm /tmp/$(basename "$DEB_FILE")"

echo ""
echo "=== Writing secrets to target ==="

ssh "$TARGET_HOST" "sudo bash -c 'cat > /opt/cmeter/.env' <<EOF
export DB_CONN_STRING=\"$CONNSTRING\"
export BOT_TOKEN=\"$TOKEN\"
EOF"

ssh "$TARGET_HOST" "sudo chmod 600 /opt/cmeter/.env && sudo chown cmeter:cmeter /opt/cmeter/.env"

ssh "$TARGET_HOST" "sudo supervisorctl reread && sudo supervisorctl update && sudo supervisorctl restart cmeter"

echo ""
echo "=== Deployed successfully ==="
echo "Check: ssh $TARGET_HOST sudo supervisorctl status cmeter"
echo "Logs:  ssh $TARGET_HOST sudo tail -f /var/log/cmeter/output.log"
