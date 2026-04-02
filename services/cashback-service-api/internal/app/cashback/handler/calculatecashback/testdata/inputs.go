package testdata

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler/calculatecashback"
)

func ValidPayload() calculatecashback.InputPayload {
	return calculatecashback.InputPayload{PurchaseID: "1"}
}

func MissingFieldsPayload() calculatecashback.InputPayload {
	return calculatecashback.InputPayload{}
}
