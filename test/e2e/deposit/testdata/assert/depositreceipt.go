//go:build e2e

package assert

import (
	"database/sql"

	e2esuite "github.com/cashback-platform/test/e2e/suite"
)

func DepositReceiptExists(db *sql.DB, txHash string) bool {
	return e2esuite.RowExists(db, `SELECT COUNT(*) FROM cashback.deposit_receipts WHERE tx_hash = $1`, txHash)
}

func CashbackFromDepositExists(db *sql.DB, walletAddress string) bool {
	return e2esuite.RowExists(db, `
		SELECT COUNT(*)
		FROM cashback.cashback_ledger cl
		JOIN cashback.deposit_receipts dr ON cl.deposit_receipt_id = dr.id
		JOIN cashback.users u ON dr.user_id = u.id
		WHERE u.wallet_address = $1
		  AND cl.deposit_receipt_id IS NOT NULL
	`, walletAddress)
}
