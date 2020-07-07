package model

import (
	"github.com/fajardm/ewallet-example/app/base"
	uuid "github.com/satori/go.uuid"
)

// Balance is balance model
type Balance struct {
	base.Model
	UserID    uuid.UUID        `json:"user_id"`
	Balance   float64          `json:"balance"`
	Histories BalanceHistories `json:"-"`
}

// Balances is list of balance model
type Balances []Balance

// BalanceHistory is balance history model
type BalanceHistory struct {
	base.Model
	BalanceID     uuid.UUID              `json:"balance_id"`
	BalanceBefore float64                `json:"balance_before"`
	BalanceAfter  float64                `json:"balance_after"`
	Activity      *string                `json:"activity"`
	Type          UserBalanceHistoryType `json:"type"`
	IP            *string                `json:"ip"`
	Location      *string                `json:"location"`
	UserAgent     *string                `json:"user_agent"`
}

// BalanceHistories is list of balance history model
type BalanceHistories []BalanceHistory
