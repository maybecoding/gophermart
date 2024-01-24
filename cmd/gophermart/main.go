package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gophermart/internal/config"
	"gophermart/internal/controller/http"
	"gophermart/internal/usecase"
	"gophermart/internal/usecase/impl"
	"gophermart/internal/usecase/repo"
	"gophermart/pkg/logger"
	"gophermart/pkg/postgres"
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
		logger.Fatal().Err(fmt.Errorf("main - postgres.NewJwt: %w", err)).Msg("error")
	}
	defer pg.Close()

	// Инициализируем код приложения
	ucAuth := usecase.NewAuth(repo.NewAuth(pg), impl.NewPwd(), impl.NewJwt(cfg.JWT))
	uc := usecase.New(ucAuth)

	r := gin.Default()
	http.NewRouter(r, uc)

}
