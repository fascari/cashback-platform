# Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for schema migrations.

## Migration files

Migrations live under `db/migrations/{service}/` and follow the naming convention:

```text
{N}_{description}.up.sql    # applies the change
{N}_{description}.down.sql  # reverts the change
```

Migrations are numbered sequentially (`001`, `002`, ...) and are applied in order.

## How migration state is tracked

golang-migrate records applied migrations in a `schema_migrations` table (created automatically on first run):

| Column | Type | Description |
|--------|------|-------------|
| `version` | int | matches the numeric prefix of the migration filename |
| `dirty` | bool | `true` if the last migration failed midway |

A `dirty = true` state means the migration was interrupted. Fix it with `force <version>` before retrying.

## Local dev quickstart

```bash
# 1. Start infrastructure (creates empty databases)
docker-compose up -d postgres

# 2. Apply schemas via migrations (from repo root)
./scripts/migrate.sh cashback-service up
./scripts/migrate.sh blockchain-adapter up
./scripts/migrate.sh mint-consumer up

# Or use mise tasks from each service directory
cd services/cashback-service-api && mise run db:migrate
```

## Running migrations

Use the helper script (requires Docker) from the repo root:

```bash
# Apply all pending migrations
./scripts/migrate.sh cashback-service up

# Roll back the last migration
./scripts/migrate.sh cashback-service down 1

# Check current version
./scripts/migrate.sh cashback-service version

# Clear a dirty state (schema was created outside migrate, or migration was interrupted)
./scripts/migrate.sh cashback-service force 1
```

Or use the `mise` tasks from within each service directory:

```bash
cd services/cashback-service-api
mise run db:migrate    # apply all pending
mise run db:rollback   # roll back last
mise run db:version    # print current version
```

## Writing a new migration

1. Create the next numbered pair in `db/migrations/{service}/`:

```text
db/migrations/cashback-service/
  002_add_cashback_tier.up.sql
  002_add_cashback_tier.down.sql
```

2. The up file applies the change and the down file reverts it exactly.

3. Never modify existing migration files once they have been applied to any environment.

## CI/CD

Run migrations before starting the service. Example in Docker Compose:

```yaml
migrate:
  image: migrate/migrate
  command: ["-path", "/migrations", "-database", "${DB_URL}", "up"]
  volumes:
    - ./db/migrations/cashback-service:/migrations
  depends_on:
    postgres:
      condition: service_healthy
```
