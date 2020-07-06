package model

import (
	"github.com/fajardm/ewallet-example/app/base"
	"github.com/fajardm/ewallet-example/validator"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Input struct {
	Username    string `json:"username" validate:"required,max=45"`
	Email       string `json:"email" validate:"required,email,max=128"`
	MobilePhone string `json:"mobile_phone" validate:"required,max=13"`
	Password    string `json:"password" validate:"required,max=10"`
}

func (i Input) Validate() error {
	return validator.Validate().Struct(i)
}

func (i Input) NewUser() (*User, error) {
	hashedPassword, err := GeneratePassword(i.Password)
	if err != nil {
		return nil, err
	}
	id := uuid.NewV4()
	return &User{
		Model: base.Model{
			ID:        id,
			CreatedBy: id,
			CreatedAt: time.Now(),
		},
		Username:       i.Username,
		Email:          i.Email,
		MobilePhone:    i.MobilePhone,
		HashedPassword: hashedPassword,
	}, nil
}
