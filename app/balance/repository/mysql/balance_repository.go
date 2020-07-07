package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fajardm/ewallet-example/app/balance"
	"github.com/fajardm/ewallet-example/app/balance/model"
	"github.com/fajardm/ewallet-example/app/errorcode"
	"github.com/fajardm/ewallet-example/database"
	uuid "github.com/satori/go.uuid"
)

const (
	// Table balances
	querySelectBalance = `
		SELECT 
			id,
			balance,
			user_id,
			created_by,
			created_at,
			updated_by,
			updated_at 
		FROM balances
	`
	queryInsertBalance = `
		INSERT INTO balances (
			id,
			balance,
			user_id,
			created_by,
			created_at
		) VALUES (?, ?, ?, ?, ?)
	`
	queryUpdateBalance = `
		UPDATE balances SET balance=?, updated_by=?, updated_at=? WHERE id=?
	`
	queryDeleteBalance = `
		DELETE FROM balances WHERE id=?
	`
	// Table balance_histories
	querySelectBalanceHistories = `
		SELECT 
			id,
			balance_before,
			balance_after,
			activity,
			type,
			ip,
			location,
			user_agent,
			balance_id,
			created_by,
			created_at,
			updated_by,
			updated_at 
		FROM balance_histories
	`
	queryInsertBalanceHistories = `
		INSERT INTO balance_histories (
			id,
			balance_before,
			balance_after,
			activity,
			type,
			ip,
			location,
			user_agent,
			balance_id,
			created_by,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	queryDeleteBalanceHistories = `
		DELETE FROM balance_histories WHERE balance_id=?
	`
)

type balanceRepository struct {
	db *database.MySQL
}

func NewBalanceRepository(conn *database.MySQL) balance.Repository {
	return &balanceRepository{db: conn}
}

func (b balanceRepository) TxStore(ctx context.Context, tx *sql.Tx, balance model.Balance) (err error) {
	_, err = tx.ExecContext(ctx, queryInsertBalance, balance.ID, balance.Balance, balance.UserID, balance.CreatedBy, balance.CreatedAt)
	return
}

func (b balanceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Balance, error) {
	q := querySelectBalance + " WHERE user_id=?"
	list, err := b.fetchContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return &list[0], nil
	}
	return nil, errorcode.ErrNotFound
}

func (b balanceRepository) TxUpdate(ctx context.Context, tx *sql.Tx, balance model.Balance) (err error) {
	res, err := tx.ExecContext(ctx, queryUpdateBalance, balance.Balance, balance.UpdatedBy, balance.UpdatedAt, balance.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affected > 1 {
		err = fmt.Errorf("Weird behaviour. Total affected: %d", affected)
		return
	}
	return
}

func (b balanceRepository) TxDelete(ctx context.Context, tx *sql.Tx, id uuid.UUID) (err error) {
	res, err := tx.ExecContext(ctx, queryDeleteBalance, id)
	if err != nil {
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affected > 1 {
		err = fmt.Errorf("Weird behaviour. Total affected: %d", affected)
		return
	}
	return
}

func (b balanceRepository) TxStoreBalanceHistory(ctx context.Context, tx *sql.Tx, history model.BalanceHistory) (err error) {
	_, err = tx.ExecContext(ctx, queryInsertBalanceHistories, history.ID, history.BalanceBefore, history.BalanceAfter, history.Activity, history.Type, history.IP, history.Location, history.UserAgent, history.BalanceID, history.CreatedBy, history.CreatedAt)
	return
}

func (b balanceRepository) FetchBalanceHistoriesByBalanceID(ctx context.Context, balanceID uuid.UUID) (model.BalanceHistories, error) {
	q := querySelectBalanceHistories + " WHERE balance_id = ? ORDER BY created_at DESC LIMIT 10"
	return b.fetchBalanceHistoriesContext(ctx, q, balanceID)
}

func (b balanceRepository) TxDeleteBalanceHistoriesByBalanceID(ctx context.Context, tx *sql.Tx, balanceID uuid.UUID) (err error) {
	res, err := tx.ExecContext(ctx, queryDeleteBalanceHistories, balanceID)
	if err != nil {
		return
	}
	_, err = res.RowsAffected()
	return
}

func (b balanceRepository) fetchContext(ctx context.Context, query string, args ...interface{}) (model.Balances, error) {
	rows, err := b.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(model.Balances, 0)
	for rows.Next() {
		r := model.Balance{}
		err = rows.Scan(&r.ID, &r.Balance, &r.UserID, &r.CreatedBy, &r.CreatedAt, &r.UpdatedBy, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func (b balanceRepository) fetchBalanceHistoriesContext(ctx context.Context, query string, args ...interface{}) (model.BalanceHistories, error) {
	rows, err := b.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(model.BalanceHistories, 0)
	for rows.Next() {
		r := model.BalanceHistory{}
		err = rows.Scan(&r.ID, &r.BalanceBefore, &r.BalanceAfter, &r.Activity, &r.Type, &r.IP, &r.Location, &r.UserAgent, &r.BalanceID, &r.CreatedBy, &r.CreatedAt, &r.UpdatedBy, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}
