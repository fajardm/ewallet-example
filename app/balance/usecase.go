package balance

import (
	"context"
	"github.com/fajardm/ewallet-example/app/balance/model"
	uuid "github.com/satori/go.uuid"
)

// Usecase represent the balance's usecase contract
type Usecase interface {
	GetBalanceByUserID(context.Context, uuid.UUID) (*model.Balance, error)
	GetBalanceHistoriesByUserID(context.Context, uuid.UUID) (model.BalanceHistories, error)
}
