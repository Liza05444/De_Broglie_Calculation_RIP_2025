package handler

import (
	"net/http"
	"strconv"

	"DeBroglieProject/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetRequestDeBroglieCalculation(ctx *gin.Context) {
	idRequsetDeBroglieCalculation := ctx.Param("id")
	_, err := strconv.Atoi(idRequsetDeBroglieCalculation)
	if err != nil {
		logrus.Error(err)
	}

	requestDeBroglieCalculation, deBroglieCalculations, err := h.Repository.GetDraftRequestDeBroglieCalculationInfo()
	if err != nil {
		logrus.Error(err)
	}

	if requestDeBroglieCalculation.Status != ds.RequestStatusDraft {
		logrus.Warnf("Attempt to access non-draft request (ID: %d, Status: %s)", requestDeBroglieCalculation.ID, requestDeBroglieCalculation.Status)
		ctx.Redirect(http.StatusFound, "/particles")
		return
	}

	var particles []ds.Particle
	particles, err = h.Repository.GetParticles()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "request_de_broglie_calculation.html", gin.H{
		"requestDeBroglieCalculation": requestDeBroglieCalculation,
		"calculations":                deBroglieCalculations,
		"calculationsCount":           len(deBroglieCalculations),
		"particles":                   particles,
	})
}

func (h *Handler) AddParticleToRequest(ctx *gin.Context) {
	particleIDStr := ctx.Param("id")
	particleID, err := strconv.Atoi(particleIDStr)
	if err != nil {
		logrus.Error("Error converting particle ID:", err)
	}

	draftRequest, _, err := h.Repository.GetDraftRequestDeBroglieCalculationInfo()
	if err != nil {
		_, err := h.Repository.CreateRequestDeBroglieCalculation(uint(particleID))
		if err != nil {
			logrus.Error("Error creating new draft request:", err)
		}
	} else {
		err = h.Repository.AddDeBroglieCalculationToRequest(draftRequest.ID, uint(particleID))
		if err != nil {
			logrus.Error("Error adding particle to existing request:", err)
		}
	}

	ctx.Redirect(http.StatusFound, "/particles")
}

func (h *Handler) DeleteRequestDeBroglieCalculation(ctx *gin.Context) {
	idRequestDeBroglieCalculation := ctx.Param("id")
	id, err := strconv.Atoi(idRequestDeBroglieCalculation)
	if err != nil {
		logrus.Errorf("Error converting request ID: %v", err)
	}

	err = h.Repository.DeleteRequestDeBroglieCalculation(id)
	if err != nil {
		logrus.Errorf("Error deleting request: %v", err)
	}

	ctx.Redirect(http.StatusFound, "/particles")
}
