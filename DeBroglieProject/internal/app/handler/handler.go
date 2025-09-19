package handler

import (
	"DeBroglieProject/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/particles", h.GetParticles)
	router.GET("/particle/:id", h.GetParticle)
	router.POST("/de-broglie-calculation/:id/add-particle", h.AddParticleToRequest)
	router.GET("/de-broglie-calculation/:id", h.GetRequestDeBroglieCalculation)
	router.POST("/de-broglie-calculation/:id/delete-request", h.DeleteRequestDeBroglieCalculation)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static/styles", "./resources/styles")
	router.Static("/static/img", "./resources/img")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
