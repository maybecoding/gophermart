package impl

import (
	"fmt"
	"gophermart/internal/entity"
	"gophermart/pkg/luna"
)

type OrderNumAlgImpl struct{}

func NewOrderNumAlgImpl() *OrderNumAlgImpl {
	return &OrderNumAlgImpl{}
}

func (ona *OrderNumAlgImpl) Check(num entity.OrderNumber) (isCorrect bool, err error) {
	isCorrect, err = luna.Check(string(num))
	if err != nil {
		return false, fmt.Errorf("OrderNumAlgImpl - Check - luna.Check: %w", err)
	}
	return isCorrect, nil
}
