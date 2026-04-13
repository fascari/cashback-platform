package ethereum

import (
	"context"

	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/usecase"
	"github.com/cashback-platform/services/blockchain-adapter/internal/chain"
)

type Client struct {
	useCase usecase.TokenUsecase
}

func New(uc usecase.TokenUsecase) Client {
	return Client{useCase: uc}
}

func (c Client) ChainID() chain.ID {
	return chain.Ethereum
}

func (c Client) MintToken(ctx context.Context, req chain.MintTokenRequest) (*chain.MintTokenResult, error) {
	result, err := c.useCase.MintToken(ctx, req.IdempotencyKey, req.WalletAddress, req.TokenAmount)
	if err != nil || result == nil {
		return nil, err
	}
	return new(chain.MintTokenResult{
		Success:         result.Success,
		TransactionHash: result.TransactionHash,
		BlockNumber:     result.BlockNumber,
		Status:          result.Status,
		ErrorCode:       result.ErrorCode,
		ErrorMessage:    result.ErrorMessage,
		Retryable:       result.Retryable,
	}), nil
}

func (c Client) FetchBalance(ctx context.Context, walletAddress string) (*chain.BalanceResult, error) {
	result, err := c.useCase.Balance(ctx, walletAddress)
	if err != nil || result == nil {
		return nil, err
	}
	return new(chain.BalanceResult{
		WalletAddress: result.WalletAddress,
		Balance:       result.Balance,
		BlockNumber:   result.BlockNumber,
	}), nil
}

func (c Client) FetchTransaction(ctx context.Context, txHash string) (*chain.TransactionResult, error) {
	result, err := c.useCase.Transaction(ctx, txHash)
	if err != nil || result == nil {
		return nil, err
	}
	return new(chain.TransactionResult{
		TransactionHash: result.TransactionHash,
		Status:          result.Status,
		BlockNumber:     result.BlockNumber,
		Confirmations:   result.Confirmations,
		GasUsed:         result.GasUsed,
		Success:         result.Success,
	}), nil
}
