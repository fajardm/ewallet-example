package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
)

// ErrInvalidUserBalanceHistoryType represent error when invalid UserBalanceHistoryType
var ErrInvalidUserBalanceHistoryType = errors.New("InvalidUserBalanceHistoryType")

type UserBalanceHistoryType int

const (
	// Credit represent credit enum
	Credit = 1 + iota
	// Debit represent debit enum
	Debit
)

// UserBalanceHistoryTypeFromString will converts a string to a UserBalanceHistoryType, will return UserBalanceHistoryType if string is
// valid representation of UserBalanceHistoryType, or error otherwise
func UserBalanceHistoryTypeFromString(s string) (res UserBalanceHistoryType, err error) {
	switch s {
	case "credit":
		res = Credit
	case "debit":
		res = Debit
	default:
		err = errors.WithMessagef(ErrInvalidUserBalanceHistoryType, "invalid value: %s", s)
	}
	return
}

// MarshalText is the custom marshalling for UserBalanceHistoryType. With this when marshalling to json
// UserBalanceHistoryType will be shown as its string representation instead of int
func (u UserBalanceHistoryType) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// String returns the string representation of UserBalanceHistoryType
func (u UserBalanceHistoryType) String() string {
	var s string
	switch u {
	case Credit:
		s = "debit"
	case Debit:
		s = "credit"
	}
	return s
}

// Value transforms UserBalanceHistoryType to its value for its column in database (MySQL)
func (u UserBalanceHistoryType) Value() (driver.Value, error) {
	return u.String(), nil
}

//// Scan transforms MySQL enum column value for type column to UserBalanceHistoryType
func (u *UserBalanceHistoryType) Scan(value interface{}) error {
	b, ok := value.([]uint8)
	if !ok {
		return fmt.Errorf("expecting a []uint8 found %T, in string: %s", value, value)
	}
	s := string(b)
	st, err := UserBalanceHistoryTypeFromString(s)
	if err != nil {
		return err
	}
	*u = st
	return nil
}
