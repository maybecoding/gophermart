package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gophermart/internal/config"
	"gophermart/internal/controller/http"
	"gophermart/internal/entity"
	"gophermart/internal/migration"
	"gophermart/internal/usecase"
	"gophermart/internal/usecase/accrual"
	"gophermart/internal/usecase/impl"
	"gophermart/internal/usecase/repo"
	"gophermart/pkg/logger"
	"gophermart/pkg/postgres"
	nethttp "net/http"
	"os"
	"os/signal"
)

func main() {
	// Читаем конфиг
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	// Инициализируем логгер
	logger.Init(cfg.Log.Level)

	// Подключаемся к БД
	pg, err := postgres.New(cfg.PG.URI)
	if err != nil {
		logger.Fatal().Err(err).Msg("main - postgres.NewJwt")
	}
	defer pg.Close()

	// Выполняем миграции
	err = migration.Run(cfg.PG.URI)
	if err != nil {
		logger.Fatal().Err(err).Msg("main - migration.Run")
	}

	// Инициализируем код приложения
	ucAuth := usecase.NewAuth(repo.NewAuth(pg), impl.NewPwd(), impl.NewJwt(cfg.JWT))
	repoBonus := repo.NewBonus(pg)
	ucOrder := usecase.NewOrder(repo.NewOrder(pg), repoBonus, impl.NewOrderNumAlgImpl(), accrual.NewOrder(cfg.AccrualSystem))
	ucBonus := usecase.NewBonus(repoBonus)
	uc := usecase.New(ucAuth, ucOrder, ucBonus)

	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})

	// Контекст, который будет отменен при выходе из приложения Ctrl + C
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	g, gCtx := errgroup.WithContext(ctx)

	// Процесс, обновляющий статусы
	g.Go(func() error {
		ucOrder.RunAccrualRefresh(gCtx)
		logger.Debug().Msg("accrual refresh stopped")
		return entity.ErrGracefulShutdown
	})
	h := http.New()

	// Запускаем сервер
	g.Go(func() error {
		err := h.Run(uc, cfg.HTTP)
		logger.Debug().Msg("server stopped")
		return err
	})
	// Запускаем выключатель для сервера
	g.Go(func() error {
		err := h.Shutdown(gCtx)
		logger.Debug().Msg("server stopper stopped")
		return err
	})

	if err = g.Wait(); err != nil && !errors.Is(err, nethttp.ErrServerClosed) && !errors.Is(err, entity.ErrGracefulShutdown) {
		logger.Error().Err(err).Msg("main - error")
	} else {
		logger.Info().Msg("app gracefully stopped")
	}
}
