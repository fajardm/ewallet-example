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
		UPDATE balances SET balance=?, updated_by=?, updated_at=? WHERE user_id=?
	`
	queryDeleteBalance = `
		DELETE FROM balances WHERE user_id=?
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
	if err != nil {
		return
	}

	for _, h := range balance.Histories {
		if _, err = tx.ExecContext(ctx, queryInsertBalanceHistories, h.ID, h.BalanceBefore, h.BalanceAfter, h.Activity, h.Type, h.IP, h.Location, h.UserAgent, h.BalanceID, h.CreatedBy, h.CreatedAt); err != nil {
			return
		}
	}
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

func (b balanceRepository) Update(ctx context.Context, m model.Balance) error {
	panic("implement me")
}

func (b balanceRepository) TxDeleteByUserID(ctx context.Context, tx *sql.Tx, userID uuid.UUID) (err error) {
	balance, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if err = b.txDeleteBalanceHistoriesByBalanceID(ctx, tx, balance.ID); err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, queryDeleteBalance, userID)
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

func (b balanceRepository) FetchBalanceHistoriesByUserID(ctx context.Context, uuid uuid.UUID) (model.BalanceHistories, error) {
	panic("implement me")
}

func (b balanceRepository) txDeleteBalanceHistoriesByBalanceID(ctx context.Context, tx *sql.Tx, balanceID uuid.UUID) (err error) {
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