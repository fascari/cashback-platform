package grpc

import (
	"context"
	"fmt"

	tokenpb "github.com/cashback-platform/proto/token"
	"github.com/cashback-platform/services/mint-consumer/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// BlockchainAdapterClient wraps the generated gRPC client for the blockchain adapter.
type BlockchainAdapterClient struct {
	conn        *grpc.ClientConn
	tokenClient tokenpb.TokenServiceClient
}

func NewBlockchainAdapterClient(cfg config.GRPC) (*BlockchainAdapterClient, error) {
	conn, err := grpc.NewClient(
		cfg.BlockchainAdapterAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to blockchain adapter: %w", err)
	}

	return &BlockchainAdapterClient{
		conn:        conn,
		tokenClient: tokenpb.NewTokenServiceClient(conn),
	}, nil
}

// MintToken calls the blockchain adapter's MintToken RPC.
func (c *BlockchainAdapterClient) MintToken(ctx context.Context, idempotencyKey, walletAddress, tokenAmount string) (*tokenpb.MintTokenResponse, error) {
	return c.tokenClient.MintToken(ctx, &tokenpb.MintTokenRequest{
		IdempotencyKey: idempotencyKey,
		WalletAddress:  walletAddress,
		TokenAmount:    tokenAmount,
	})
}

func (c *BlockchainAdapterClient) Connection() *grpc.ClientConn {
	return c.conn
}

func (c *BlockchainAdapterClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
