package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fajardm/ewallet-example/app/errorcode"
	"github.com/fajardm/ewallet-example/app/user"
	"github.com/fajardm/ewallet-example/app/user/model"
	"github.com/fajardm/ewallet-example/database"
	uuid "github.com/satori/go.uuid"
)

const (
	// Table users
	querySelectUser = `
		SELECT 
			id,
			username,
			email,
			mobile_phone,
			hashed_password,
			created_by,
			created_at,
			updated_by,
			updated_at 
		FROM users
	`
	queryInsertUser = `
		INSERT INTO users (
			id,
			username,
			email,
			mobile_phone,
			hashed_password,
			created_by,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	queryUpdateUser = `
		UPDATE users SET email=?, hashed_password=?, updated_by=?, updated_at=? WHERE id=?
	`
	queryDeleteUser = `
		DELETE FROM users WHERE id=?
	`
)

type userRepository struct {
	db *database.MySQL
}

func NewUserRepository(conn *database.MySQL) user.Repository {
	return &userRepository{db: conn}
}

func (u userRepository) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return u.db.WithTransaction(ctx, fn)
}

func (u userRepository) TxStore(ctx context.Context, tx *sql.Tx, user model.User) error {
	_, err := tx.ExecContext(ctx, queryInsertUser, user.ID, user.Username, user.Email, user.MobilePhone, user.HashedPassword, user.CreatedBy, user.CreatedAt)
	return err
}

func (u userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	q := querySelectUser + " WHERE id=?"
	list, err := u.fetchContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return &list[0], nil
	}
	return nil, errorcode.ErrNotFound
}

func (u userRepository) GetByUsernameOrEmail(ctx context.Context, username string, email string) (*model.User, error) {
	q := querySelectUser + " WHERE username=? OR email=?"
	list, err := u.fetchContext(ctx, q, username, email)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return &list[0], nil
	}
	return nil, errorcode.ErrNotFound
}

func (u userRepository) Update(ctx context.Context, user model.User) (err error) {
	res, err := u.db.ExecContext(ctx, queryUpdateUser, user.Email, user.HashedPassword, user.UpdatedBy, user.UpdatedAt, user.ID)
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

func (u userRepository) TxDelete(ctx context.Context, tx *sql.Tx, id uuid.UUID) (err error) {
	res, err := tx.ExecContext(ctx, queryDeleteUser, id)
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

func (u userRepository) fetchContext(ctx context.Context, query string, args ...interface{}) (model.Users, error) {
	rows, err := u.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(model.Users, 0)
	for rows.Next() {
		r := model.User{}
		err = rows.Scan(&r.ID, &r.Username, &r.Email, &r.MobilePhone, &r.HashedPassword, &r.CreatedBy, &r.CreatedAt, &r.UpdatedBy, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}
