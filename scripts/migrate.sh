#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/migrate.sh <service> <command> [flags]
# Services: cashback-service, blockchain-adapter, mint-consumer
# Commands: up, down, version, force VERSION

SERVICE="${1:?Usage: migrate.sh <service> up|down|version}"
CMD="${2:?Usage: migrate.sh <service> up|down|version}"
shift 2

# host.docker.internal resolves to the host machine from inside a Docker container
# (works on Docker Desktop for Mac/Windows and Linux with --add-host flag).
# Port 15432 is the host-side mapping defined in docker-compose.yml.
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

case "$SERVICE" in
  cashback-service)
    DB_URL="${CASHBACK_DB_URL:-postgres://cashback_app:cashback_app@host.docker.internal:15432/cashback_service_db?sslmode=disable&search_path=cashback}"
    MIGRATIONS_DIR="$REPO_ROOT/services/cashback-service-api/db/migrations"
    ;;
  blockchain-adapter)
    DB_URL="${BLOCKCHAIN_DB_URL:-postgres://cashback_app:cashback_app@host.docker.internal:15432/blockchain_adapter_db?sslmode=disable&search_path=blockchain}"
    MIGRATIONS_DIR="$REPO_ROOT/services/blockchain-adapter/db/migrations"
    ;;
  mint-consumer)
    DB_URL="${MINT_CONSUMER_DB_URL:-postgres://cashback_app:cashback_app@host.docker.internal:15432/mint_consumer_db?sslmode=disable&search_path=mint}"
    MIGRATIONS_DIR="$REPO_ROOT/services/mint-consumer/db/migrations"
    ;;
  *)
    echo "Unknown service: $SERVICE. Valid: cashback-service, blockchain-adapter, mint-consumer" >&2
    exit 1
    ;;
esac

docker run --rm \
  --add-host=host.docker.internal:host-gateway \
  -v "$MIGRATIONS_DIR:/migrations" \
  migrate/migrate \
  -path /migrations \
  -database "$DB_URL" \
  "$CMD" "$@"
