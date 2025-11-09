# Database Migrations

This directory contains SQL migration files managed by [golang-migrate](https://github.com/golang-migrate/migrate).

## Migration File Naming Convention

Migration files must follow this naming format:

```
{version}_{title}.up.sql
{version}_{title}.down.sql
```

### Restrictions:
- Version numbers should be sequential (000001, 000002, 000003, etc.)
- Each migration needs both `.up.sql` (apply) and `.down.sql` (rollback)
- Filenames must match exactly (same version and title)

## Running Migrations

Migrations are automatically applied when the server starts. They can also be run manually:

### Via Server Startup
The server automatically runs migrations on startup


### Via CLI (if you have golang-migrate installed)
```bash
migrate -path ./backend/migrations -database "postgresql://user:pass@localhost:5432/cm_db?sslmode=disable" up
```

## Migration Status

The `schema_migrations` table tracks which migrations have been applied:

```sql
SELECT * FROM schema_migrations;
```

## Rollback

To rollback the last migration (requires golang-migrate CLI):

```bash
migrate -path ./backend/migrations -database "postgresql://..." down 1
```
