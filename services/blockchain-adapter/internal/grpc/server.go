package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	tokenpb "github.com/cashback-platform/proto/token"
	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
	"github.com/cashback-platform/services/blockchain-adapter/internal/usecase"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type TokenServer struct {
	tokenpb.UnimplementedTokenServiceServer
	tokenUsecase usecase.TokenUsecase
}

func NewTokenServer(tokenUsecase usecase.TokenUsecase) *TokenServer {
	return &TokenServer{tokenUsecase: tokenUsecase}
}

func (s *TokenServer) MintToken(ctx context.Context, req *tokenpb.MintTokenRequest) (*tokenpb.MintTokenResponse, error) {
	result, err := s.tokenUsecase.MintToken(ctx, req.IdempotencyKey, req.WalletAddress, req.TokenAmount)
	if err != nil {
		return nil, err
	}

	response := &tokenpb.MintTokenResponse{
		Success:         result.Success,
		TransactionHash: result.TransactionHash,
		BlockNumber:     result.BlockNumber,
	}

	if !result.Success {
		response.Error = &tokenpb.MintError{
			Code:      result.ErrorCode,
			Message:   result.ErrorMessage,
			Retryable: result.Retryable,
		}
	}

	return response, nil
}

func (s *TokenServer) GetBalance(ctx context.Context, req *tokenpb.GetBalanceRequest) (*tokenpb.GetBalanceResponse, error) {
	result, err := s.tokenUsecase.Balance(ctx, req.WalletAddress)
	if err != nil {
		return nil, err
	}

	return &tokenpb.GetBalanceResponse{
		WalletAddress: result.WalletAddress,
		Balance:       result.Balance,
		BlockNumber:   result.BlockNumber,
	}, nil
}

func (s *TokenServer) GetTransaction(ctx context.Context, req *tokenpb.GetTransactionRequest) (*tokenpb.GetTransactionResponse, error) {
	result, err := s.tokenUsecase.Transaction(ctx, req.TransactionHash)
	if err != nil {
		return nil, err
	}

	return &tokenpb.GetTransactionResponse{
		TransactionHash: result.TransactionHash,
		BlockNumber:     result.BlockNumber,
		Confirmations:   result.Confirmations,
		GasUsed:         result.GasUsed,
		Success:         result.Success,
	}, nil
}

func StartServer(lc fx.Lifecycle, tokenServer *TokenServer, cfg *config.Config) {
	server := grpc.NewServer()

	tokenpb.RegisterTokenServiceServer(server, tokenServer)
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
