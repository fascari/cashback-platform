DROP INDEX IF EXISTS cashback.uq_cashback_ledger_deposit_receipt_id;
DROP INDEX IF EXISTS cashback.idx_cashback_ledger_deposit_receipt_id;
DROP INDEX IF EXISTS cashback.uq_cashback_ledger_purchase_id;

ALTER TABLE cashback.cashback_ledger
    DROP COLUMN IF EXISTS deposit_receipt_id;

ALTER TABLE cashback.cashback_ledger
    ADD CONSTRAINT cashback_ledger_purchase_id_key UNIQUE (purchase_id),
    ALTER COLUMN purchase_id SET NOT NULL;
