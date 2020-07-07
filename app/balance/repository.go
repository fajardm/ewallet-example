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
	TxUpdate(context.Context, *sql.Tx, model.Balance) error
	TxDelete(context.Context, *sql.Tx, uuid.UUID) error
	TxStoreBalanceHistory(context.Context, *sql.Tx, model.BalanceHistory) error
	FetchBalanceHistoriesByBalanceID(context.Context, uuid.UUID) (model.BalanceHistories, error)
	TxDeleteBalanceHistoriesByBalanceID(context.Context, *sql.Tx, uuid.UUID) error
	WithTransaction(context.Context, func(tx *sql.Tx) error) error
}
