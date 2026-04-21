//go:build integration

package assert

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func DepositSaved(t *testing.T, db *gorm.DB, chainID, txHash string) {
	t.Helper()
	var count int64
	err := db.Raw(
		"SELECT COUNT(*) FROM blockchain.detected_deposits WHERE chain_id = ? AND transaction_hash = ?",
		chainID, txHash,
	).Scan(&count).Error
	require.NoError(t, err)
	require.Equal(t, int64(1), count, "deposit %s/%s should exist in the database", chainID, txHash)
}
