package findpurchase_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/findpurchase"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/findpurchase/mocks"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/findpurchase/testdata"
)

func TestFindPurchase_ShouldReturnPurchaseWhenFound(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().FindByID(mock.Anything, testdata.PurchaseID).Return(testdata.FoundPurchase(), nil)

	uc := findpurchase.New(repo)
	result, err := uc.Execute(t.Context(), testdata.PurchaseID)

	require.NoError(t, err)
	require.Equal(t, testdata.FoundPurchase(), result)
}

func TestFindPurchase_ShouldReturnErrorWhenNotFound(t *testing.T) {
	repo := mocks.NewRepository(t)

	const id int64 = 99

	repo.EXPECT().FindByID(mock.Anything, id).
		Return(purchasedomain.Purchase{}, purchasedomain.ErrPurchaseNotFound)

	uc := findpurchase.New(repo)
	_, err := uc.Execute(t.Context(), id)

	require.ErrorIs(t, err, purchasedomain.ErrPurchaseNotFound)
}
