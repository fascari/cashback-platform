#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres <<-EOSQL
    CREATE DATABASE cashback_service_db;
    CREATE DATABASE mint_consumer_db;
    CREATE DATABASE blockchain_adapter_db;
EOSQL

for db in cashback_service_db mint_consumer_db blockchain_adapter_db; do
    echo "Applying schema to $db..."
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$db" -f /schema.sql
done

echo "Database initialization complete."
