package entity

import (
	"errors"
	"time"
)

type (
	OrderNumber string
	OrderStatus string

	Order struct {
		UserID             UserID
		Number             OrderNumber
		Status             OrderStatus
		Accrual            BonusAmount
		AccrualProcessedAt *time.Time
		UploadedAt         time.Time
	}

	OrderAccrual struct {
		Number OrderNumber
		Status OrderStatus
	}
)

var (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

var (
	ErrOrderNotFound                 = errors.New("order not found")
	ErrOrderNumberFormat             = errors.New("incorrect order number format")
	ErrOrderNumberOwnedByAnotherUser = errors.New("order number is used by another user")
	ErrOrderNumberAlreadyLoaded      = errors.New("order number is already loaded")
	ErrGracefulShutdown              = errors.New("graceful shutdown committed")
)
