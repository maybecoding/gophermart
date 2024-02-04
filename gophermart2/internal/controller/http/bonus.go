package http

import (
	"github.com/gin-gonic/gin"
	"gophermart/internal/entity"
	"gophermart/internal/usecase"
	"gophermart/pkg/logger"
	"net/http"
	"time"
)

type BonusRoutes struct {
	uc usecase.Bonus
}

func bonusRoutes(r *gin.RouterGroup, uc usecase.Bonus) {
	br := &BonusRoutes{uc}

	{
		r.GET("balance", br.GetBalance)
		r.GET("withdrawals", br.GetWithdrawal)
	}
}

type resBalance struct {
	Available entity.BonusAmount `json:"current"`
	Withdrawn entity.BonusAmount `json:"withdrawn"`
}

// GetBalance godoc
// @Summary      Get Balance
// @Description  Returns balance on user. current is all bonuses minus withdrawn
// @Tags         bonus
// @Accept       text/plain
// @Produce      json
// @Success      200   {object} resBalance
// @Failure      401   "No Content"
// @Failure      500   {object}  response
// @Router       /balance [get]
// @Security     JWT
// @SecurityDefinitions JWT
func (u *BonusRoutes) GetBalance(c *gin.Context) {
	userID, ok := getUser(c)
	if !ok {
		return
	}

	balance, err := u.uc.Get(c, userID)
	if err != nil {
		logger.Error().Err(err).Msg("http - BonusRoutes - Get - u.uc.Get")
		errorResponse(c, err, http.StatusInternalServerError)
		return
	}
	res := resBalance{
		Available: balance.Available,
		Withdrawn: balance.Withdrawn,
	}
	c.JSON(http.StatusOK, res)
}

type resWithdraw struct {
	Order       entity.OrderNumber `json:"order"`
	Amount      entity.BonusAmount `json:"sum"`
	ProcessedAt time.Time          `json:"processed_at"`
}

// GetWithdrawal godoc
// @Summary      Get Withdrawals
// @Description  Returns orders for withdrawals
// @Tags         bonus
// @Accept       text/plain
// @Produce      json
// @Success      200   {object} []resWithdraw
// @Failure      204   "No Content"
// @Failure      401   "No Content"
// @Failure      500   {object}  response
// @Router       /withdrawals [get]
// @Security     JWT
// @SecurityDefinitions JWT
func (u *BonusRoutes) GetWithdrawal(c *gin.Context) {
	userID, ok := getUser(c)
	if !ok {
		return
	}

	ws, err := u.uc.GetWithdrawals(c, userID)
	if err != nil {
		logger.Error().Err(err).Msg("http - BonusRoutes - GetWithdrawal - u.uc.GetWithdrawals")
		errorResponse(c, err, http.StatusInternalServerError)
		return
	}
	if len(ws) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	res := make([]resWithdraw, 0, len(ws))
	for _, w := range ws {
		res = append(res, resWithdraw{
			Order:       w.Order,
			Amount:      w.Amount,
			ProcessedAt: w.ProcessedAt,
		})
	}
	c.JSON(http.StatusOK, res)
}
