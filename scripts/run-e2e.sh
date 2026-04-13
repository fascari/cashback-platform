#!/usr/bin/env bash
# run-e2e.sh — orchestrates the full e2e lifecycle:
#   1. Start isolated infrastructure (postgres, nats, redis)
#   2. Run migrations
#   3. Start services natively
#   4. Run tests
#   5. Tear everything down
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
COMPOSE="docker-compose -f $REPO_ROOT/docker-compose.e2e.yml"

# Load base env then override infra connection strings for e2e.
set -a
# shellcheck source=../.env
source "$REPO_ROOT/.env"
# shellcheck source=../.env.e2e
source "$REPO_ROOT/.env.e2e"
set +a

ADAPTER_PID=""
API_PID=""
CONSUMER_PID=""

cleanup() {
  echo "Stopping e2e services..."
  [[ -n "$ADAPTER_PID" ]] && kill "$ADAPTER_PID" 2>/dev/null || true
  [[ -n "$API_PID" ]]     && kill "$API_PID"     2>/dev/null || true
  [[ -n "$CONSUMER_PID" ]] && kill "$CONSUMER_PID" 2>/dev/null || true
  echo "Stopping e2e infrastructure..."
  $COMPOSE down -v
}
trap cleanup EXIT

# Kill any leftover processes from a previous interrupted run.
lsof -ti:18080 2>/dev/null | xargs kill 2>/dev/null || true

echo "Starting e2e infrastructure..."
$COMPOSE up -d

echo "Waiting for e2e PostgreSQL..."
until $COMPOSE exec -T postgres-e2e pg_isready -U cashback_app 2>/dev/null; do
  sleep 2
done

echo "Waiting for e2e NATS..."
until curl -sf http://localhost:8322/healthz >/dev/null 2>&1; do
  sleep 2
done

echo "Running migrations..."
mise run db:migrate -- --e2e

echo "Starting e2e services..."
(cd "$REPO_ROOT/services/blockchain-adapter"   && go run ./cmd/main.go) &
ADAPTER_PID=$!
(cd "$REPO_ROOT/services/cashback-service-api" && go run ./cmd/api/main.go) &
API_PID=$!
(cd "$REPO_ROOT/services/mint-consumer"        && go run ./cmd/main.go) &
CONSUMER_PID=$!

echo "Waiting for e2e API..."
until curl -sf http://localhost:18080/health >/dev/null 2>&1; do
  echo "  still waiting..."
  sleep 3
done

echo "Running e2e tests..."
cd "$REPO_ROOT/test/e2e" && GOWORK=off go test -v -tags=e2e -timeout=120s ./...
