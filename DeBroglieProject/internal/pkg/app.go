package pkg

import (
	"fmt"

	"DeBroglieProject/internal/app/config"
	"DeBroglieProject/internal/app/handler"
	"DeBroglieProject/internal/app/redis"
	_ "DeBroglieProject/docs" 

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
	Redis   *redis.Client
}

func NewApp(c *config.Config, r *gin.Engine, h *handler.Handler, redisClient *redis.Client) *Application {
	return &Application{
		Config:  c,
		Router:  r,
		Handler: h,
		Redis:   redisClient,
	}
}

func (a *Application) RunApp() {
	logrus.Info("Server start up")

	a.Handler.RegisterAPI(a.Router)
	
	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	if err := a.Router.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server down")
}
