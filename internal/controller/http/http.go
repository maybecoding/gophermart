package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gophermart/internal/config"
	"gophermart/internal/usecase"
	"net/http"
)

type HTTP struct {
	server *http.Server
}

func New() *HTTP {
	return &HTTP{}
}

func (h *HTTP) Run(uc *usecase.UseCase, cfg config.HTTP) error {
	r := newRouter(gin.Default(), uc)
	h.server = &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	return fmt.Errorf("http - Run - http.ListenAndServe: %w", h.server.ListenAndServe())
}

func (h *HTTP) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return h.server.Shutdown(ctx)
}
