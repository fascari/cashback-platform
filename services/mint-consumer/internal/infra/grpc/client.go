package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/cashback-platform/services/mint-consumer/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	// MintResult represents the result of a mint operation
	MintResult struct {
		Success         bool
		TransactionHash string
		BlockNumber     int64
		ErrorCode       string
		ErrorMessage    string
		Retryable       bool
	}

	BlockchainAdapterClient struct {
		conn    *grpc.ClientConn
		address string
	}
)

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
		conn:    conn,
		address: cfg.GRPC.BlockchainAdapterAddress,
	}, nil
}

func (*BlockchainAdapterClient) MintToken(_ context.Context, idempotencyKey, walletAddress, tokenAmount string) (*MintResult, error) {
	// TODO: Use generated gRPC client from proto files
	// For now, return a mock successful response
	log.Printf("Minting token: idempotencyKey=%s, wallet=%s, amount=%s", idempotencyKey, walletAddress, tokenAmount)

	// Simulated successful mint
	return &MintResult{
		Success:         true,
		TransactionHash: fmt.Sprintf("0x%s", idempotencyKey[:32]),
		BlockNumber:     12345678,
	}, nil
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
