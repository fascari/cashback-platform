package balance

import (
	"context"
	"fmt"
	"math/big"

	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

type (
	UserRepository interface {
		FindByID(ctx context.Context, id int64) (userdomain.User, error)
	}

	BlockchainClient interface {
		Balance(ctx context.Context, walletAddress string) (TokenBalance, error)
	}

	UseCase struct {
		userRepository   UserRepository
		blockchainClient BlockchainClient
	}
)

func New(userRepository UserRepository, blockchainClient BlockchainClient) UseCase {
	return UseCase{
		userRepository:   userRepository,
		blockchainClient: blockchainClient,
	}
}

// Execute returns the on-chain token balance for the given user.
func (u UseCase) Execute(ctx context.Context, userID int64) (Output, error) {
	user, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		return Output{}, fmt.Errorf("find user: %w", err)
	}

	bal, err := u.blockchainClient.Balance(ctx, user.WalletAddress)
	if err != nil {
		return Output{}, fmt.Errorf("fetch balance: %w", err)
	}

	return Output{
		UserID:        user.ID,
		WalletAddress: user.WalletAddress,
		Balance:       bal.Amount,
		BalanceTokens: toTokens(bal.Amount),
		BlockNumber:   bal.BlockNumber,
	}, nil
}

// toTokens converts a wei string to a human-readable token amount (divides by 10^18).
func toTokens(wei string) string {
	amount, ok := new(big.Int).SetString(wei, 10)
	if !ok || amount.Sign() == 0 {
		return "0"
	}
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	tokens := new(big.Float).Quo(
		new(big.Float).SetInt(amount),
		new(big.Float).SetInt(divisor),
	)
	return tokens.Text('f', 2)
}
