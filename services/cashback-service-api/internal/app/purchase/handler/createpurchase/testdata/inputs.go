package testdata

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/handler/createpurchase"
)

func ValidPayload() createpurchase.InputPayload {
	return createpurchase.InputPayload{UserID: "42", Amount: 100.0, Merchant: "shopA"}
}

func MissingFieldsPayload() createpurchase.InputPayload {
	return createpurchase.InputPayload{}
}

func ZeroAmountPayload() createpurchase.InputPayload {
	return createpurchase.InputPayload{UserID: "42", Amount: 0, Merchant: "shopA"}
}
