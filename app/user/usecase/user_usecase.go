package usecase

import (
	"context"
	"github.com/fajardm/ewallet-example/app/errorcode"
	"github.com/fajardm/ewallet-example/app/user"
	"github.com/fajardm/ewallet-example/app/user/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type userUsecase struct {
	userRepository user.Repository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository user.Repository, contextTimeout time.Duration) user.Usecase {
	return userUsecase{userRepository: userRepository, contextTimeout: contextTimeout}
}

func (u userUsecase) Store(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existed, _ := u.userRepository.GetByUsernameOrEmail(ctx, user.Username, user.Email)
	if existed != nil {
		return errorcode.ErrConflict
	}

	return u.userRepository.Store(ctx, user)
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

	return u.userRepository.Delete(ctx, id)
}
