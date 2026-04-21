package events

// Subject constants used when publishing or subscribing to individual event types.
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

// Subject wildcard constants — used when declaring stream subject filters.
const (
	SubjectPurchaseAll = "purchase.>"
	SubjectCashbackAll = "cashback.>"
	SubjectTokenAll    = "token.>"
	SubjectDepositAll  = "deposit.>"
)
