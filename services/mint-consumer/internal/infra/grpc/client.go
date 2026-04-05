package grpc

import (
	"context"
	"fmt"

	tokenpb "github.com/cashback-platform/proto/token"
	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/mint-consumer/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	// MintTokenRequest carries the inputs for a mint operation.
	MintTokenRequest struct {
		IdempotencyKey string
		WalletAddress  string
		TokenAmount    string
		ChainID        string
	}

	// MintResult carries the outcome of a mint operation.
	MintResult struct {
		TransactionHash string
		BlockNumber     int64
		ErrorCode       string
		ErrorMessage    string
		Retryable       bool
	}

	// Client wraps the blockchain adapter gRPC client.
	Client struct {
		conn        *grpc.ClientConn
		tokenClient tokenpb.TokenServiceClient
	}
)

// New connects to the blockchain adapter gRPC server.
func New(cfg config.GRPC) (*Client, error) {
	conn, err := grpc.NewClient(
		cfg.BlockchainAdapterAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to blockchain adapter: %w", err)
	}

	logger.Info("connected to blockchain adapter", "address", cfg.BlockchainAdapterAddress)
	return &Client{
		conn:        conn,
		tokenClient: tokenpb.NewTokenServiceClient(conn),
	}, nil
}

// MintToken calls the blockchain adapter and returns a domain-level result.
func (c *Client) MintToken(ctx context.Context, req MintTokenRequest) (MintResult, error) {
	resp, err := c.tokenClient.MintToken(ctx, &tokenpb.MintTokenRequest{
		IdempotencyKey: req.IdempotencyKey,
		WalletAddress:  req.WalletAddress,
		TokenAmount:    req.TokenAmount,
	})
	if err != nil {
		return MintResult{}, fmt.Errorf("mint token grpc: %w", err)
	}

	result := MintResult{
		TransactionHash: resp.GetTransactionHash(),
		BlockNumber:     resp.GetBlockNumber(),
	}

	if e := resp.GetError(); e != nil {
		result.ErrorCode = e.GetCode()
		result.ErrorMessage = e.GetMessage()
		result.Retryable = e.GetRetryable()
	}

	if resp.GetStatus() == tokenpb.MintStatus_MINT_STATUS_FAILED && result.ErrorCode == "" {
		result.ErrorCode = "mint_failed"
		result.ErrorMessage = "mint operation failed without details"
	}

	return result, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
