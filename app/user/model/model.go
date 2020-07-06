package model

import (
	"github.com/fajardm/ewallet-example/app/base"
	"golang.org/x/crypto/bcrypt"
)

// User is user model
type User struct {
	base.Model
	Username       string `json:"username"`
	Email          string `json:"email"`
	MobilePhone    string `json:"mobile_phone"`
	HashedPassword []byte `json:"-"`
}

// Users represent list of User
type Users []User

// ValidatePassword will check if passwords are matched
func (u User) ValidatePassword(password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}

// GeneratePassword will generate a hashed password for us based on the user input
func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
