package events

const (
	CashbackApproved = "cashback.approved"
	DepositDetected  = "deposit.detected"
)

// Stream name constants — these must match the streams created by cmd/nats-setup.
const (
	StreamPurchaseEvents = "PURCHASE_EVENTS"
	StreamCashbackEvents = "CASHBACK_EVENTS"
	StreamTokenEvents    = "TOKEN_EVENTS"
	StreamDepositEvents  = "DEPOSIT_EVENTS"
)

const (
	SubjectPurchaseAll = "purchase.>"
	SubjectCashbackAll = "cashback.>"
	SubjectTokenAll    = "token.>"
	SubjectDepositAll  = "deposit.>"
)
