#!/usr/bin/env bash
# run-e2e.sh — orchestrates the full e2e lifecycle:
#   1. Start isolated infrastructure (postgres, nats, redis)
#   2. Run migrations
#   3. Build and start services natively
#   4. Run tests
#   5. Tear everything down
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
COMPOSE="docker-compose -f $REPO_ROOT/docker-compose.e2e.yml"
BIN_DIR="$REPO_ROOT/.bin/e2e"

# Load base env then override infra connection strings for e2e.
set -a
# shellcheck source=../.env
source "$REPO_ROOT/.env"
# shellcheck source=../.e2e.env
source "$REPO_ROOT/.e2e.env"
set +a

ADAPTER_PID=""
API_PID=""
CONSUMER_PID=""

cleanup() {
  echo "Stopping e2e services..."
  [[ -n "$ADAPTER_PID" ]] && kill "$ADAPTER_PID" 2>/dev/null || true
  [[ -n "$API_PID" ]]     && kill "$API_PID"     2>/dev/null || true
  [[ -n "$CONSUMER_PID" ]] && kill "$CONSUMER_PID" 2>/dev/null || true
  wait "$ADAPTER_PID" "$API_PID" "$CONSUMER_PID" 2>/dev/null || true
  echo "Stopping e2e infrastructure..."
  $COMPOSE down -v
}
trap cleanup EXIT

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

echo "Waiting for e2e Anvil..."
until curl -sf http://localhost:8545 -X POST \
    -H 'Content-Type: application/json' \
    -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' >/dev/null 2>&1; do
  sleep 2
done

echo "Deploying CashbackToken to Anvil..."
(cd "$REPO_ROOT/contracts" && npm run deploy:local >/dev/null)

echo "Setting up NATS streams..."
(cd "$REPO_ROOT/cmd/nats-setup" && NATS_URL=nats://localhost:4322 go run .)

echo "Running migrations..."
mise run db:migrate -- --e2e

echo "Building e2e service binaries..."
mkdir -p "$BIN_DIR"
(cd "$REPO_ROOT/services/blockchain-adapter"   && go build -o "$BIN_DIR/blockchain-adapter" ./cmd/main.go)
(cd "$REPO_ROOT/services/cashback-service-api" && go build -o "$BIN_DIR/cashback-service-api" ./cmd/api/main.go)
(cd "$REPO_ROOT/services/mint-consumer"        && go build -o "$BIN_DIR/mint-consumer" ./cmd/main.go)

echo "Starting e2e services..."
"$BIN_DIR/blockchain-adapter" &
ADAPTER_PID=$!
"$BIN_DIR/cashback-service-api" &
API_PID=$!
"$BIN_DIR/mint-consumer" &
CONSUMER_PID=$!

echo "Waiting for e2e API..."
until curl -sf http://localhost:18080/health >/dev/null 2>&1; do
  echo "  still waiting..."
  sleep 3
done

echo "Waiting for e2e blockchain-adapter (gRPC:15051)..."
until nc -z localhost 15051 2>/dev/null; do
  echo "  still waiting..."
  sleep 2
done

echo "Running e2e tests..."
cd "$REPO_ROOT/test/e2e" && GOWORK=off go test -v -tags=e2e -timeout=120s ./...
