package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	inner *ethclient.Client
}

func New(rpcURL string) (*Client, error) {
	c, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial ethereum rpc %s: %w", rpcURL, err)
	}
	return &Client{inner: c}, nil
}

func (c *Client) Inner() *ethclient.Client {
	return c.inner
}

func (c *Client) Close() {
	c.inner.Close()
}

func (c *Client) ChainID(ctx context.Context) (int64, error) {
	id, err := c.inner.ChainID(ctx)
	if err != nil {
		return 0, fmt.Errorf("get chain id: %w", err)
	}
	return id.Int64(), nil
}

func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if err := c.inner.SendTransaction(ctx, tx); err != nil {
		return fmt.Errorf("send transaction: %w", err)
	}
	return nil
}

func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	price, err := c.inner.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("suggest gas price: %w", err)
	}
	return price, nil
}

func (c *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := c.inner.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("get transaction receipt: %w", err)
	}
	return receipt, nil
}

func (c *Client) PendingNonceAt(ctx context.Context, addr common.Address) (uint64, error) {
	n, err := c.inner.PendingNonceAt(ctx, addr)
	if err != nil {
		return 0, fmt.Errorf("get pending nonce: %w", err)
	}
	return n, nil
}
