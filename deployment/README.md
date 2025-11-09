# C-Meter Deployment

## Prerequisites

**Build/Deploy machine:**
- Go 1.24+
- Yandex Cloud CLI: `curl https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash && yc init`
- `jq`: `sudo apt-get install -y jq`

**Target host:**
- Ubuntu/Debian with `supervisor` installed

**Yandex Cloud Lockbox:**
- `connstring` - PostgreSQL connection string
- `token` - Telegram bot token

## Deploy

```bash
cd deployment

# Build DEB package
./build-deb.sh

# Deploy to target (fetches secrets locally, pushes to target)
./deploy.sh
# Enter: user@host, lockbox-id
```

## Manage

```bash
# Status
ssh user@host sudo supervisorctl status cmeter

# Restart
ssh user@host sudo supervisorctl restart cmeter

# Logs
ssh user@host sudo tail -f /var/log/cmeter/output.log

# Uninstall
ssh user@host "sudo supervisorctl stop cmeter && sudo apt remove cmeter"
```
