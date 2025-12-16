package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
	"github.com/cashback-platform/services/blockchain-adapter/internal/usecase"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TokenServer implements the gRPC TokenService
type (
	TokenServer struct {
		tokenUsecase *usecase.TokenUsecase
	}

	// MintTokenRequest represents a request to mint tokens
	MintTokenRequest struct {
		IdempotencyKey string
		WalletAddress  string
		TokenAmount    string
	}

	// MintTokenResponse represents the response from a mint operation
	MintTokenResponse struct {
		Success         bool
		TransactionHash string
		BlockNumber     int64
		Status          string
		Error           *MintError
	}

	// MintError represents an error in mint operation
	MintError struct {
		Code      string
		Message   string
		Retryable bool
	}

	// GetBalanceRequest represents a request to get balance
	GetBalanceRequest struct {
		WalletAddress string
	}

	// GetBalanceResponse represents the response from a balance query
	GetBalanceResponse struct {
		WalletAddress string
		Balance       string
		BlockNumber   int64
	}

	// GetTransactionRequest represents a request to get transaction status
	GetTransactionRequest struct {
		TransactionHash string
	}

	// GetTransactionResponse represents the response from a transaction query
	GetTransactionResponse struct {
		TransactionHash string
		Status          string
		BlockNumber     int64
		Confirmations   int64
		GasUsed         int64
		Success         bool
	}
)

func NewTokenServer(tokenUsecase *usecase.TokenUsecase) *TokenServer {
	return &TokenServer{tokenUsecase: tokenUsecase}
}

// MintToken handles the MintToken gRPC call
func (s *TokenServer) MintToken(ctx context.Context, req *MintTokenRequest) (*MintTokenResponse, error) {
	result, err := s.tokenUsecase.MintToken(ctx, req.IdempotencyKey, req.WalletAddress, req.TokenAmount)
	if err != nil {
		return nil, err
	}

	response := &MintTokenResponse{
		Success:         result.Success,
		TransactionHash: result.TransactionHash,
		BlockNumber:     result.BlockNumber,
		Status:          result.Status,
	}

	if !result.Success {
		response.Error = &MintError{
			Code:      result.ErrorCode,
			Message:   result.ErrorMessage,
			Retryable: result.Retryable,
		}
	}

	return response, nil
}

// GetBalance handles the GetBalance gRPC call
func (s *TokenServer) GetBalance(ctx context.Context, req *GetBalanceRequest) (*GetBalanceResponse, error) {
	result, err := s.tokenUsecase.GetBalance(ctx, req.WalletAddress)
	if err != nil {
		return nil, err
	}

	return &GetBalanceResponse{
		WalletAddress: result.WalletAddress,
		Balance:       result.Balance,
		BlockNumber:   result.BlockNumber,
	}, nil
}

// GetTransaction handles the GetTransaction gRPC call
func (s *TokenServer) GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error) {
	result, err := s.tokenUsecase.GetTransaction(ctx, req.TransactionHash)
	if err != nil {
		return nil, err
	}

	return &GetTransactionResponse{
		TransactionHash: result.TransactionHash,
		Status:          result.Status,
		BlockNumber:     result.BlockNumber,
		Confirmations:   result.Confirmations,
		GasUsed:         result.GasUsed,
		Success:         result.Success,
	}, nil
}

func StartServer(lc fx.Lifecycle, _ *TokenServer, cfg *config.Config) {
	server := grpc.NewServer()

	// Register reflection for debugging
	reflection.Register(server)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}

			go func() {
				log.Printf("gRPC server starting on port %s", cfg.GRPC.Port)
				if err := server.Serve(listener); err != nil {
					log.Printf("gRPC server error: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Println("Shutting down gRPC server...")
			server.GracefulStop()
			return nil
		},
	})
}
