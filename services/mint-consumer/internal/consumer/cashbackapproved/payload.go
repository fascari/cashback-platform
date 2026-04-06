package cashbackapproved

import (
	"fmt"
	"strconv"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback"
	"github.com/google/uuid"
)

type cashbackApprovedPayload struct {
	EventID       string `json:"event_id"`
	CashbackID    string `json:"cashback_id"`
	UserID        string `json:"user_id"`
	WalletAddress string `json:"wallet_address"`
	PurchaseID    string `json:"purchase_id"`
	TokenAmount   string `json:"token_amount"`
	ChainID       string `json:"chain_id"`
}

func (p cashbackApprovedPayload) toDomain() (mintcashback.Input, error) {
	eventID, err := uuid.Parse(p.EventID)
	if err != nil {
		return mintcashback.Input{}, fmt.Errorf("event_id %q: %w", p.EventID, err)
	}

	cashbackID, err := strconv.ParseInt(p.CashbackID, 10, 64)
	if err != nil {
		return mintcashback.Input{}, fmt.Errorf("cashback_id %q: %w", p.CashbackID, err)
	}

	userID, err := strconv.ParseInt(p.UserID, 10, 64)
	if err != nil {
		return mintcashback.Input{}, fmt.Errorf("user_id %q: %w", p.UserID, err)
	}

	return mintcashback.Input{
		EventID:       eventID,
		CashbackID:    cashbackID,
		UserID:        userID,
		WalletAddress: p.WalletAddress,
		PurchaseID:    p.PurchaseID,
		TokenAmount:   p.TokenAmount,
		ChainID:       p.ChainID,
	}, nil
}
