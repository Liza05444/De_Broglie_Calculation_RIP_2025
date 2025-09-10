package api

import (
	"DeBroglieProject/internal/app/handler"
	"DeBroglieProject/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/particles", handler.GetParticles)
	r.GET("/particle/:id", handler.GetParticle)
	r.GET("/de-broglie-calculation/:id", handler.GetDeBroglieCalculation)

	r.Run()
	log.Println("Server down")
}
