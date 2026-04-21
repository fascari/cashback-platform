package ethereum_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cashback-platform/services/blockchain-adapter/internal/chain"
	ethereuminfra "github.com/cashback-platform/services/blockchain-adapter/internal/infra/ethereum"
	"github.com/cashback-platform/services/blockchain-adapter/internal/infra/ethereum/mocks"
)

func TestDepositMonitor_Watch_ShouldCallHandlerForEachDeposit(t *testing.T) {
	filter := mocks.NewTokenFilter(t)
	blocks := mocks.NewBlockReader(t)

	deposit := chain.Deposit{
		ChainID:         chain.Ethereum,
		TransactionHash: "0xabc",
		FromAddress:     "0xfrom",
		ToAddress:       "0xto",
		TokenAmount:     "1000000000000000000",
		BlockNumber:     100,
		DetectedAt:      time.Now().UTC(),
	}

	blocks.EXPECT().BlockNumber(mock.Anything).Return(uint64(100), nil).Once()
	filter.EXPECT().FilterTransfers(mock.Anything, uint64(1), uint64(100)).Return([]chain.Deposit{deposit}, nil).Once()

	// stop after the first tick
	blocks.EXPECT().BlockNumber(mock.Anything).Return(uint64(100), nil).Maybe()
	filter.EXPECT().FilterTransfers(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	var received []chain.Deposit
	handler := func(_ context.Context, d chain.Deposit) error {
		received = append(received, d)
		return nil
	}

	m := ethereuminfra.NewDepositMonitor(filter, blocks, 0, ethereuminfra.WithPollInterval(5*time.Millisecond))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = m.Watch(ctx, handler)

	require.Len(t, received, 1)
	require.Equal(t, deposit.TransactionHash, received[0].TransactionHash)
}

func TestDepositMonitor_Watch_ShouldContinueWhenBlockNumberFails(t *testing.T) {
	filter := mocks.NewTokenFilter(t)
	blocks := mocks.NewBlockReader(t)

	blocks.EXPECT().BlockNumber(mock.Anything).Return(uint64(0), errors.New("rpc error")).Maybe()

	handler := func(_ context.Context, _ chain.Deposit) error { return nil }

	m := ethereuminfra.NewDepositMonitor(filter, blocks, 0, ethereuminfra.WithPollInterval(5*time.Millisecond))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = m.Watch(ctx, handler)
	// no panic — monitor swallows the error and continues
}

func TestDepositMonitor_Watch_ShouldContinueWhenFilterFails(t *testing.T) {
	filter := mocks.NewTokenFilter(t)
	blocks := mocks.NewBlockReader(t)

	blocks.EXPECT().BlockNumber(mock.Anything).Return(uint64(100), nil).Maybe()
	filter.EXPECT().FilterTransfers(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("filter error")).Maybe()

	handler := func(_ context.Context, _ chain.Deposit) error { return nil }

	m := ethereuminfra.NewDepositMonitor(filter, blocks, 0, ethereuminfra.WithPollInterval(5*time.Millisecond))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = m.Watch(ctx, handler)
}

func TestDepositMonitor_Watch_ShouldResumeFromStartBlock(t *testing.T) {
	filter := mocks.NewTokenFilter(t)
	blocks := mocks.NewBlockReader(t)

	blocks.EXPECT().BlockNumber(mock.Anything).Return(uint64(200), nil).Maybe()
	filter.EXPECT().FilterTransfers(mock.Anything, uint64(150), uint64(200)).Return(nil, nil).Maybe()

	m := ethereuminfra.NewDepositMonitor(filter, blocks, 150, ethereuminfra.WithPollInterval(5*time.Millisecond))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = m.Watch(ctx, func(_ context.Context, _ chain.Deposit) error { return nil })
}

func TestDepositMonitor_Stop_ShouldBeIdempotent(t *testing.T) {
	filter := mocks.NewTokenFilter(t)
	blocks := mocks.NewBlockReader(t)

	m := ethereuminfra.NewDepositMonitor(filter, blocks, 0)

	// multiple Stop calls must not panic (sync.Once guards the close)
	m.Stop()
	m.Stop()
}
