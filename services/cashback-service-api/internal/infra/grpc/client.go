package grpc

import (
	"context"
	"fmt"
	"log"

	tokenpb "github.com/cashback-platform/proto/token"
	"github.com/cashback-platform/services/cashback-service-api/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// BlockchainAdapterClient wraps the generated gRPC client for the blockchain adapter.
type BlockchainAdapterClient struct {
	conn        *grpc.ClientConn
	tokenClient tokenpb.TokenServiceClient
}

func NewBlockchainAdapterClient(cfg *config.Config) (*BlockchainAdapterClient, error) {
	conn, err := grpc.Dial(
		cfg.GRPC.BlockchainAdapterAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain adapter: %w", err)
	}

	log.Printf("Connected to blockchain adapter at %s", cfg.GRPC.BlockchainAdapterAddress)
	return &BlockchainAdapterClient{
		conn:        conn,
		tokenClient: tokenpb.NewTokenServiceClient(conn),
	}, nil
}

// MintToken calls the blockchain adapter's MintToken RPC.
func (c *BlockchainAdapterClient) MintToken(ctx context.Context, req *tokenpb.MintTokenRequest) (*tokenpb.MintTokenResponse, error) {
	return c.tokenClient.MintToken(ctx, req)
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
