package processdeposit

import "time"

type Input struct {
	ChainID         string
	TransactionHash string
	FromAddress     string
	ToAddress       string
	TokenAmount     string // wei string, e.g. "1000000000000000000"
	BlockNumber     int64
	DetectedAt      time.Time
}
