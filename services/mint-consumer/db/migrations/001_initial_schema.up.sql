CREATE SCHEMA IF NOT EXISTS mint;

CREATE TYPE mint.mint_request_status AS ENUM ('pending', 'processing', 'completed', 'failed');

CREATE TABLE mint.processed_events (
    id BIGSERIAL PRIMARY KEY,
    event_id UUID UNIQUE NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE mint.mint_requests (
    id BIGSERIAL PRIMARY KEY,
    cashback_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    wallet_address VARCHAR(42) NOT NULL,
    token_amount VARCHAR(78) NOT NULL,
    idempotency_key UUID UNIQUE NOT NULL,
    status mint.mint_request_status NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 5,
    error_code VARCHAR(100),
    error_message TEXT,
    transaction_hash VARCHAR(66),
    block_number BIGINT,
    next_retry_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_mint_requests_cashback_id ON mint.mint_requests(cashback_id);
CREATE INDEX idx_mint_requests_status ON mint.mint_requests(status);
CREATE INDEX idx_mint_requests_next_retry_at ON mint.mint_requests(next_retry_at) WHERE status = 'failed';
