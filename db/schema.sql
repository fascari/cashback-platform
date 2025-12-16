-- ============================================================================
CREATE INDEX idx_wallet_nonces_wallet_address ON wallet_nonces(wallet_address);

);
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    current_nonce BIGINT NOT NULL DEFAULT 0,
    wallet_address VARCHAR(42) UNIQUE NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE wallet_nonces (
-- Wallet nonces: nonce tracking for transactions

CREATE INDEX idx_blockchain_transactions_status ON blockchain_transactions(status);
CREATE INDEX idx_blockchain_transactions_transaction_hash ON blockchain_transactions(transaction_hash);
CREATE INDEX idx_blockchain_transactions_idempotency_key ON blockchain_transactions(idempotency_key);

);
    confirmed_at TIMESTAMP WITH TIME ZONE
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    nonce BIGINT,
    error_message TEXT,
    error_code VARCHAR(100),
    -- pending, submitted, confirmed, failed
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    gas_price VARCHAR(78),
    gas_used BIGINT,
    block_number BIGINT,
    transaction_hash VARCHAR(66),
    token_amount VARCHAR(78) NOT NULL,
    wallet_address VARCHAR(42) NOT NULL,
    idempotency_key UUID UNIQUE NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE blockchain_transactions (
-- Blockchain transactions: status of on-chain mint operations

-- ============================================================================
-- Database: blockchain_adapter_db
-- BLOCKCHAIN ADAPTER DATABASE
-- ============================================================================

CREATE INDEX idx_mint_requests_next_retry_at ON mint_requests(next_retry_at) WHERE status = 'failed';
CREATE INDEX idx_mint_requests_idempotency_key ON mint_requests(idempotency_key);
CREATE INDEX idx_mint_requests_status ON mint_requests(status);
CREATE INDEX idx_mint_requests_cashback_id ON mint_requests(cashback_id);

);
    completed_at TIMESTAMP WITH TIME ZONE
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    next_retry_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    error_code VARCHAR(100),
    block_number BIGINT,
    transaction_hash VARCHAR(66),
    max_retries INT NOT NULL DEFAULT 5,
    retry_count INT NOT NULL DEFAULT 0,
    -- pending, processing, completed, failed
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    idempotency_key UUID UNIQUE NOT NULL,
    token_amount VARCHAR(78) NOT NULL,
    wallet_address VARCHAR(42) NOT NULL,
    user_id UUID NOT NULL,
    cashback_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE mint_requests (
-- Mint requests: state of mint operations

CREATE INDEX idx_processed_events_event_id ON processed_events(event_id);

);
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    event_type VARCHAR(100) NOT NULL,
    event_id UUID UNIQUE NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE processed_events (
-- Processed events: tracking for idempotency

-- ============================================================================
-- Database: mint_consumer_db
-- MINT CONSUMER DATABASE
-- ============================================================================

CREATE INDEX idx_outbox_events_aggregate ON outbox_events(aggregate_type, aggregate_id);
CREATE INDEX idx_outbox_events_created_at ON outbox_events(created_at);
CREATE INDEX idx_outbox_events_status ON outbox_events(status);

);
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    error_message TEXT,
    max_retries INT NOT NULL DEFAULT 5,
    retry_count INT NOT NULL DEFAULT 0,
    -- pending, published, failed
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    payload JSONB NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE outbox_events (
-- Outbox events: events pending publication (Outbox Pattern)

CREATE INDEX idx_cashback_ledger_status ON cashback_ledger(status);
CREATE INDEX idx_cashback_ledger_purchase_id ON cashback_ledger(purchase_id);
CREATE INDEX idx_cashback_ledger_user_id ON cashback_ledger(user_id);

);
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    calculation_basis JSONB,
    -- pending, approved, minting, minted, failed
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    token_amount VARCHAR(78) NOT NULL, -- Wei representation (uint256)
    amount DECIMAL(18, 8) NOT NULL,
    purchase_id UUID NOT NULL REFERENCES purchases(id),
    user_id UUID NOT NULL REFERENCES users(id),
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE cashback_ledger (
-- Cashback ledger: off-chain representation of generated cashback

CREATE INDEX idx_purchases_created_at ON purchases(created_at);
CREATE INDEX idx_purchases_status ON purchases(status);
CREATE INDEX idx_purchases_user_id ON purchases(user_id);

);
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    description TEXT,
    merchant_name VARCHAR(255),
    merchant_id VARCHAR(255),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    amount DECIMAL(18, 2) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE purchases (
-- Purchases table: stores purchase records

CREATE INDEX idx_users_external_id ON users(external_id);
CREATE INDEX idx_users_wallet_address ON users(wallet_address);

);
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    wallet_address VARCHAR(42) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    external_id VARCHAR(255) UNIQUE NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
CREATE TABLE users (
-- Users table: stores user accounts and wallet addresses

-- ============================================================================
-- Database: cashback_service_db
-- CASHBACK SERVICE API DATABASE
-- ============================================================================

-- ============================================================================
-- conceptual and separated per service. There is NO shared database.
-- IMPORTANT: Each service owns its own database. These schemas are
-- ============================================================================
-- Web3 Cashback Platform - Conceptual Database Schema

