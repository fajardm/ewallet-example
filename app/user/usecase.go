package user

import (
	"context"
	"github.com/fajardm/ewallet-example/app/user/model"
	uuid "github.com/satori/go.uuid"
)

// Usecase represent the user's usecase contract
type Usecase interface {
	Store(context.Context, model.User) error
	GetByID(context.Context, uuid.UUID) (*model.User, error)
	Update(context.Context, model.User) error
	Delete(context.Context, uuid.UUID) error
}
