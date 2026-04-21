//go:build e2e

package assert

import "database/sql"

func IsDetected(db *sql.DB, txHash string) bool {
	var count int
	_ = db.QueryRow(
		"SELECT COUNT(*) FROM detected_deposits WHERE transaction_hash = $1",
		txHash,
	).Scan(&count)
	return count > 0
}
