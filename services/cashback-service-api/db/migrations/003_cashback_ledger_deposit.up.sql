ALTER TABLE cashback.cashback_ledger
    ALTER COLUMN purchase_id DROP NOT NULL,
    ALTER COLUMN purchase_id DROP DEFAULT;

ALTER TABLE cashback.cashback_ledger
    DROP CONSTRAINT IF EXISTS cashback_ledger_purchase_id_key;

CREATE UNIQUE INDEX uq_cashback_ledger_purchase_id
    ON cashback.cashback_ledger(purchase_id)
    WHERE purchase_id IS NOT NULL;

ALTER TABLE cashback.cashback_ledger
    ADD COLUMN deposit_receipt_id BIGINT
        REFERENCES cashback.deposit_receipts(id);

CREATE UNIQUE INDEX uq_cashback_ledger_deposit_receipt_id
    ON cashback.cashback_ledger(deposit_receipt_id)
    WHERE deposit_receipt_id IS NOT NULL;

CREATE INDEX idx_cashback_ledger_deposit_receipt_id
    ON cashback.cashback_ledger(deposit_receipt_id);
