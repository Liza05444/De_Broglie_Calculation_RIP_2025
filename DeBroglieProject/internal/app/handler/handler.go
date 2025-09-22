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

func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		profile := api.Group("/profile")
		{
			profile.POST("/register", h.RegisterUserAPI)
			profile.POST("/login", h.LoginAPI)
			profile.POST("/logout", h.LogoutAPI)
			profile.GET("/me", h.GetMeAPI)
			profile.PUT("/me", h.UpdateMeAPI)
		}

		particles := api.Group("/particles")
		{
			particles.GET("", h.GetParticlesAPI)
			particles.GET("/:id", h.GetParticleAPI)
			particles.POST("", h.CreateParticleAPI)
			particles.PUT("/:id", h.UpdateParticleAPI)
			particles.DELETE("/:id", h.DeleteParticleAPI)
			particles.POST("/:id/image", h.UploadParticleImageAPI)
			particles.POST("/:id/addParticle", h.AddParticleToRequestDeBroglieCalculationAPI)
		}

		requests := api.Group("/requestdebrogliecalculations")
		{
			requests.GET("/cart", h.DraftRequestDeBroglieCalculationInfoAPI)
			requests.GET("", h.GetRequestDeBroglieCalculationsAPI)
			requests.GET("/:id", h.GetRequestDeBroglieCalculationAPI)
			requests.PUT("/:id", h.UpdateRequestDeBroglieCalculationAPI)
			requests.DELETE("/:id", h.DeleteRequestDeBroglieCalculationAPI)
			requests.PUT("/:id/form", h.FormRequestDeBroglieCalculationAPI)
			requests.PUT("/:id/complete", h.CompleteRequestDeBroglieCalculationAPI)
		}

		calculations := api.Group("/debrogliecalculations")
		{
			calculations.PUT("/:id/particle/:particleId", h.UpdateCalculationSpeedAPI)
			calculations.DELETE("/:id/particle/:particleId", h.RemoveParticleFromRequestDeBroglieCalculationAPI)
		}
	}
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
