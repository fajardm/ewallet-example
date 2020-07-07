package usecase

import (
	"context"
	"database/sql"
	"github.com/fajardm/ewallet-example/app/balance"
	_balanceModel "github.com/fajardm/ewallet-example/app/balance/model"
	"github.com/fajardm/ewallet-example/app/base"
	"github.com/fajardm/ewallet-example/app/errorcode"
	"github.com/fajardm/ewallet-example/app/user"
	"github.com/fajardm/ewallet-example/app/user/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type userUsecase struct {
	userRepository    user.Repository
	balanceRepository balance.Repository
	contextTimeout    time.Duration
}

func NewUserUsecase(userRepository user.Repository, balanceRepository balance.Repository, contextTimeout time.Duration) user.Usecase {
	return userUsecase{userRepository: userRepository, balanceRepository: balanceRepository, contextTimeout: contextTimeout}
}

func (u userUsecase) Login(ctx context.Context, username, email, password string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.GetByUsernameOrEmail(ctx, username, email)
	if err != nil {
		return nil, err
	}

	valid, err := user.ValidatePassword(password)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errorcode.ErrUnauthorized
	}

	return user, nil
}

func (u userUsecase) Store(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existed, _ := u.userRepository.GetByUsernameOrEmail(ctx, user.Username, user.Email)
	if existed != nil {
		return errorcode.ErrConflict
	}

	now := time.Now()
	balanceID := uuid.NewV4()
	activity := "initial balance"
	userBalance := _balanceModel.Balance{
		Model: base.Model{
			ID:        balanceID,
			CreatedBy: user.ID,
			CreatedAt: now,
		},
		UserID:  user.ID,
		Balance: 0,
		Histories: _balanceModel.BalanceHistories{
			_balanceModel.BalanceHistory{
				Model: base.Model{
					ID:        uuid.NewV4(),
					CreatedBy: user.ID,
					CreatedAt: now,
				},
				BalanceID:     balanceID,
				BalanceBefore: 0,
				BalanceAfter:  0,
				Activity:      &activity,
				Type:          _balanceModel.Credit,
				IP:            nil,
				Location:      nil,
				UserAgent:     nil,
			},
		},
	}

	return u.userRepository.WithTransaction(ctx, func(tx *sql.Tx) (err error) {
		if err = u.userRepository.TxStore(ctx, tx, user); err != nil {
			return err
		}
		if err = u.balanceRepository.TxStore(ctx, tx, userBalance); err != nil {
			return err
		}
		return err
	})
}

func (u userUsecase) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepository.GetByID(ctx, id)
}

func (u userUsecase) Update(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existed, _ := u.GetByID(ctx, user.ID)
	if existed == nil {
		return errorcode.ErrNotFound
	}

	return u.userRepository.Update(ctx, user)
}

func (u userUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existed, _ := u.GetByID(ctx, id)
	if existed == nil {
		return errorcode.ErrNotFound
	}

	return u.userRepository.WithTransaction(ctx, func(tx *sql.Tx) (err error) {
		if err = u.balanceRepository.TxDeleteByUserID(ctx, tx, id); err != nil {
			return err
		}
		if err = u.userRepository.TxDelete(ctx, tx, id); err != nil {
			return err
		}
		return err
	})
}
