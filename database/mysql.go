package database

import (
	"context"
	"database/sql"
)

type MySQL struct {
	*sql.DB
}

func (m MySQL) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}
