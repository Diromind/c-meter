# C-Meter Backend

Go backend service for c-meter application with Telegram bot interface.

## Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── config/
│   └── config.go             # Configuration management
├── internal/
│   ├── bot/
│   │   └── handlers.go       # Telegram bot command handlers
│   ├── database/
│   │   ├── database.go       # Database connection & golang-migrate integration
│   │   └── operations.go     # Database operations (CRUD)
│   └── models/               # Data models
├── migrations/               # SQL migration files (golang-migrate format)
├── go.mod                    # Go module dependencies
└── README.md                 # This file
```

## Prerequisites

- Go 1.24.0 or later
- PostgreSQL database
- Telegram Bot Token (get it from [@BotFather](https://t.me/BotFather))

## Configuration

The application uses environment variables for configuration:

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `BOT_TOKEN` | - | **Yes** | Telegram bot token from BotFather |
| `DB_HOST` | `localhost` | No | PostgreSQL host |
| `DB_PORT` | `5432` | No | PostgreSQL port |
| `DB_USER` | `postgres` | No | Database user |
| `DB_PASSWORD` | `postgres` | No | Database password |
| `DB_NAME` | `cm_db` | No | Database name |
| `DB_SSLMODE` | `disable` | No | SSL mode |

## Setup

1. Create a Telegram bot:
   - Message [@BotFather](https://t.me/BotFather) on Telegram
   - Send `/newbot` and follow the instructions
   - Copy the bot token

2. Create the database:
```bash
createdb cm_db
```

3. Install dependencies:
```bash
cd backend
go mod download
```

4. Set up environment variables:
```bash
export BOT_TOKEN="your-bot-token-here"
```

5. Add migration files to the `migrations/` directory following golang-migrate naming convention:
```
000001_initial_schema.up.sql
000001_initial_schema.down.sql
000002_add_feature.up.sql
000002_add_feature.down.sql
...
```

See `migrations/README.md` for detailed migration guide.

## Running

From the backend directory:

```bash
export BOT_TOKEN="your-bot-token-here"
go run cmd/server/main.go
```

Or build and run:

```bash
go build -o server cmd/server/main.go
export BOT_TOKEN="your-bot-token-here"
./server
```

## Telegram Bot Commands

The bot supports the following commands:

### /start
Welcome message and bot introduction.

### /help
Shows all available commands.

### /ping
Health check that returns database connection status and current schema version.

**Example response:**
```
✅ Database connected
📦 Schema version: 01
```

## Development

### Adding New Commands

1. Add handler method in `internal/bot/handlers.go`:
```go
func (h *BotHandler) HandleMyCommand(c tele.Context) error {
    // Your logic here
    return c.Send("Response")
}
```

2. Register the command in `cmd/server/main.go`:
```go
b.Handle("/mycommand", handler.HandleMyCommand)
```

### Adding Database Operations

Add your database queries in `internal/database/operations.go`:

```go
func (db *DB) CreateSomething(data string) error {
    query := `INSERT INTO table_name (column) VALUES ($1)`
    _, err := db.Exec(query, data)
    return err
}
```

### Database Models

Define your data structures in `internal/models/models.go`.

## Migrations

Migrations are managed using [golang-migrate](https://github.com/golang-migrate/migrate) and automatically applied on server startup.

### Migration Files

Migration files must be placed in `migrations/` and follow this naming convention:
```
{version}_{title}.up.sql    # Apply migration
{version}_{title}.down.sql  # Rollback migration
```

Example:
```
000001_initial_schema.up.sql
000001_initial_schema.down.sql
```

### How It Works

- Migrations are automatically applied when the bot starts
- The `schema_migrations` table tracks applied migrations
- Only new migrations are applied (safe to run multiple times)
- Migrations are applied in order by version number
- Each migration runs in a transaction (atomic)

See `migrations/README.md` for detailed migration guide and best practices.

## Architecture

The application follows a layered architecture:

```
Telegram Bot Commands
        ↓
    Handlers (internal/bot)
        ↓
   Operations (internal/database/operations.go)
        ↓
    Database (PostgreSQL)
```

This structure allows for:
- Clean separation of concerns
- Easy testing (can mock the operations layer)
- Future extensibility (can add HTTP API alongside Telegram bot)
