package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gophermart/internal/entity"
	"gophermart/internal/usecase"
	"gophermart/pkg/logger"
	"io"
	"net/http"
	"time"
)

type OrderRoutes struct {
	uc usecase.Order
}

func orderRoutes(r *gin.RouterGroup, uc usecase.Order) {
	orr := &OrderRoutes{uc}

	{
		r.POST("orders", orr.OrderAdd)
		r.GET("orders", orr.GetOrders)
	}
}

type respOrder struct {
	UserID     entity.UserID       `json:"-"`
	Number     entity.OrderNumber  `json:"number"`
	Status     entity.OrderStatus  `json:"status"`
	Accrual    entity.OrderAccrual `json:"accrual,omitempty"`
	UploadedAt string              `json:"uploaded_at"`
}

// OrderAdd godoc
// @Summary      Order Create
// @Description  Creates New Order
// @Tags         order
// @Accept       text/plain
// @Produce      json
// @Param        request	body	string	true	"order ID"
// @Success      200   {object} respOrder
// @Success      202   {object} respOrder
// @Failure      400   {object}  response
// @Failure      401   "No Content"
// @Failure      409   {object}  response
// @Failure      422   {object}  response
// @Failure      500   {object}  response
// @Router       /orders [post]
// @Security     JWT
// @SecurityDefinitions JWT
func (u *OrderRoutes) OrderAdd(c *gin.Context) {
	fmt.Println("hi")
	b, err := io.ReadAll(c.Request.Body)
	defer func() {
		_ = c.Request.Body.Close()
	}()
	number := entity.OrderNumber(b)
	if err != nil {
		logger.Error().Err(err).Msg("http - OrderRoutes - OrderAdd - io.ReadAll")
		errorResponse(c, err, http.StatusBadRequest)
		return
	}
	if number == "" {
		logger.Error().Err(err).Msg("http - OrderRoutes - OrderAdd - number==\"\"")
		errorResponse(c, errors.New("order number is blank"), http.StatusBadRequest)
		return
	}
	logger.Debug().Str("number", string(number)).Msg("Order number")
	userID, ok := getUser(c)
	if !ok {
		return
	}
	order, err := u.uc.Add(c, userID, number)
	if err == nil || err != nil && errors.Is(err, entity.ErrOrderNumberAlreadyLoadeed) {
		code := http.StatusAccepted
		if err != nil {
			code = http.StatusOK
		}

		resp := respOrder{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		}
		c.JSON(code, resp)
		return
	}

	errCode := http.StatusInternalServerError
	if errors.Is(err, entity.ErrOrderNumberFormat) {
		errCode = http.StatusUnprocessableEntity
	} else if errors.Is(err, entity.ErrOrderNumberOwnedByAnotherUser) {
		errCode = http.StatusConflict
	}
	logger.Error().Err(err).Int("error code", errCode).Msg("http - OrderRoutes - OrderAdd - u.uc.Add")
	errorResponse(c, err, errCode)
	return
}

type respOrderList []respOrder

// GetOrders godoc
// @Summary      Get user orders
// @Description  Returns orders authorised by user
// @Tags         order
// @Accept       text/plain
// @Produce      json
// @Success      200 {object} respOrderList
// @Success      204   {object} respOrderList
// @Failure      401  "No Content"
// @Failure      500  {object}  response
// @Router       /orders [get]
// @Security     JWT
// @SecurityDefinitions JWT
func (u *OrderRoutes) GetOrders(c *gin.Context) {
	userID, ok := getUser(c)
	if !ok {
		return
	}
	orders, err := u.uc.GetByUser(c, userID)
	if err != nil {
		logger.Error().Err(err).Msg("http - OrderRoutes - GetOrders - u.uc.GetByUser")
		errorResponse(c, err, http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	res := make(respOrderList, 0, len(orders))
	for _, order := range orders {
		res = append(res, respOrder{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, res)
}
