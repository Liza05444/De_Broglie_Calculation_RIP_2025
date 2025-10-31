package handler

import (
	"DeBroglieProject/internal/app/ds"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateCalculationSpeedAPI godoc
// @Summary Обновление скорости частицы
// @Description Обновляет скорость частицы в расчете де Бройля
// @Tags DeBroglieCalculations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param particleId path int true "ID частицы"
// @Param request body object{speed=number} true "Новая скорость"
// @Success 200 {object} successResponse "Успешное обновление"
// @Failure 400 {object} errorResponse "Неверные параметры или статус заявки"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 404 {object} errorResponse "Заявка не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /debrogliecalculations/{id}/particle/{particleId} [put]
func (h *Handler) UpdateCalculationSpeedAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	particleIdStr := ctx.Param("particleId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	particleID, err := strconv.ParseUint(particleIdStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetRequestDeBroglieCalculation(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if request.Status != ds.RequestStatusDraft {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("заявка должна быть в статусе черновика"))
		return
	}

	var body struct {
		Speed float64 `json:"speed" binding:"required,gt=0"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("скорость не может быть отрицательной"))
		return
	}

	if err := h.Repository.UpdateCalculationValue(uint(id), uint(particleID), body.Speed); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "скорость обновлена",
	})
}

// RemoveParticleFromRequestDeBroglieCalculationAPI godoc
// @Summary Удаление частицы из заявки
// @Description Удаляет частицу из черновика заявки
// @Tags DeBroglieCalculations
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param particleId path int true "ID частицы"
// @Success 200 {object} successResponse "Успешное удаление"
// @Failure 400 {object} errorResponse "Неверные параметры или статус заявки"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 404 {object} errorResponse "Заявка или частица не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /debrogliecalculations/{id}/particle/{particleId} [delete]
func (h *Handler) RemoveParticleFromRequestDeBroglieCalculationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	particleIdStr := ctx.Param("particleId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	particleID, err := strconv.ParseUint(particleIdStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetRequestDeBroglieCalculation(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if request.Status != ds.RequestStatusDraft {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("заявка должна быть в статусе черновика"))
		return
	}

	rowsAffected, err := h.Repository.RemoveCalculationFromRequest(uint(id), uint(particleID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		h.errorHandler(ctx, http.StatusNotFound, fmt.Errorf("частица не найдена в заявке"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "частица удалена из заявки",
	})
}

// UpdateCalculationDeBroglieLengthAPI godoc
// @Summary Обновление длины волны де Бройля
// @Description Обновляет длину волны де Бройля для расчета (вызывается асинхронным сервисом)
// @Tags DeBroglieCalculations
// @Accept json
// @Produce json
// @Param id path int true "ID расчета"
// @Param request body object{de_broglie_length=number,token=string} true "Длина волны и токен авторизации"
// @Success 200 {object} successResponse "Успешное обновление"
// @Failure 400 {object} errorResponse "Неверные параметры или неверный токен"
// @Failure 404 {object} errorResponse "Расчет не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /debrogliecalculations/{id}/update-length [put]
func (h *Handler) UpdateCalculationDeBroglieLengthAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req struct {
		DeBroglieLength *float64 `json:"de_broglie_length"`
		Rejected        *bool    `json:"rejected"`
		Token           string   `json:"token"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	const AUTH_TOKEN = "12345678"
	if req.Token != AUTH_TOKEN {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("invalid token"))
		return
	}

	calc, err := h.Repository.GetCalculationByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	requestID := calc.RequestDeBroglieCalculationID

	if req.Rejected != nil && *req.Rejected {
		request, err := h.Repository.GetRequestDeBroglieCalculation(requestID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		if request.Status == ds.RequestStatusFormed {
			err = h.Repository.UpdateDeBroglieRequestStatus(requestID, ds.RequestStatusRejected, nil)
			if err != nil {
				h.errorHandler(ctx, http.StatusInternalServerError, err)
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "заявка отклонена",
		})
		return
	}

	if req.DeBroglieLength == nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("de_broglie_length is required"))
		return
	}

	err = h.Repository.UpdateCalculationDeBroglieLength(uint(id), *req.DeBroglieLength)
	if err != nil {
		if err.Error() == "record not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	calculatedCount, err := h.Repository.CountCalculationsWithDeBroglieLength(requestID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	totalCount, err := h.Repository.CountTotalCalculationsForRequest(requestID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if calculatedCount > 0 && calculatedCount == totalCount {
		request, err := h.Repository.GetRequestDeBroglieCalculation(requestID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		if request.Status == ds.RequestStatusFormed {
			err = h.Repository.UpdateDeBroglieRequestStatus(requestID, ds.RequestStatusCompleted, nil)
			if err != nil {
				h.errorHandler(ctx, http.StatusInternalServerError, err)
				return
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "заявка обработана",
	})
}
