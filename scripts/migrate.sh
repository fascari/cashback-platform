#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/migrate.sh <service> <command> [--e2e] [flags]
# Services: cashback-service, blockchain-adapter, mint-consumer
# Commands: up, down, version, force VERSION
# --e2e    Target the e2e database (port 25432) instead of dev (port 15432)

SERVICE="${1:?Usage: migrate.sh <service> up|down|version [--e2e]}"
CMD="${2:?Usage: migrate.sh <service> up|down|version [--e2e]}"
shift 2

PG_PORT=15432
EXTRA_ARGS=()
for arg in "$@"; do
  if [[ "$arg" == "--e2e" ]]; then
    PG_PORT=25432
  else
    EXTRA_ARGS+=("$arg")
  fi
done

# host.docker.internal resolves to the host machine from inside a Docker container
# (works on Docker Desktop for Mac/Windows and Linux with --add-host flag).
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

case "$SERVICE" in
  cashback-service)
    DB_URL="${CASHBACK_DB_URL:-postgres://cashback_app:cashback_app@host.docker.internal:${PG_PORT}/cashback_service_db?sslmode=disable&search_path=cashback}"
    MIGRATIONS_DIR="$REPO_ROOT/services/cashback-service-api/db/migrations"
    ;;
  blockchain-adapter)
    DB_URL="${BLOCKCHAIN_DB_URL:-postgres://cashback_app:cashback_app@host.docker.internal:${PG_PORT}/blockchain_adapter_db?sslmode=disable&search_path=blockchain}"
    MIGRATIONS_DIR="$REPO_ROOT/services/blockchain-adapter/db/migrations"
    ;;
  mint-consumer)
    DB_URL="${MINT_CONSUMER_DB_URL:-postgres://cashback_app:cashback_app@host.docker.internal:${PG_PORT}/mint_consumer_db?sslmode=disable&search_path=mint}"
    MIGRATIONS_DIR="$REPO_ROOT/services/mint-consumer/db/migrations"
    ;;
  *)
    echo "Unknown service: $SERVICE. Valid: cashback-service, blockchain-adapter, mint-consumer" >&2
    exit 1
    ;;
esac

MIGRATE_ARGS=("$CMD")
[[ ${#EXTRA_ARGS[@]} -gt 0 ]] && MIGRATE_ARGS+=("${EXTRA_ARGS[@]}")

docker run --rm \
  --add-host=host.docker.internal:host-gateway \
  -v "$MIGRATIONS_DIR:/migrations" \
  migrate/migrate \
  -path /migrations \
  -database "$DB_URL" \
  "${MIGRATE_ARGS[@]}"
