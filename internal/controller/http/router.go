package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "gophermart/docs"
	"gophermart/internal/usecase"
)

// @title           Gopher Mart
// @version         1.0
// @description     Market for gophers

// @contact.name   Dmitry Petukhov
// @contact.email  maybecoding@gmail.com

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// NewRouter - Создает новый роутер
func NewRouter(r *gin.Engine, uc *usecase.UseCase) {

	api := r.Group("/api")
	{
		userRoutes(api, uc.Auth)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
