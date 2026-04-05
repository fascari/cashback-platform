package usecase

//go:generate mockery --all

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"

	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
	"github.com/cashback-platform/services/blockchain-adapter/internal/contracts"
	"github.com/cashback-platform/services/blockchain-adapter/internal/domain"
	"github.com/cashback-platform/kit/ethereum"
)

const errCodeSendFailed = "send_failed"

var (
	ErrTransactionFailed = errors.New("transaction permanently failed")
	ErrLockUnavailable   = errors.New("nonce lock unavailable — retry later")
)

type (
	NonceRepository interface {
		Increment(ctx context.Context, walletAddress string) (int64, error)
		SyncFromChain(ctx context.Context, walletAddress string, nonce int64) error
	}

	TransactionRepository interface {
		Create(ctx context.Context, tx domain.BlockchainTransaction) (domain.BlockchainTransaction, error)
		FindByIdempotencyKey(ctx context.Context, key uuid.UUID) (*domain.BlockchainTransaction, error)
		UpdateStatus(ctx context.Context, id int64, status domain.TransactionStatus, txHash string, blockNumber int64) error
		MarkFailed(ctx context.Context, id int64, errCode, errMsg string) error
	}

	EthereumClient interface {
		SendTransaction(ctx context.Context, tx *types.Transaction) error
		SuggestGasPrice(ctx context.Context) (*big.Int, error)
		TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
		PendingNonceAt(ctx context.Context, addr common.Address) (uint64, error)
	}

	TokenUsecase struct {
		nonceRepo       NonceRepository
		transactionRepo TransactionRepository
		ethClient       EthereumClient
		wallet          *ethereum.Wallet
		token           *contracts.CashbackToken
		cfg             *config.Config
	}

	MintResult struct {
		Success         bool
		TransactionHash string
		BlockNumber     int64
		Status          string
		ErrorCode       string
		ErrorMessage    string
		Retryable       bool
	}

	BalanceResult struct {
		WalletAddress string
		Balance       string
		BlockNumber   int64
	}

	TransactionResult struct {
		TransactionHash string
		Status          string
		BlockNumber     int64
		Confirmations   int64
		GasUsed         int64
		Success         bool
	}
)

func NewToken(
	nonceRepo NonceRepository,
	transactionRepo TransactionRepository,
	ethClient EthereumClient,
	wallet *ethereum.Wallet,
	token *contracts.CashbackToken,
	cfg *config.Config,
) TokenUsecase {
	return TokenUsecase{
		nonceRepo:       nonceRepo,
		transactionRepo: transactionRepo,
		ethClient:       ethClient,
		wallet:          wallet,
		token:           token,
		cfg:             cfg,
	}
}

func (u TokenUsecase) MintToken(ctx context.Context, idempotencyKeyStr, walletAddress, tokenAmount string) (*MintResult, error) {
	key, err := uuid.Parse(idempotencyKeyStr)
	if err != nil {
		return nil, fmt.Errorf("parse idempotency key: %w", err)
	}

	existing, early, err := u.checkIdempotency(ctx, key)
	if early != nil || err != nil {
		return early, err
	}

	nonce, err := u.nonceRepo.Increment(ctx, walletAddress)
	if err != nil {
		return &MintResult{Success: false, Retryable: true}, fmt.Errorf("%w: %w", ErrLockUnavailable, err)
	}

	signedTx, err := u.buildSignedTransaction(ctx, walletAddress, tokenAmount, nonce)
	if err != nil {
		return nil, err
	}
	txHash := signedTx.Hash().Hex()

	recordID, err := u.persistPending(ctx, key, existing, walletAddress, tokenAmount, txHash, nonce)
	if err != nil {
		_ = u.nonceRepo.SyncFromChain(ctx, walletAddress, nonce)
		return nil, err
	}

	if err := u.sendWithRetry(ctx, signedTx, recordID, walletAddress, nonce); err != nil {
		return &MintResult{
			Success:      false,
			Retryable:    true,
			ErrorCode:    errCodeSendFailed,
			ErrorMessage: err.Error(),
		}, nil
	}

	_ = u.transactionRepo.UpdateStatus(ctx, recordID, domain.TransactionStatusSubmitted, txHash, 0)

	return &MintResult{
		Success:         true,
		TransactionHash: txHash,
		Status:          string(domain.TransactionStatusSubmitted),
	}, nil
}

