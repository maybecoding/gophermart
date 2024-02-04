package numalg

import (
	"fmt"
	"gophermart/internal/entity"
	"gophermart/pkg/luna"
)

type NumAlg struct{}

func New() *NumAlg {
	return &NumAlg{}
}

func (ona *NumAlg) Check(num entity.OrderNumber) (isCorrect bool, err error) {
	isCorrect, err = luna.Check(string(num))
	if err != nil {
		return false, fmt.Errorf("NumAlg - Check - luna.Check: %w", err)
	}
	return isCorrect, nil
}
