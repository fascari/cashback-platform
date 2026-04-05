package testdata

import cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"

const UserID int64 = 10

func Cashbacks() []cashdomain.Cashback {
	return []cashdomain.Cashback{
		{ID: 1, UserID: UserID, Amount: 5.0},
		{ID: 2, UserID: UserID, Amount: 3.0},
	}
}
