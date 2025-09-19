package handler

import (
	"net/http"
	"strconv"

	"DeBroglieProject/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetParticles(ctx *gin.Context) {
	var particles []ds.Particle
	var err error

	searchParticle := ctx.Query("particle")
	if searchParticle == "" {
		particles, err = h.Repository.GetParticles()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		particles, err = h.Repository.GetParticlesByName(searchParticle)
		if err != nil {
			logrus.Error(err)
		}
	}

	draftRequest, deBroglieCalculations, err := h.Repository.GetDraftRequestDeBroglieCalculationInfo()
	var draftRequestID uint = 0
	var deBroglieCalculationsCount int64 = 0
	if err == nil {
		draftRequestID = draftRequest.ID
		deBroglieCalculationsCount = int64(len(deBroglieCalculations))
	}

	ctx.HTML(http.StatusOK, "particles.html", gin.H{
		"particles":                  particles,
		"particle":                   searchParticle,
		"deBroglieCalculationsCount": deBroglieCalculationsCount,
		"draftRequestID":             draftRequestID,
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

func (h *Handler) AddParticleToRequest(ctx *gin.Context) {
	particleIDStr := ctx.Param("id")
	particleID, err := strconv.Atoi(particleIDStr)
	if err != nil {
		logrus.Error("Error converting particle ID:", err)
	}

	draftRequestDeBroglieCalculation, _, err := h.Repository.GetDraftRequestDeBroglieCalculationInfo()
	if err != nil {
		_, err := h.Repository.CreateRequestDeBroglieCalculation(uint(particleID))
		if err != nil {
			logrus.Error("Error creating new draft request:", err)
		}
	} else {
		err = h.Repository.AddDeBroglieCalculationToRequest(draftRequestDeBroglieCalculation.ID, uint(particleID))
		if err != nil {
			logrus.Error("Error adding particle to existing request:", err)
		}
	}

	ctx.Redirect(http.StatusFound, "/particles")
}
