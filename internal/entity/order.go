package entity

import (
	"errors"
	"time"
)

type (
	OrderNumber  string
	OrderStatus  string
	OrderAccrual int

	Order struct {
		UserID     UserID
		Number     OrderNumber
		Status     OrderStatus
		Accrual    OrderAccrual
		UploadedAt time.Time
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
	ErrOrderNumberAlreadyLoadeed     = errors.New("order number is already loaded")
)
