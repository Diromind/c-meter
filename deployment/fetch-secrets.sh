#!/bin/bash

set -e

if [ -z "$LOCKBOX_ID" ]; then
    echo "Error: LOCKBOX_ID is required"
    echo "Usage: sudo LOCKBOX_ID=<id> ./fetch-secrets.sh"
    exit 1
fi

if ! command -v yc &> /dev/null; then
    echo "Error: Yandex Cloud CLI not installed"
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

cat > /opt/cmeter/.env <<EOF
export DB_CONN_STRING="$CONNSTRING"
export BOT_TOKEN="$TOKEN"
EOF

chmod 600 /opt/cmeter/.env
chown cmeter:cmeter /opt/cmeter/.env 2>/dev/null || true

echo "Secrets configured successfully"
