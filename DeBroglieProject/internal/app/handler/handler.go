package handler

import (
	"DeBroglieProject/internal/app/config"
	"DeBroglieProject/internal/app/redis"
	"DeBroglieProject/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository        *repository.Repository
	Config            *config.Config
	AuthCheck         func(requireModerator bool) gin.HandlerFunc
	OptionalAuthCheck func() gin.HandlerFunc
	Redis             *redis.Client
}

func NewHandler(r *repository.Repository, cfg *config.Config, authCheck func(requireModerator bool) gin.HandlerFunc, optionalAuthCheck func() gin.HandlerFunc, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository:        r,
		Config:            cfg,
		AuthCheck:         authCheck,
		OptionalAuthCheck: optionalAuthCheck,
		Redis:             redisClient,
	}
}

func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		h.registerPublicEndpoints(api)
		h.registerProfileEndpoints(api)
		h.registerParticleEndpoints(api)
		h.registerDeBroglieRequestEndpoints(api)
	}
}

func (h *Handler) registerPublicEndpoints(api *gin.RouterGroup) {
	api.GET("/particles", h.GetParticlesAPI)
	api.GET("/particles/:id", h.GetParticleAPI)
	
	debroglieCart := api.Group("/requestdebrogliecalculations")
	debroglieCart.Use(h.OptionalAuthCheck())
	{
		debroglieCart.GET("/debrogliecart", h.DraftRequestDeBroglieCalculationInfoAPI)
	}
}

func (h *Handler) registerProfileEndpoints(api *gin.RouterGroup) {
	profile := api.Group("/profile")
	{
		profile.POST("/register", h.Register)
		profile.POST("/login", h.Login)
		profile.POST("/logout", h.Logout)

		profileProtected := profile.Group("")
		profileProtected.Use(h.AuthCheck(false))
		{
			profileProtected.GET("/me", h.GetMeAPI)
			profileProtected.PUT("/me", h.UpdateMeAPI)
		}
	}
}

func (h *Handler) registerParticleEndpoints(api *gin.RouterGroup) {
	particlesModerator := api.Group("/particles")
	particlesModerator.Use(h.AuthCheck(true))
	{
		particlesModerator.POST("", h.CreateParticleAPI)
		particlesModerator.PUT("/:id", h.UpdateParticleAPI)
		particlesModerator.DELETE("/:id", h.DeleteParticleAPI)
		particlesModerator.POST("/:id/image", h.UploadParticleImageAPI)
	}

	particlesUser := api.Group("/particles")
	particlesUser.Use(h.AuthCheck(false))
	{
		particlesUser.POST("/:id/addParticle", h.AddParticleToRequestDeBroglieCalculationAPI)
	}
}

func (h *Handler) registerDeBroglieRequestEndpoints(api *gin.RouterGroup) {
	deBroglieRequests := api.Group("/requestdebrogliecalculations")
	deBroglieRequests.Use(h.AuthCheck(false))
	{
		deBroglieRequests.GET("", h.GetRequestDeBroglieCalculationsAPI)
		deBroglieRequests.GET("/:id", h.GetRequestDeBroglieCalculationAPI)
		deBroglieRequests.PUT("/:id", h.UpdateRequestDeBroglieCalculationAPI)
		deBroglieRequests.DELETE("/:id", h.DeleteRequestDeBroglieCalculationAPI)
		deBroglieRequests.PUT("/:id/form", h.FormRequestDeBroglieCalculationAPI)
	}

	deBroglieRequestsModerator := api.Group("/requestdebrogliecalculations")
	deBroglieRequestsModerator.Use(h.AuthCheck(true))
	{
		deBroglieRequestsModerator.PUT("/:id/complete", h.CompleteRequestDeBroglieCalculationAPI)
	}

	deBroglieCalculations := api.Group("/debrogliecalculations")
	deBroglieCalculations.Use(h.AuthCheck(false))
	{
		deBroglieCalculations.PUT("/:id/particle/:particleId", h.UpdateCalculationSpeedAPI)
		deBroglieCalculations.DELETE("/:id/particle/:particleId", h.RemoveParticleFromRequestDeBroglieCalculationAPI)
	}
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