func (u TokenUsecase) buildSignedTransaction(ctx context.Context, walletAddress, tokenAmount string, nonce int64) (*types.Transaction, error) {
	amount, ok := new(big.Int).SetString(tokenAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid token amount: %s", tokenAmount)
	}

	gasPrice, err := u.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("get gas price: %w", err)
	}

	chainID := big.NewInt(u.cfg.Ethereum.ChainID)
	auth, err := bind.NewKeyedTransactorWithChainID(u.wallet.PrivateKey(), chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}
	auth.Nonce = big.NewInt(nonce)
	auth.GasPrice = gasPrice
	auth.NoSend = true

	tx, err := u.token.Mint(auth, common.HexToAddress(walletAddress), amount)
	if err != nil {
		return nil, fmt.Errorf("build mint transaction: %w", err)
	}
	return tx, nil
}

func (u TokenUsecase) persistPending(ctx context.Context, key uuid.UUID, existing *domain.BlockchainTransaction, walletAddress, tokenAmount, txHash string, nonce int64) (int64, error) {
	if existing != nil {
		return existing.ID, nil
	}
	record := domain.BlockchainTransaction{
		IdempotencyKey:  key,
		WalletAddress:   walletAddress,
		TokenAmount:     tokenAmount,
		TransactionHash: txHash,
		Status:          domain.TransactionStatusPending,
		Nonce:           nonce,
	}
	created, err := u.transactionRepo.Create(ctx, record)
	if err != nil {
		return 0, fmt.Errorf("persist transaction: %w", err)
	}
	return created.ID, nil
}

func (u TokenUsecase) checkIdempotency(ctx context.Context, key uuid.UUID) (*domain.BlockchainTransaction, *MintResult, error) {
	tx, err := u.transactionRepo.FindByIdempotencyKey(ctx, key)
	if err != nil {
		return nil, nil, fmt.Errorf("check idempotency: %w", err)
	}
	if tx == nil {
		return nil, nil, nil
	}
	if tx.IsFinalized() {
		return tx, &MintResult{Success: true, TransactionHash: tx.TransactionHash, Status: string(tx.Status)}, nil
	}
	if tx.IsFailed() {
		return tx, &MintResult{Success: false, ErrorCode: tx.ErrorCode, ErrorMessage: tx.ErrorMessage}, ErrTransactionFailed
	}
	if tx.Status == domain.TransactionStatusPending && time.Since(tx.CreatedAt) < 30*time.Second {
		return tx, &MintResult{Success: false, Retryable: true}, nil
	}
	return tx, nil, nil
}

func (u TokenUsecase) sendWithRetry(ctx context.Context, tx *types.Transaction, recordID int64, walletAddress string, nonce int64) error {
	const maxRetries = 3
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		lastErr = u.ethClient.SendTransaction(ctx, tx)
		if lastErr == nil {
			return nil
		}
		if i < maxRetries-1 {
			time.Sleep(time.Duration(1<<uint(i)) * 100 * time.Millisecond)
		}
	}

	syncNonce := nonce
	if onChain, err := u.ethClient.PendingNonceAt(ctx, u.wallet.Address()); err == nil {
		syncNonce = int64(onChain)
	}
	_ = u.nonceRepo.SyncFromChain(ctx, walletAddress, syncNonce)
	_ = u.transactionRepo.MarkFailed(ctx, recordID, errCodeSendFailed, lastErr.Error())

	return lastErr
}

func (u TokenUsecase) Balance(ctx context.Context, walletAddress string) (*BalanceResult, error) {
	addr := common.HexToAddress(walletAddress)
	balance, err := u.token.BalanceOf(&bind.CallOpts{Context: ctx}, addr)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}
	return &BalanceResult{
		WalletAddress: walletAddress,
		Balance:       balance.String(),
	}, nil
}

func (u TokenUsecase) Transaction(ctx context.Context, txHash string) (*TransactionResult, error) {
	receipt, err := u.ethClient.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, fmt.Errorf("get receipt: %w", err)
	}
	return &TransactionResult{
		TransactionHash: txHash,
		BlockNumber:     receipt.BlockNumber.Int64(),
		GasUsed:         int64(receipt.GasUsed),
		Success:         receipt.Status == 1,
		Status:          fmt.Sprintf("%d", receipt.Status),
	}, nil
}
