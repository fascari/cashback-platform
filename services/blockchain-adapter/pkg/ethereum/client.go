package ethereum

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

// Client wraps ethclient.Client.
type Client struct {
	inner *ethclient.Client
}

// New dials an Ethereum RPC endpoint and returns a connected Client.
func New(rpcURL string) (*Client, error) {
	c, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial ethereum rpc %s: %w", rpcURL, err)
	}
	return &Client{inner: c}, nil
}

// Inner returns the underlying ethclient.Client for direct use.
func (c *Client) Inner() *ethclient.Client {
	return c.inner
}

// Close terminates the client connection.
func (c *Client) Close() {
	c.inner.Close()
}

// ChainID retrieves the chain ID from the connected node.
func (c *Client) ChainID(ctx context.Context) (int64, error) {
	id, err := c.inner.ChainID(ctx)
	if err != nil {
		return 0, fmt.Errorf("get chain id: %w", err)
	}
	return id.Int64(), nil
}
