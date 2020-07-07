package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fajardm/ewallet-example/app/balance"
	"github.com/fajardm/ewallet-example/app/balance/model"
	"github.com/fajardm/ewallet-example/app/base"
	uuid "github.com/satori/go.uuid"
	"time"
)

type balanceUsecase struct {
	balanceRepository balance.Repository
	contextTimeout    time.Duration
}

func NewBalanceUsecase(balanceRepository balance.Repository, contextTimeout time.Duration) balance.Usecase {
	return balanceUsecase{balanceRepository: balanceRepository, contextTimeout: contextTimeout}
}

func (b balanceUsecase) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*model.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	return b.balanceRepository.GetByUserID(ctx, userID)
}

func (b balanceUsecase) GetBalanceHistoriesByUserID(ctx context.Context, userID uuid.UUID) (model.BalanceHistories, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	balance, err := b.balanceRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return b.balanceRepository.FetchBalanceHistoriesByBalanceID(ctx, balance.ID)
}

func (b balanceUsecase) TransferBalance(ctx context.Context, fromUserID, toUserID uuid.UUID, amount float64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	now := time.Now()

	sender, err := b.balanceRepository.GetByUserID(ctx, fromUserID)
	if err != nil {
		return err
	}
	bfs := sender.Balance
	if err := sender.Reduce(amount); err != nil {
		return err
	}
	senderActivity := fmt.Sprintf("transfer amount %d to %s", amount, toUserID)
	sender.Histories = model.BalanceHistories{
		model.BalanceHistory{
			Model: base.Model{
				ID:        uuid.NewV4(),
				CreatedBy: sender.UserID,
				CreatedAt: now,
			},
			BalanceID:     sender.ID,
			BalanceBefore: bfs,
			BalanceAfter:  sender.Balance,
			Activity:      &senderActivity,
			Type:          model.Debit,
			IP:            nil,
			Location:      nil,
			UserAgent:     nil,
		},
	}

	reciever, err := b.balanceRepository.GetByUserID(ctx, toUserID)
	if err != nil {
		return err
	}
	bfr := reciever.Balance
	if err := reciever.Add(amount); err != nil {
		return err
	}
	rActivity := fmt.Sprintf("retrieve amount %d from %s", amount, fromUserID)
	reciever.Histories = model.BalanceHistories{
		model.BalanceHistory{
			Model: base.Model{
				ID:        uuid.NewV4(),
				CreatedBy: reciever.UserID,
				CreatedAt: now,
			},
			BalanceID:     reciever.ID,
			BalanceBefore: bfr,
			BalanceAfter:  reciever.Balance,
			Activity:      &rActivity,
			Type:          model.Credit,
			IP:            nil,
			Location:      nil,
			UserAgent:     nil,
		},
	}

	return b.balanceRepository.WithTransaction(ctx, func(tx *sql.Tx) (err error) {
		if err = b.balanceRepository.TxUpdate(ctx, tx, *sender); err != nil {
			return err
		}
		if err = b.balanceRepository.TxStoreBalanceHistory(ctx, tx, sender.Histories[0]); err != nil {
			return err
		}
		if err = b.balanceRepository.TxUpdate(ctx, tx, *reciever); err != nil {
			return err
		}
		if err = b.balanceRepository.TxStoreBalanceHistory(ctx, tx, reciever.Histories[0]); err != nil {
			return err
		}
		return
	})
}

func (b balanceUsecase) TopUp(ctx context.Context, userID uuid.UUID, amount float64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	now := time.Now()

	balance, err := b.balanceRepository.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	bfs := balance.Balance
	if err := balance.Add(amount); err != nil {
		return err
	}
	activity := fmt.Sprintf("topup amount %d", amount)
	balance.Histories = model.BalanceHistories{
		model.BalanceHistory{
			Model: base.Model{
				ID:        uuid.NewV4(),
				CreatedBy: balance.UserID,
				CreatedAt: now,
			},
			BalanceID:     balance.ID,
			BalanceBefore: bfs,
			BalanceAfter:  balance.Balance,
			Activity:      &activity,
			Type:          model.Credit,
			IP:            nil,
			Location:      nil,
			UserAgent:     nil,
		},
	}

	return b.balanceRepository.WithTransaction(ctx, func(tx *sql.Tx) (err error) {
		if err = b.balanceRepository.TxUpdate(ctx, tx, *balance); err != nil {
			return err
		}
		if err = b.balanceRepository.TxStoreBalanceHistory(ctx, tx, balance.Histories[0]); err != nil {
			return err
		}
		return
	})
}
