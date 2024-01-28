package usecase

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/entity"
	"time"
)

type OrderUseCase struct {
	repo   OrderRepo
	numAlg OrderNumAlg
}

func NewOrder(repo OrderRepo, numAlg OrderNumAlg) *OrderUseCase {
	return &OrderUseCase{
		repo:   repo,
		numAlg: numAlg,
	}
}

func (uc *OrderUseCase) Add(ctx context.Context, userID entity.UserID, number entity.OrderNumber) (*entity.Order, error) {
	// Проверяем, что номер заказа допустим
	ok, err := uc.numAlg.Check(number)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - Add - uc.numAlg.Check: %w", err)
	}
	if !ok {
		return nil, entity.ErrOrderNumberFormat
	}
	// Получаем заказ если он уже есть
	order, err := uc.repo.Get(ctx, number)
	if err != nil {
		if !errors.Is(err, entity.ErrOrderNotFound) {
			return nil, fmt.Errorf("OrderUseCase - Add - uc.repo.Get: %w", err)
		}
	} else {
		if order.UserID != userID {
			return order, entity.ErrOrderNumberOwnedByAnotherUser
		}
		if order.UserID == userID {
			return order, entity.ErrOrderNumberAlreadyLoadeed
		}
	}
	// Добавляем новый номер заказа
	o := entity.Order{
		UserID:     userID,
		Number:     number,
		Status:     entity.OrderStatusNew,
		UploadedAt: time.Now(),
	}
	outO, err := uc.repo.Add(ctx, o)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - Add - uc.repo.Add: %w", err)
	}

	return outO, nil
}

func (uc *OrderUseCase) GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error) {
	orders, err := uc.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - GetByUser - uc.repo.GetByUser: %w", err)
	}
	return orders, nil
}
