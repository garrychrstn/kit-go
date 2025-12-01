# Database Migration Manager

A simple Go-based database migration tool using `golang-migrate` for PostgreSQL.

## Prerequisites

- Go 1.16+
- PostgreSQL database
- `.env.local` file with `DB_URL` environment variable

## Environment Setup

Create a `.env.local` file in your project root:

```env
DB_URL=postgres://username:password@localhost:5432/database?sslmode=disable
```

## Installation

Install required dependencies:

```bash
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
go get github.com/joho/godotenv
```

## Usage

### Create a New Migration

```bash
go run mgr.go create <migration_name>
```

Example:
```bash
go run mgr.go create add_users_table
```

This creates two files in `./db/migrations/`:
- `000001_add_users_table.up.sql` - Apply migration
- `000001_add_users_table.down.sql` - Rollback migration

### Run Migrations

Apply all pending migrations:
```bash
go run mgr.go do up
```

Rollback all migrations:
```bash
go run mgr.go do down
```

### Force Migration Version

If migrations are in a dirty state, force to a specific version:
```bash
go run mgr.go force <version>
```

Example:
```bash
go run mgr.go force 3
```

### Check Current Version

```bash
go run mgr.go version
```

## Migration File Structure

Migrations are stored in `./db/migrations/` with sequential numbering:

```
db/
└── migrations/
    ├── 000001_create_users.up.sql
    ├── 000001_create_users.down.sql
    ├── 000002_add_posts.up.sql
    └── 000002_add_posts.down.sql
```

## Commands Summary

| Command | Description |
|---------|-------------|
| `create <name>` | Create new migration files |
| `do up` | Apply all pending migrations |
| `do down` | Rollback all migrations |
| `force <version>` | Force database to specific version |
| `version` | Show current migration version |

## Notes

- The tool automatically creates the `./db/migrations/` directory if it doesn't exist
- Migration files are numbered sequentially (000001, 000002, etc.)
- If the database is in a "dirty" state, use `force` to resolve conflicts
