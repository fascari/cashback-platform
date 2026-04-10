package testdata

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"

	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/domain"
)

const (
	WalletAddress = "0x71C7656EC7ab88b098defB751B7401B5f6d8976F"
	TokenAmount   = "1000000000000000000"
	TxHash        = "0xabc123"
	RecordID      = int64(1)
	Nonce         = int64(42)
	OnChainNonce  = int64(43)
)

func SubmittedTransaction(key uuid.UUID) *domain.BlockchainTransaction {
	return &domain.BlockchainTransaction{
		ID:              RecordID,
		IdempotencyKey:  key,
		Status:          domain.TransactionStatusSubmitted,
		TransactionHash: TxHash,
	}
}

func FailedTransaction(key uuid.UUID) *domain.BlockchainTransaction {
	return &domain.BlockchainTransaction{
		ID:             RecordID,
		IdempotencyKey: key,
		Status:         domain.TransactionStatusFailed,
		ErrorCode:      "send_failed",
		ErrorMessage:   "timeout",
	}
}

func RecentPendingTransaction(key uuid.UUID) *domain.BlockchainTransaction {
	return &domain.BlockchainTransaction{
		ID:             RecordID,
		IdempotencyKey: key,
		Status:         domain.TransactionStatusPending,
		CreatedAt:      time.Now(),
	}
}

func SignedTransaction() *types.Transaction {
	return types.NewTx(&types.LegacyTx{Nonce: uint64(Nonce)})
}
