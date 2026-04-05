#!/bin/bash
set -e

# Creates the databases and schemas.
# Schemas are needed before migrations run (search_path in DB_URL requires them to exist).
# Table structures are applied by running migrations:
#   ./scripts/migrate.sh cashback-service up
#   ./scripts/migrate.sh blockchain-adapter up
#   ./scripts/migrate.sh mint-consumer up

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres <<-EOSQL
    CREATE DATABASE cashback_service_db;
    CREATE DATABASE mint_consumer_db;
    CREATE DATABASE blockchain_adapter_db;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname cashback_service_db <<-EOSQL
    CREATE SCHEMA IF NOT EXISTS cashback;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname blockchain_adapter_db <<-EOSQL
    CREATE SCHEMA IF NOT EXISTS blockchain;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname mint_consumer_db <<-EOSQL
    CREATE SCHEMA IF NOT EXISTS mint;
EOSQL

echo "Databases and schemas created. Run ./scripts/migrate.sh <service> up to apply table structures."
