#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_ENV="$SCRIPT_DIR/../../.env"

output=$(npm run deploy:sepolia 2>&1)
echo "$output"

address=$(echo "$output" | grep "CashbackToken deployed to:" | awk '{print $NF}')
if [ -n "$address" ]; then
  if grep -q "^CONTRACT_ADDRESS=" "$ROOT_ENV" 2>/dev/null; then
    sed -i.bak "s|^CONTRACT_ADDRESS=.*|CONTRACT_ADDRESS=$address|" "$ROOT_ENV" && rm -f "$ROOT_ENV.bak"
  else
    echo "CONTRACT_ADDRESS=$address" >> "$ROOT_ENV"
  fi
  echo ""
  echo "CONTRACT_ADDRESS updated in .env: $address"
fi
