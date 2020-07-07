package user

import (
	"context"
	"database/sql"
	"github.com/fajardm/ewallet-example/app/user/model"
	uuid "github.com/satori/go.uuid"
)

// Repository represent the user's repository contract
type Repository interface {
	TxStore(context.Context, *sql.Tx, model.User) error
	GetByID(context.Context, uuid.UUID) (*model.User, error)
	GetByUsernameOrEmail(context.Context, string, string) (*model.User, error)
	Update(context.Context, model.User) error
	TxDelete(context.Context, *sql.Tx, uuid.UUID) error
	WithTransaction(context.Context, func(tx *sql.Tx) error) error
}
