package usecase

import "context"

type MintUsecase struct{}

func NewMintUsecase() *MintUsecase {
	return &MintUsecase{}
}

func (MintUsecase) ProcessCashbackApproved(_ context.Context, _ []byte) error {
	// TODO: Implementar lógica de processamento
	return nil
}

func (MintUsecase) RetryFailedMints(_ context.Context) error {
	// TODO: Implementar lógica de retry
	return nil
}
