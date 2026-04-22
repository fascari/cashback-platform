//go:build e2e

package assert

import (
	"database/sql"

	e2esuite "github.com/cashback-platform/test/e2e/suite"
)

func IsDetected(db *sql.DB, txHash string) bool {
	return e2esuite.RowExists(db,
		"SELECT COUNT(*) FROM detected_deposits WHERE transaction_hash = $1",
		txHash,
	)
}
