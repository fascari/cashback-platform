package grpc

import (
	"fmt"
	"log"

	"github.com/cashback-platform/services/cashback-service-api/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BlockchainAdapterClient struct {
	conn *grpc.ClientConn
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
	return &BlockchainAdapterClient{conn: conn}, nil
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
