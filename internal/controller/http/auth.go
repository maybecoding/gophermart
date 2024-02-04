package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gophermart/internal/entity"
	"gophermart/internal/usecase"
	"gophermart/pkg/logger"
	"net/http"
)

type AuthRoutes struct {
	uc usecase.Auth
}

func authRoutes(r *gin.RouterGroup, uc usecase.Auth) {
	ur := &AuthRoutes{uc}
	{
		r.POST("register", ur.Register)
		r.POST("login", ur.Login)
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
// @Header       200 {string} Set-Cookie "sets cookie"
// @Failure      400  {object}  response
// @Failure      409  {object}  response
// @Failure      500  {object}  response
// @Router       /register [post]
func (u *AuthRoutes) Register(c *gin.Context) {
	var req registerRequest
	err := c.BindJSON(&req)
	if err != nil {
		logger.Error().Err(err).Msg("http - AuthRoutes - register - req validation")
		errorResponse(c, err, http.StatusBadRequest)
		return
	}
	token, err := u.uc.Register(c, req.Login, req.Password)
	if err != nil {
		logger.Error().Err(err).Msg("http - AuthRoutes - register - u.uc.Register")
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
// @Header       200 {string} Set-Cookie "sets cookie"
// @Failure      400  {object}  response
// @Failure      401  {object}  response
// @Failure      500  {object}  response
// @Router       /login [post]
func (u *AuthRoutes) Login(c *gin.Context) {
	var req loginRequest
	err := c.BindJSON(&req)
	if err != nil {
		logger.Error().Err(err).Msg("http - AuthRoutes - login - req validation")
		errorResponse(c, err, http.StatusBadRequest)
		return
	}
	token, err := u.uc.Login(c, req.Login, req.Password)
	if err != nil {
		logger.Error().Err(err).Msg("http - AuthRoutes - login - u.uc.Login")
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
