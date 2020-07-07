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
		return
	}
	senderActivity := fmt.Sprintf("transfer amount %d to %s", amount, m.BalanceReceiver.UserID)
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
		return
	}
	rActivity := fmt.Sprintf("retrieve amount %d from %s", amount, m.BalanceSender.UserID)
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
