package cashbackapproved_test

import (
	"context"
	"testing"

	natsgo "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/handler/cashbackapproved"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback"
	mintcashmocks "github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback/mocks"
)

func txManagerPassthrough(t *testing.T) *mintcashmocks.TransactionManager {
	t.Helper()
	tm := mintcashmocks.NewTransactionManager(t)
	tm.EXPECT().WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	return tm
}

func TestHandle_ShouldCallUseCaseWithParsedPayload(t *testing.T) {
	msg := &natsgo.Msg{
		Data: []byte(`{"event_id":"11111111-1111-1111-1111-111111111111","cashback_id":"1","user_id":"10","wallet_address":"0xABC","purchase_id":"P1","token_amount":"100","chain_id":"sepolia"}`),
	}

	repo := mintcashmocks.NewRepository(t)
	repo.EXPECT().ExistsProcessedEvent(mock.Anything, mock.Anything).Return(false, nil)
	repo.EXPECT().CreateMintRequest(mock.Anything, mock.Anything).Return(domain.MintRequest{ID: 42}, nil)
	repo.EXPECT().CreateProcessedEvent(mock.Anything, mock.Anything, mock.Anything).Return(nil)

	blockchain := mintcashmocks.NewBlockchainClient(t)
	blockchain.EXPECT().MintToken(mock.Anything, mock.Anything).Return(domain.MintResult{TransactionHash: "0xhash", BlockNumber: 1}, nil)

	repo.EXPECT().MarkCompleted(mock.Anything, int64(42), "0xhash", int64(1)).Return(nil)

	uc := mintcashback.NewUseCase(repo, blockchain, txManagerPassthrough(t))
	h := cashbackapproved.New(uc)

	err := h.Handle(context.Background(), msg)
	require.NoError(t, err)
}

func TestHandle_ShouldReturnErrorWhenPayloadIsInvalidJSON(t *testing.T) {
	msg := &natsgo.Msg{
		Data: []byte("{not valid json"),
	}

	repo := mintcashmocks.NewRepository(t)
	blockchain := mintcashmocks.NewBlockchainClient(t)
	tm := mintcashmocks.NewTransactionManager(t)

	uc := mintcashback.NewUseCase(repo, blockchain, tm)
	h := cashbackapproved.New(uc)

	err := h.Handle(context.Background(), msg)
	require.Error(t, err)
}

func TestHandle_ShouldReturnErrorWhenEventIDIsNotValidUUID(t *testing.T) {
	msg := &natsgo.Msg{
		Data: []byte(`{"event_id":"not-a-uuid","cashback_id":"1","user_id":"10","wallet_address":"0xABC","purchase_id":"P1","token_amount":"100","chain_id":"sepolia"}`),
	}

	repo := mintcashmocks.NewRepository(t)
	blockchain := mintcashmocks.NewBlockchainClient(t)
	tm := mintcashmocks.NewTransactionManager(t)

	uc := mintcashback.NewUseCase(repo, blockchain, tm)
	h := cashbackapproved.New(uc)

	err := h.Handle(context.Background(), msg)
	require.Error(t, err)
}

func TestHandle_ShouldReturnErrorWhenCashbackIDIsNotValidInt(t *testing.T) {
	msg := &natsgo.Msg{
		Data: []byte(`{"event_id":"11111111-1111-1111-1111-111111111111","cashback_id":"abc","user_id":"10","wallet_address":"0xABC","purchase_id":"P1","token_amount":"100","chain_id":"sepolia"}`),
	}

	repo := mintcashmocks.NewRepository(t)
	blockchain := mintcashmocks.NewBlockchainClient(t)
	tm := mintcashmocks.NewTransactionManager(t)

	uc := mintcashback.NewUseCase(repo, blockchain, tm)
	h := cashbackapproved.New(uc)

	err := h.Handle(context.Background(), msg)
	require.Error(t, err)
}
