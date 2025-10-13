package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"DeBroglieProject/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetRequestDeBroglieCalculationsAPI godoc
// @Summary Получение списка заявок
// @Description Возвращает список заявок на расчет с возможностью фильтрации
// @Tags RequestDeBroglieCalculations
// @Produce json
// @Security BearerAuth
// @Param status query string false "Статус заявки"
// @Param start_date query string false "Начальная дата (YYYY-MM-DD)"
// @Param end_date query string false "Конечная дата (YYYY-MM-DD)"
// @Success 200 {array} ds.RequestDeBroglieCalculation "Список заявок"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /requestdebrogliecalculations [get]
func (h *Handler) GetRequestDeBroglieCalculationsAPI(ctx *gin.Context) {
	var status *ds.RequestStatus
	var startDate, endDate *time.Time

	if statusStr := ctx.Query("status"); statusStr != "" {
		requestStatus := ds.RequestStatus(statusStr)
		status = &requestStatus
	}

	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	creatorID, ok := userUUID.(uuid.UUID)
	if !ok {
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("invalid user UUID in context"))
		return
	}

	isModerator, exists := ctx.Get("is_moderator")
	if !exists {
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("moderator status not found in context"))
		return
	}

	moderatorStatus, ok := isModerator.(bool)
	if !ok {
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("invalid moderator status in context"))
		return
	}

	requests, err := h.Repository.GetRequestDeBroglieCalculations(status, startDate, endDate, creatorID, moderatorStatus)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, requests)
}

// GetRequestDeBroglieCalculationAPI godoc
// @Summary Получение заявки по ID
// @Description Возвращает детальную информацию о заявке с расчетами
// @Tags RequestDeBroglieCalculations
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200 {object} deBroglieRequestDetailResponse "Детальная информация о заявке"
// @Failure 400 {object} errorResponse "Неверный ID"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 404 {object} errorResponse "Заявка не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /requestdebrogliecalculations/{id} [get]
func (h *Handler) GetRequestDeBroglieCalculationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, calcs, err := h.Repository.GetRequestWithCalculations(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	var deBroglieCalcs []gin.H
	for _, calc := range calcs {
		deBroglieCalcs = append(deBroglieCalcs, gin.H{
			"id":                calc.ID,
			"particle_id":       calc.ParticleID,
			"particle_name":     calc.Particle.Name,
			"particle_mass":     calc.Particle.Mass,
			"particle_image":    calc.Particle.Image,
			"speed":             calc.Speed,
			"de_broglie_length": calc.DeBroglieLength,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":                    request.ID,
		"name":                  request.Name,
		"status":                request.Status,
		"created_at":            request.CreatedAt,
		"formed_at":             request.FormedAt,
		"completed_at":          request.CompletedAt,
		"debrogliecalculations": deBroglieCalcs,
	})
}

// UpdateRequestDeBroglieCalculationAPI godoc
// @Summary Обновление заявки
// @Description Обновляет информацию о заявке
// @Tags RequestDeBroglieCalculations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param request body ds.RequestDeBroglieCalculation true "Обновленные данные заявки"
// @Success 200 {object} successResponse "Успешное обновление"
// @Failure 400 {object} errorResponse "Неверный формат запроса"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 404 {object} errorResponse "Заявка не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /requestdebrogliecalculations/{id} [put]
func (h *Handler) UpdateRequestDeBroglieCalculationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var request ds.RequestDeBroglieCalculation
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateRequestDeBroglieCalculation(uint(id), request)
	if err != nil {
		if err.Error() == "record not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "заявка обновлена",
	})
}

// DeleteRequestDeBroglieCalculationAPI godoc
// @Summary Удаление заявки
// @Description Удаляет заявку на расчет
// @Tags RequestDeBroglieCalculations
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200 {object} successResponse "Успешное удаление"
// @Failure 400 {object} errorResponse "Неверный ID"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 404 {object} errorResponse "Заявка не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /requestdebrogliecalculations/{id} [delete]
func (h *Handler) DeleteRequestDeBroglieCalculationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	rowsAffected, err := h.Repository.DeleteRequestDeBroglieCalculation(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	if rowsAffected == 0 {
		h.errorHandler(ctx, http.StatusNotFound, fmt.Errorf("заявка не найдена"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "заявка удалена",
	})
}

func (h *Handler) UpdateRequestStatusAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req struct {
		Status      ds.RequestStatus `json:"status" binding:"required"`
		ModeratorID *uuid.UUID       `json:"moderator_id,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateDeBroglieRequestStatus(uint(id), req.Status, req.ModeratorID)
	if err != nil {
		if err.Error() == "record not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusBadRequest, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "статус заявки обновлен",
	})
}

// DraftRequestDeBroglieCalculationInfoAPI godoc
// @Summary Информация о черновике заявки
// @Description Возвращает информацию о черновике заявки пользователя
// @Tags RequestDeBroglieCalculations
// @Produce json
// @Success 200 {object} deBroglieDraftInfoResponse "Информация о черновике"
// @Router /requestdebrogliecalculations/debrogliecart [get]
func (h *Handler) DraftRequestDeBroglieCalculationInfoAPI(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{
			"draft_id":      0,
			"particles_cnt": 0,
		})
		return
	}

	creatorID, ok := userUUID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"draft_id":      0,
			"particles_cnt": 0,
		})
		return
	}

	draft, calcs, err := h.Repository.GetDraftRequestDeBroglieCalculationInfo(creatorID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"draft_id":      0,
			"particles_cnt": 0,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"draft_id":      draft.ID,
		"particles_cnt": len(calcs),
	})
}

// FormRequestDeBroglieCalculationAPI godoc
// @Summary Формирование заявки
// @Description Переводит черновик заявки в статус "сформирована"
// @Tags RequestDeBroglieCalculations
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200 {object} successResponse "Успешное формирование"
// @Failure 400 {object} errorResponse "Неверный ID или статус заявки"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /requestdebrogliecalculations/{id}/form [put]
func (h *Handler) FormRequestDeBroglieCalculationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	creatorID, ok := userUUID.(uuid.UUID)
	if !ok {
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("invalid user UUID in context"))
		return
	}

	if err := h.Repository.FormDeBroglieRequestDraft(uint(id), creatorID); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "заявка сформирована",
	})
}

// CompleteRequestDeBroglieCalculationAPI godoc
// @Summary Завершение заявки модератором
// @Description Одобряет или отклоняет заявку (только для модераторов)
// @Tags RequestDeBroglieCalculations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param request body object{approve=boolean} true "Решение модератора"
// @Success 200 {object} successResponse "Успешное завершение"
// @Failure 400 {object} errorResponse "Неверный формат запроса"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 403 {object} errorResponse "Недостаточно прав"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /requestdebrogliecalculations/{id}/complete [put]
func (h *Handler) CompleteRequestDeBroglieCalculationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req struct {
		Approve bool `json:"approve"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	moderatorID, ok := userUUID.(uuid.UUID)
	if !ok {
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("invalid user UUID in context"))
		return
	}

	if err := h.Repository.CompleteDeBroglieRequest(uint(id), req.Approve, moderatorID); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	message := "заявка отклонена"
	if req.Approve {
		message = "заявка одобрена и обработана"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
