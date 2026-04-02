#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres <<-EOSQL
    CREATE DATABASE cashback_service_db;
    CREATE DATABASE mint_consumer_db;
    CREATE DATABASE blockchain_adapter_db;
EOSQL

echo "Applying schema to cashback_service_db..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname cashback_service_db -f /db/cashback_service.sql

echo "Applying schema to mint_consumer_db..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname mint_consumer_db -f /db/mint_consumer.sql

echo "Applying schema to blockchain_adapter_db..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname blockchain_adapter_db -f /db/blockchain_adapter.sql

echo "Database initialization complete."
