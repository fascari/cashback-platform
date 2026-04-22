CREATE TABLE cashback.deposit_receipts (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES cashback.users(id),
    tx_hash      VARCHAR(66) NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    amount       VARCHAR(78) NOT NULL,
    chain_id     VARCHAR(32) NOT NULL,
    block_number BIGINT NOT NULL,
    detected_at  TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_deposit_receipts_tx_hash UNIQUE (tx_hash)
);

CREATE INDEX idx_deposit_receipts_user_id ON cashback.deposit_receipts(user_id);
CREATE INDEX idx_deposit_receipts_tx_hash ON cashback.deposit_receipts(tx_hash);
