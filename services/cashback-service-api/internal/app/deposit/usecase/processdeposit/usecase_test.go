package processdeposit_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/usecase/processdeposit"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/usecase/processdeposit/mocks"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/usecase/processdeposit/testdata"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

func TestProcessDeposit_ShouldCreditCashbackWhenDepositIsNew(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	depositRepo := mocks.NewDepositRepository(t)
	cashbackRepo := mocks.NewCashbackRepository(t)

	userRepo.EXPECT().FindByWalletAddress(mock.Anything, testdata.FromAddress).
		Return(testdata.FoundUser(), nil)
	depositRepo.EXPECT().ExistsByTxHash(mock.Anything, testdata.TxHash).
		Return(false, nil)
	depositRepo.EXPECT().Save(mock.Anything, mock.Anything).
		Return(testdata.SavedReceipt(), nil)
	cashbackRepo.EXPECT().CreateWithEvent(mock.Anything, mock.Anything, mock.Anything).
		Return(testdata.ApprovedCashback(), nil)

	uc := processdeposit.New(userRepo, depositRepo, cashbackRepo)
	err := uc.Execute(t.Context(), testdata.ValidInput())

	require.NoError(t, err)
}

func TestProcessDeposit_ShouldReturnErrorWhenUserNotFound(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	depositRepo := mocks.NewDepositRepository(t)
	cashbackRepo := mocks.NewCashbackRepository(t)

	userRepo.EXPECT().FindByWalletAddress(mock.Anything, testdata.FromAddress).
		Return(userdomain.User{}, userdomain.ErrUserNotFound)

	uc := processdeposit.New(userRepo, depositRepo, cashbackRepo)
	err := uc.Execute(t.Context(), testdata.ValidInput())

	require.Error(t, err)
	require.True(t, errors.Is(err, userdomain.ErrUserNotFound))
}

func TestProcessDeposit_ShouldReturnErrAlreadyProcessedWhenDuplicate(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	depositRepo := mocks.NewDepositRepository(t)
	cashbackRepo := mocks.NewCashbackRepository(t)

	userRepo.EXPECT().FindByWalletAddress(mock.Anything, testdata.FromAddress).
		Return(testdata.FoundUser(), nil)
	depositRepo.EXPECT().ExistsByTxHash(mock.Anything, testdata.TxHash).
		Return(true, nil)

	uc := processdeposit.New(userRepo, depositRepo, cashbackRepo)
	err := uc.Execute(t.Context(), testdata.ValidInput())

	require.ErrorIs(t, err, domain.ErrDepositAlreadyProcessed)
}

func TestProcessDeposit_ShouldReturnErrorWhenTokenAmountInvalid(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	depositRepo := mocks.NewDepositRepository(t)
	cashbackRepo := mocks.NewCashbackRepository(t)

	input := testdata.ValidInput()
	input.TokenAmount = "not-a-number"

	userRepo.EXPECT().FindByWalletAddress(mock.Anything, testdata.FromAddress).
		Return(testdata.FoundUser(), nil)
	depositRepo.EXPECT().ExistsByTxHash(mock.Anything, testdata.TxHash).
		Return(false, nil)

	uc := processdeposit.New(userRepo, depositRepo, cashbackRepo)
	err := uc.Execute(t.Context(), input)

	require.Error(t, err)
}
