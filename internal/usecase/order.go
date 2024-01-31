package usecase

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/entity"
	"gophermart/pkg/logger"
	"time"
)

type OrderUseCase struct {
	repo      OrderRepo
	repoBonus BonusRepo
	numAlg    OrderNumAlg
	accrual   OrderAccrual
}

func NewOrder(repo OrderRepo, repoBonus BonusRepo, numAlg OrderNumAlg, accrual OrderAccrual) *OrderUseCase {
	return &OrderUseCase{
		repo:      repo,
		repoBonus: repoBonus,
		numAlg:    numAlg,
		accrual:   accrual,
	}
}

func (uc *OrderUseCase) add(ctx context.Context, order entity.Order) (*entity.Order, error) {
	// Проверяем, что номер заказа допустим
	ok, err := uc.numAlg.Check(order.Number)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - Add - uc.numAlg.Check: %w", err)
	}
	if !ok {
		return nil, entity.ErrOrderNumberFormat
	}
	// Получаем заказ если он уже есть
	existsOrder, err := uc.repo.Get(ctx, order.Number)
	if err != nil {
		if !errors.Is(err, entity.ErrOrderNotFound) {
			return nil, fmt.Errorf("OrderUseCase - Add - uc.repo.Get: %w", err)
		}
	} else {
		if existsOrder.UserID != order.UserID {
			return nil, entity.ErrOrderNumberOwnedByAnotherUser
		}
		if existsOrder.UserID == order.UserID {
			return nil, entity.ErrOrderNumberAlreadyLoaded
		}
	}
	// Добавляем новый номер заказа
	outO, err := uc.repo.Add(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - Add - uc.repo.Add: %w", err)
	}

	return outO, nil
}

func (uc *OrderUseCase) AddNew(ctx context.Context, userID entity.UserID, number entity.OrderNumber) (*entity.Order, error) {

	o := entity.Order{
		UserID:     userID,
		Number:     number,
		Status:     entity.OrderStatusNew,
		Accrual:    0,
		UploadedAt: time.Now(),
	}
	return uc.add(ctx, o)
}

func (uc *OrderUseCase) AddForBonuses(ctx context.Context, userID entity.UserID, number entity.OrderNumber, amount entity.BonusAmount) (*entity.Order, error) {
	accrualAt := time.Now()
	// Проверяем баланс
	balance, err := uc.repoBonus.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - AddForBonuses - uc.repoBonus.GetBalance: %w", err)
	}
	if balance.Available < amount {
		return nil, entity.ErrBonusNotEnough
	}
	o := entity.Order{
		UserID:             userID,
		Number:             number,
		Status:             entity.OrderStatusProcessed,
		Accrual:            -amount,
		AccrualProcessedAt: &accrualAt,
		UploadedAt:         time.Now(),
	}
	return uc.add(ctx, o)
}

func (uc *OrderUseCase) GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error) {
	orders, err := uc.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - GetByUser - uc.repo.GetByUser: %w", err)
	}
	return orders, nil
}

func (uc *OrderUseCase) RunAccrualRefresh(ctx context.Context) {
	// Получаем список к обновлению
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(20 * time.Second):
			uc.accrualRefreshAll(ctx)
		}
	}
}

func (uc *OrderUseCase) accrualRefreshAll(ctx context.Context) {
	// Получаем список к обновлению
	orders, err := uc.repo.GetUnAccrued(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("OrderUseCase - accrualRefreshAll - uc.repo.GetUnAccrued")
		return
	}
	for _, order := range orders {
		// Получаем информацию о начислении бонусов
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Millisecond):
			uc.accrualRefresh(ctx, order)
		}
	}
}

func (uc *OrderUseCase) accrualRefresh(ctx context.Context, order entity.Order) {
	accInf, err := uc.accrual.GetStatus(order.Number)
	if err != nil {
		logger.Error().Err(err).Msg("OrderUseCase - accrualRefresh - uc.accrual.GetStatus")
		return
	}
	if order.Status == accInf.Status {
		logger.Debug().Interface("order", order).Msg("status is same - return")
		return
	}
	_, err = uc.repo.Accrual(ctx, *accInf)
	if err != nil {
		logger.Error().Err(err).Msg("OrderUseCase - accrualRefresh - uc.repo.Accrual")
		return
	}
}
