package handler

import (
	"context"

	tokenpb "github.com/cashback-platform/proto/token"
	"google.golang.org/grpc"

	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/usecase"
)

type Handler struct {
	tokenpb.UnimplementedTokenServiceServer
	useCase usecase.TokenUsecase
}

func NewHandler(useCase usecase.TokenUsecase) Handler {
	return Handler{useCase: useCase}
}

func RegisterServer(s *grpc.Server, h Handler) {
	tokenpb.RegisterTokenServiceServer(s, h)
}

func (h Handler) MintToken(ctx context.Context, req *tokenpb.MintTokenRequest) (*tokenpb.MintTokenResponse, error) {
	result, err := h.useCase.MintToken(ctx, req.IdempotencyKey, req.WalletAddress, req.TokenAmount)
	if err != nil || result == nil {
		return nil, err
	}

	resp := new(tokenpb.MintTokenResponse{
		Success:         result.Success,
		TransactionHash: result.TransactionHash,
		BlockNumber:     result.BlockNumber,
	})

	if !result.Success {
		resp.Error = new(tokenpb.MintError{
			Code:      result.ErrorCode,
			Message:   result.ErrorMessage,
			Retryable: result.Retryable,
		})
	}

	return resp, nil
}

func (h Handler) GetBalance(ctx context.Context, req *tokenpb.GetBalanceRequest) (*tokenpb.GetBalanceResponse, error) {
	result, err := h.useCase.Balance(ctx, req.WalletAddress)
	if err != nil || result == nil {
		return nil, err
	}

	return new(tokenpb.GetBalanceResponse{
		WalletAddress: result.WalletAddress,
		Balance:       result.Balance,
		BlockNumber:   result.BlockNumber,
	}), nil
}

func (h Handler) GetTransaction(ctx context.Context, req *tokenpb.GetTransactionRequest) (*tokenpb.GetTransactionResponse, error) {
	result, err := h.useCase.Transaction(ctx, req.TransactionHash)
	if err != nil || result == nil {
		return nil, err
	}

	return new(tokenpb.GetTransactionResponse{
		TransactionHash: result.TransactionHash,
		BlockNumber:     result.BlockNumber,
		Confirmations:   result.Confirmations,
		GasUsed:         result.GasUsed,
		Success:         result.Success,
	}), nil
}
