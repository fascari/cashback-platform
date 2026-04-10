package grpc

import (
	"context"
	"fmt"
	"log/slog"

	tokenpb "github.com/cashback-platform/proto/token"
	"github.com/cashback-platform/services/cashback-service-api/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
		return nil, fmt.Errorf("failed to connect to blockchain adapter: %w", err)
	}

	slog.Info("Connected to blockchain adapter", "address", cfg.BlockchainAdapterAddress)
	return &BlockchainAdapterClient{
		conn:        conn,
		tokenClient: tokenpb.NewTokenServiceClient(conn),
	}, nil
}

func (c *BlockchainAdapterClient) MintToken(ctx context.Context, req *tokenpb.MintTokenRequest) (*tokenpb.MintTokenResponse, error) {
	return c.tokenClient.MintToken(ctx, req)
}

func (c *BlockchainAdapterClient) Balance(ctx context.Context, walletAddress string) (*tokenpb.GetBalanceResponse, error) {
	return c.tokenClient.GetBalance(ctx, &tokenpb.GetBalanceRequest{WalletAddress: walletAddress})
}

func (c *BlockchainAdapterClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
