package handler

import (
	"DeBroglieProject/internal/app/repository"
	"net/http"
	"strconv"

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

func (h *Handler) GetParticles(ctx *gin.Context) {
	var particles []repository.Particle
	var err error

	searchParticle := ctx.Query("particle")
	if searchParticle == "" {
		particles, err = h.Repository.GetParticles()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		particles, err = h.Repository.GetParticleByName(searchParticle)
		if err != nil {
			logrus.Error(err)
		}
	}

	requestID := 1
	_, deBroglieCalculations, err := h.Repository.GetDeBroglieCalculationsForRequest(requestID)
	if err != nil {
		logrus.Error(err)
	}
	deBroglieCalculationsCount := len(deBroglieCalculations)

	ctx.HTML(http.StatusOK, "particles.html", gin.H{
		"particles":                  particles,
		"particle":                   searchParticle,
		"deBroglieCalculationsCount": deBroglieCalculationsCount,
	})
}

func (h *Handler) GetParticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	particle, err := h.Repository.GetParticle(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "particle.html", gin.H{
		"particle": particle,
	})
}

func (h *Handler) GetDeBroglieCalculation(ctx *gin.Context) {
	idRequset := ctx.Param("id")
	id, err := strconv.Atoi(idRequset)
	if err != nil {
		logrus.Error(err)
	}

	requestDeBroglieCalculation, deBroglieCalculations, err := h.Repository.GetDeBroglieCalculationsForRequest(id)
	if err != nil {
		logrus.Error(err)
	}

	var particles []repository.Particle
	particles, err = h.Repository.GetParticles()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "request_de_broglie_calculation.html", gin.H{
		"requestDeBroglieCalculation": requestDeBroglieCalculation,
		"calculations": deBroglieCalculations,
		"particles": particles,
	})
}
