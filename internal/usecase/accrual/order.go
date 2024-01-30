package accrual

import (
	"encoding/json"
	"fmt"
	"gophermart/internal/config"
	"gophermart/internal/entity"
	"gophermart/pkg/logger"
	"io"
	"net/http"
)

type OrderAccrual struct {
	cfg config.AccrualSystem
}

func NewOrder(cfg config.AccrualSystem) *OrderAccrual {
	return &OrderAccrual{cfg}
}

type resAccrualInfo struct {
	Order   entity.OrderNumber `json:"order"`
	Status  entity.OrderStatus `json:"status"`
	Accrual entity.BonusAmount `json:"accrual"`
}

func (oa *OrderAccrual) GetStatus(orderNum entity.OrderNumber) (*entity.AccrualInfo, error) {
	endpoint := fmt.Sprintf("http://%s/api/orders/%s", oa.cfg.Address, orderNum)
	logger.Debug().Str("endpoint", endpoint).Msg("trying to get accrual")
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("OrderAccrual - GetStatus - http.Get: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("OrderAccrual - GetStatus - io.ReadAll: %w", err)
		}
		return nil, fmt.Errorf("fetch for %s, got invalid code %d, error message %s", endpoint, resp.StatusCode, string(errBody))
	}

	dec := json.NewDecoder(resp.Body)

	res := resAccrualInfo{}
	err = dec.Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("OrderAccrual - GetStatus - dec.Decode: %w", err)
	}
	accrualInfo := entity.AccrualInfo{
		Order:   res.Order,
		Status:  res.Status,
		Accrual: res.Accrual,
	}
	return &accrualInfo, nil
}
