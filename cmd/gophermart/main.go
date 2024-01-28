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

	//f, err := os.ReadFile("/users/d.petukhov/papa/img.jpg")
	//if err != nil {
	//	panic(err)
	//}
	//
	//bs64 := base64.StdEncoding.EncodeToString(f)
	//fmt.Println(bs64)

	// Инициализируем код приложения
	ucAuth := usecase.NewAuth(repo.NewAuth(pg), impl.NewPwd(), impl.NewJwt(cfg.JWT))
	ucOrder := usecase.NewOrder(repo.NewOrder(pg), impl.NewOrderNumAlgImpl())
	uc := usecase.New(ucAuth, ucOrder)

	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})
	http.NewRouter(r, uc)

}
