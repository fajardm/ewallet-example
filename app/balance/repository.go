package balance

import (
	"context"
	"database/sql"
	"github.com/fajardm/ewallet-example/app/balance/model"
	uuid "github.com/satori/go.uuid"
)

// Repository represent the balance's repository contract
type Repository interface {
	TxStore(context.Context, *sql.Tx, model.Balance) error
	GetByUserID(context.Context, uuid.UUID) (*model.Balance, error)
	Update(context.Context, model.Balance) error
	TxDeleteByUserID(context.Context, *sql.Tx, uuid.UUID) error
	FetchBalanceHistoriesByUserID(context.Context, uuid.UUID) (model.BalanceHistories, error)
}
