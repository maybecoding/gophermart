package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gophermart/internal/entity"
	"gophermart/internal/usecase"
	"gophermart/pkg/logger"
	"net/http"
)

type UserRoutes struct {
	uc usecase.Auth
}

func userRoutes(r *gin.RouterGroup, uc usecase.Auth) {
	ur := &UserRoutes{uc}

	h := r.Group("/user")
	{
		h.POST("register", ur.Register)
		h.POST("login", ur.Login)
	}
}

type registerRequest struct {
	Login    entity.UserLogin    `json:"login" binding:"required"`
	Password entity.UserPassword `json:"password" binding:"required"`
}

// Register godoc
// @Summary      Register
// @Description  Register using login and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request	body  registerRequest	true	"login and password"
// @Success      200   "No Content"
// @Header Set-Cookie string true "sets cookie"
// @Failure      400  {object}  response
// @Failure      409  {object}  response
// @Failure      500  {object}  response
// @Router       /user/register [post]
func (u *UserRoutes) Register(c *gin.Context) {
	var req registerRequest
	err := c.BindJSON(&req)
	if err != nil {
		logger.Error().Err(err).Msg("http - UserRoutes - register - req validation")
		errorResponse(c, err, http.StatusBadRequest)
		return
	}
	token, err := u.uc.Register(c, req.Login, req.Password)
	if err != nil {
		logger.Error().Err(err).Msg("http - UserRoutes - register - u.uc.Register")
		errCode := http.StatusInternalServerError
		if errors.Is(err, entity.ErrUserNotAvailable) {
			errCode = http.StatusConflict
		}
		errorResponse(c, err, errCode)
		return
	}
	setAuthCookie(c, token)
	c.Status(http.StatusOK)
}

type loginRequest struct {
	Login    entity.UserLogin    `json:"login" binding:"required"`
	Password entity.UserPassword `json:"password" binding:"required"`
}

// Login godoc
// @Summary      Login
// @Description  Login using login and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request	body  registerRequest	true	"login and password"
// @Success      200   "No Content"
// @Header Set-Cookie string true "sets cookie"
// @Failure      400  {object}  response
// @Failure      401  {object}  response
// @Failure      500  {object}  response
// @Router       /user/login [post]
func (u *UserRoutes) Login(c *gin.Context) {
	var req loginRequest
	err := c.BindJSON(&req)
	if err != nil {
		logger.Error().Err(err).Msg("http - UserRoutes - login - req validation")
		errorResponse(c, err, http.StatusBadRequest)
		return
	}
	token, err := u.uc.Login(c, req.Login, req.Password)
	if err != nil {
		logger.Error().Err(err).Msg("http - UserRoutes - login - u.uc.Login")
		errCode := http.StatusInternalServerError
		if errors.Is(err, entity.ErrUserNotFound) || errors.Is(err, entity.ErrIncorrectPassword) {
			errCode = http.StatusUnauthorized
		}
		errorResponse(c, err, errCode)
		return
	}
	setAuthCookie(c, token)
	c.Status(http.StatusOK)
}
