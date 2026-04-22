package processdeposit

//go:generate mockery --all

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/google/uuid"

	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	calculatecashback "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/domain"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

const (
	// DefaultCashbackPercent is the percentage of deposited tokens credited as cashback.
	DefaultCashbackPercent = 1.0
	defaultChainID         = "ethereum-sepolia"
)

type (
	UserRepository interface {
		FindByWalletAddress(ctx context.Context, walletAddress string) (userdomain.User, error)
	}

	// DepositRepository provides deposit receipt persistence and duplicate detection.
	DepositRepository interface {
		ExistsByTxHash(ctx context.Context, txHash string) (bool, error)
		Save(ctx context.Context, receipt domain.DepositReceipt) (domain.DepositReceipt, error)
	}

	// CashbackRepository persists cashback entries and publishes the approval event atomically.
	CashbackRepository interface {
		CreateWithEvent(ctx context.Context, cashback cashdomain.Cashback, eventFn func(cashdomain.Cashback) any) (cashdomain.Cashback, error)
	}

	UseCase struct {
		userRepository     UserRepository
		depositRepository  DepositRepository
		cashbackRepository CashbackRepository
	}
)

func New(
	userRepo UserRepository,
	depositRepo DepositRepository,
	cashbackRepo CashbackRepository,
) UseCase {
	return UseCase{
		userRepository:     userRepo,
		depositRepository:  depositRepo,
		cashbackRepository: cashbackRepo,
	}
}

func (u UseCase) Execute(ctx context.Context, input Input) error {
	user, err := u.userRepository.FindByWalletAddress(ctx, input.FromAddress)
	if err != nil {
		return fmt.Errorf("finding user by wallet address: %w", err)
	}

	exists, err := u.depositRepository.ExistsByTxHash(ctx, input.TransactionHash)
	if err != nil {
		return fmt.Errorf("checking existing deposit: %w", err)
	}
	if exists {
		return domain.ErrDepositAlreadyProcessed
	}

	tokens, err := weiToTokens(input.TokenAmount)
	if err != nil {
		return fmt.Errorf("parsing token amount %q: %w", input.TokenAmount, err)
	}

	receipt, err := domain.NewDepositReceipt(
		user.ID,
		input.TransactionHash,
		input.FromAddress,
		input.TokenAmount,
		input.ChainID,
		input.BlockNumber,
		input.DetectedAt,
	)
	if err != nil {
		return fmt.Errorf("building deposit receipt: %w", err)
	}

	saved, err := u.depositRepository.Save(ctx, receipt)
	if err != nil {
		return fmt.Errorf("saving deposit receipt: %w", err)
	}

	cashback, err := cashdomain.NewCashbackFromDeposit(user.ID, saved.ID, tokens, DefaultCashbackPercent)
	if err != nil {
		return fmt.Errorf("building cashback from deposit: %w", err)
	}

	cashback = cashback.Approve()

	_, err = u.cashbackRepository.CreateWithEvent(ctx, cashback, func(created cashdomain.Cashback) any {
		return calculatecashback.CashbackApprovedEvent{
			EventID:       uuid.New().String(),
			CashbackID:    strconv.FormatInt(created.ID, 10),
			UserID:        strconv.FormatInt(created.UserID, 10),
			WalletAddress: user.WalletAddress,
			PurchaseID:    "",
			TokenAmount:   fmt.Sprintf("%.0f", created.Amount*1e18),
			ChainID:       defaultChainID,
		}
	})
	if err != nil {
		return fmt.Errorf("creating cashback with event: %w", err)
	}

	return nil
}

func weiToTokens(weiStr string) (float64, error) {
	wei := new(big.Int)
	if _, ok := wei.SetString(weiStr, 10); !ok {
		return 0, fmt.Errorf("invalid wei string: %s", weiStr)
	}
	weiFloat := new(big.Float).SetInt(wei)
	divisor := new(big.Float).SetFloat64(1e18)
	result, _ := new(big.Float).Quo(weiFloat, divisor).Float64()
	return result, nil
}
