-- Database: blockchain_adapter_db
CREATE SCHEMA IF NOT EXISTS blockchain;
SET search_path TO blockchain;

-- Enums
CREATE TYPE transaction_status AS ENUM ('pending', 'submitted', 'confirmed', 'failed');
CREATE TYPE deposit_status AS ENUM ('pending', 'processed', 'failed');

-- Nonce tracking per wallet — serialized via Redis distributed lock
CREATE TABLE wallet_nonces (
    id BIGSERIAL PRIMARY KEY,
    wallet_address VARCHAR(42) UNIQUE NOT NULL,
    current_nonce BIGINT NOT NULL DEFAULT 0,
    fence_token BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_wallet_nonces_wallet_address ON wallet_nonces(wallet_address);

CREATE TABLE blockchain_transactions (
    id BIGSERIAL PRIMARY KEY,
    idempotency_key UUID UNIQUE NOT NULL, -- chave de negócio para deduplicação gRPC
    wallet_address VARCHAR(42) NOT NULL,
    token_amount VARCHAR(78) NOT NULL,
    chain_id VARCHAR(50) NOT NULL DEFAULT 'ethereum',
    transaction_hash VARCHAR(66),
    block_number BIGINT,
    gas_used BIGINT,
    gas_price VARCHAR(78),
    status transaction_status NOT NULL DEFAULT 'pending',
    error_code VARCHAR(100),
    error_message TEXT,
    nonce BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    confirmed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_blockchain_transactions_idempotency_key ON blockchain_transactions(idempotency_key);
CREATE INDEX idx_blockchain_transactions_transaction_hash ON blockchain_transactions(transaction_hash);
CREATE INDEX idx_blockchain_transactions_status ON blockchain_transactions(status);
CREATE INDEX idx_blockchain_transactions_chain_id ON blockchain_transactions(chain_id);

-- Deposit monitor: on-chain deposits detected by the blockchain monitor
CREATE TABLE detected_deposits (
    id BIGSERIAL PRIMARY KEY,
    chain_id VARCHAR(50) NOT NULL,
    transaction_hash VARCHAR(100) NOT NULL,
    wallet_address VARCHAR(100) NOT NULL,
    token_amount VARCHAR(78) NOT NULL,
    block_number BIGINT NOT NULL,
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status deposit_status NOT NULL DEFAULT 'pending',
    processed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(chain_id, transaction_hash)
);

CREATE INDEX idx_detected_deposits_wallet ON detected_deposits(wallet_address);
CREATE INDEX idx_detected_deposits_status ON detected_deposits(status);
CREATE INDEX idx_detected_deposits_chain ON detected_deposits(chain_id);
