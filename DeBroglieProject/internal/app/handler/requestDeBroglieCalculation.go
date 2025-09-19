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
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Заявка на расчёт длины де Бройля для частиц не найдена",
		})
		return
	}

	if requestDeBroglieCalculation.Status != ds.RequestStatusDraft {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Заявка на расчёт длины де Бройля для частиц не найдена",
		})
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
