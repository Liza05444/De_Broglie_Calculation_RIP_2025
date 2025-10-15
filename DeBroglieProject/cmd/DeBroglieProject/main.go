package main

import (
	"DeBroglieProject/internal/app/config"
	"DeBroglieProject/internal/app/dsn"
	"DeBroglieProject/internal/app/handler"
	"DeBroglieProject/internal/app/redis"
	"DeBroglieProject/internal/app/repository"
	"DeBroglieProject/internal/pkg"
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title De Broglie Project API
// @version 1.0
// @description API для расчета длины волны де Бройля частиц

// @host 127.0.0.1:8080
// @schemes http https
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://127.0.0.1:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	router.MaxMultipartMemory = 8 << 20 // 8 MB

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()

	rep, errRep := repository.New(
		postgresString,
		conf.MinIO.Endpoint,
		conf.MinIO.AccessKeyID,
		conf.MinIO.SecretAccessKey,
		conf.MinIO.BucketName,
		conf.MinIO.UseSSL,
	)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	ctx := context.Background()
	redisClient, err := redis.New(ctx, conf.Redis)
	if err != nil {
		logrus.Fatalf("error initializing redis: %v", err)
	}

	hand := handler.NewHandler(rep, conf, func(requireModerator bool) gin.HandlerFunc {
		tempApp := pkg.NewApp(conf, router, nil, redisClient)
		return tempApp.WithAuthCheck(requireModerator)
	}, redisClient)

	application := pkg.NewApp(conf, router, hand, redisClient)
	application.RunApp()
}
