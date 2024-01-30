package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gophermart/internal/config"
	"gophermart/internal/usecase"
	"net/http"
)

type Http struct {
	server *http.Server
}

func New() *Http {
	return &Http{}
}

func (h *Http) Run(uc *usecase.UseCase, cfg config.HTTP) error {
	r := newRouter(gin.Default(), uc)
	h.server = &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	return fmt.Errorf("http - Run - http.ListenAndServe: %w", h.server.ListenAndServe())
}

func (h *Http) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return h.server.Shutdown(ctx)
}
