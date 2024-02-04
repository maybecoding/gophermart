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
	rTx      TxRepo
	rOrder   OrderRepo
	rBonus   BonusRepo
	rBalance BalanceRepo
	numAlg   OrderNumAlg
	accrual  OrderAccrual
}

func NewOrder(rTx TxRepo, order OrderRepo, bns BonusRepo, bls BalanceRepo, numAlg OrderNumAlg, accrual OrderAccrual) *OrderUseCase {
	return &OrderUseCase{
		rTx:      rTx,
		rOrder:   order,
		rBonus:   bns,
		rBalance: bls,
		numAlg:   numAlg,
		accrual:  accrual,
	}
}

func (uc *OrderUseCase) Add(ctx context.Context, userID entity.UserID, number entity.OrderNumber) error {
	// Проверяем, что номер заказа допустим
	ok, err := uc.numAlg.Check(number)
	if err != nil {
		return fmt.Errorf("OrderUseCase - Add - uc.numAlg.Check: %w", err)
	}
	if !ok {
		return entity.ErrOrderNumberFormat
	}
	// Получаем заказ если он уже есть
	existsOrder, err := uc.rOrder.Get(ctx, number)
	if err != nil {
		if !errors.Is(err, entity.ErrOrderNotFound) {
			return fmt.Errorf("OrderUseCase - Add - uc.order.Get: %w", err)
		}
	} else {
		if existsOrder.UserID != userID {
			return entity.ErrOrderNumberOwnedByAnotherUser
		}
		if existsOrder.UserID == userID {
			return entity.ErrOrderNumberAlreadyLoaded
		}
	}
	// Добавляем новый номер заказа
	o := entity.Order{
		UserID:     userID,
		Number:     number,
		Status:     entity.OrderStatusNew,
		UploadedAt: time.Now(),
	}
	logger.Debug().Interface("order", o).Msg("OrderUseCase - Add")
	err = uc.rOrder.Add(ctx, o)

	if err != nil {
		return fmt.Errorf("OrderUseCase - Add - uc.order.Add: %w", err)
	}

	return nil
}

func (uc *OrderUseCase) AddForBonuses(ctx context.Context, userID entity.UserID, number entity.OrderNumber, amount entity.BonusAmount) error {
	err := uc.rTx.WithTx(ctx, func(ctx context.Context) error {
		// Проверяем баланс
		balance, err := uc.rBalance.Get(ctx, userID)
		if err != nil {
			return fmt.Errorf("OrderUseCase - AddForBonuses - uc.bonus.Get: %w", err)
		}
		if balance.Available < amount {
			return entity.ErrBonusNotEnough
		}
		// Добавляем заказ
		err = uc.Add(ctx, userID, number)
		if err != nil {
			return fmt.Errorf("OrderUseCase - AddForBonuses - uc.Add: %w", err)
		}

		// Записываем списание бонусов в заказе
		err = uc.rBonus.Set(ctx, number, -amount, entity.OrderStatusNew)
		if err != nil {
			return fmt.Errorf("OrderUseCase - AddForBonuses - uc.rBonus.Set: %w", err)
		}

		// Списываем баланс
		newBalance := entity.BonusBalance{
			Available: balance.Available - amount,
			Withdrawn: balance.Withdrawn + amount,
		}
		err = uc.rBalance.Set(ctx, userID, newBalance)
		if err != nil {
			return fmt.Errorf("OrderUseCase - AddForBonuses - uc.rBalance.Set: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("OrderUseCase - AddForBonuses - uc.rTx.WithTx: %w", err)
	}

	return err
}

func (uc *OrderUseCase) GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error) {
	orders, err := uc.rOrder.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("OrderUseCase - GetByUser - uc.order.GetByUser: %w", err)
	}
	return orders, nil
}

func (uc *OrderUseCase) RunAccrualRefresh(ctx context.Context) {
	// Получаем список к обновлению
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			uc.accrualRefreshAll(ctx)
		}
	}
}

func (uc *OrderUseCase) accrualRefreshAll(ctx context.Context) {
	// Получаем список к обновлению
	orders, err := uc.rOrder.GetUnAccrued(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("OrderUseCase - accrualRefreshAll - uc.order.GetUnAccrued")
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

func (uc *OrderUseCase) accrualRefresh(ctx context.Context, order entity.OrderAccrual) {
	accInf, err := uc.accrual.GetStatus(order.Number)
	if err != nil {
		logger.Error().Err(err).Msg("OrderUseCase - accrualRefresh - uc.accrual.GetStatus")
		return
	}
	if order.Status == accInf.Status {
		logger.Debug().Interface("order", order).Msg("status is same - return")
		return
	}

	err = uc.rTx.WithTx(ctx, func(ctx context.Context) error {
		// Получаем предыдущее значение accrual и user_id
		currOrder, err := uc.rOrder.Get(ctx, order.Number)
		if err != nil {
			return fmt.Errorf("OrderUseCase - accrualRefresh - uc.rOrder.Get: %w", err)
		}
		// Устанавливаем начисление бонусов в заказ
		err = uc.rBonus.Set(ctx, order.Number, accInf.Accrual, entity.OrderStatusProcessed)
		if err != nil {
			return fmt.Errorf("OrderUseCase - accrualRefresh - uc.rBonus.Set: %w", err)
		}
		// Получаем текущй баланс
		currBalance, err := uc.rBalance.Get(ctx, currOrder.UserID)
		if err != nil {
			return fmt.Errorf("OrderUseCase - accrualRefresh - uc.rBalance.Get: %w", err)
		}
		// Устанавливаем баланс
		newBalance := entity.BonusBalance{
			Available: currBalance.Available + accInf.Accrual - currOrder.Accrual,
			Withdrawn: currBalance.Withdrawn,
		}
		err = uc.rBalance.Set(ctx, currOrder.UserID, newBalance)
		if err != nil {
			return fmt.Errorf("OrderUseCase - accrualRefresh - uc.rBalance.Set: %w", err)
		}
		return nil
	})
	if err != nil {
		logger.Error().Err(err).Msg("OrderUseCase - accrualRefresh - uc.rTx.WithTx")
	}
}
