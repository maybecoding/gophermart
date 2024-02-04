package entity

import (
	"errors"
	"time"
)

type BonusAmount float32

type BonusBalance struct {
	Available BonusAmount
	Withdrawn BonusAmount
}

type BonusProcessedAt time.Time

type BonusWithdraw struct {
	Amount      BonusAmount
	Order       OrderNumber
	ProcessedAt time.Time
}

var (
	ErrBonusNotEnough = errors.New("not enough bonuses")
)

type AccrualInfo struct {
	Order   OrderNumber
	Status  OrderStatus
	Accrual BonusAmount
}
