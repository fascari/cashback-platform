package depositdetected

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/usecase/processdeposit"
)

type depositPayload struct {
	ChainID         string    `json:"chain_id"`
	TransactionHash string    `json:"transaction_hash"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	TokenAmount     string    `json:"token_amount"`
	BlockNumber     int64     `json:"block_number"`
	DetectedAt      time.Time `json:"detected_at"`
}

func parsePayload(data []byte) (processdeposit.Input, error) {
	var p depositPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return processdeposit.Input{}, fmt.Errorf("unmarshaling deposit payload: %w", err)
	}
	if p.TransactionHash == "" {
		return processdeposit.Input{}, fmt.Errorf("missing transaction_hash in payload")
	}
	if p.FromAddress == "" {
		return processdeposit.Input{}, fmt.Errorf("missing from_address in payload")
	}
	return processdeposit.Input{
		ChainID:         p.ChainID,
		TransactionHash: p.TransactionHash,
		FromAddress:     p.FromAddress,
		ToAddress:       p.ToAddress,
		TokenAmount:     p.TokenAmount,
		BlockNumber:     p.BlockNumber,
		DetectedAt:      p.DetectedAt,
	}, nil
}
