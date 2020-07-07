package usecase

import (
	"context"
	"github.com/fajardm/ewallet-example/app/balance"
	"github.com/fajardm/ewallet-example/app/balance/model"
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
