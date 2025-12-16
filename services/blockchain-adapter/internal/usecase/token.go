package usecase

import "context"

type (
	TokenUsecase struct{}

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

func NewTokenUsecase() *TokenUsecase {
	return &TokenUsecase{}
}

func (TokenUsecase) MintToken(_ context.Context, _, _, _ string) (*MintResult, error) {
	// TODO: Implementar lógica de mint
	return &MintResult{Success: false}, nil
}

func (TokenUsecase) GetBalance(_ context.Context, _ string) (*BalanceResult, error) {
	// TODO: Implementar lógica de obtenção de saldo
	return &BalanceResult{}, nil
}

func (TokenUsecase) GetTransaction(_ context.Context, _ string) (*TransactionResult, error) {
	// TODO: Implementar lógica de obtenção de transação
	return &TransactionResult{}, nil
}
