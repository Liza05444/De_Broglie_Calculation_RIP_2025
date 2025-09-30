package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"DeBroglieProject/internal/app/ds"

	"github.com/gin-gonic/gin"
)

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

	requests, err := h.Repository.GetRequestDeBroglieCalculations(status, startDate, endDate)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	var simplifiedRequests []gin.H
	for _, req := range requests {
		simplifiedRequests = append(simplifiedRequests, gin.H{
			"id":           req.ID,
			"name":         req.Name,
			"status":       req.Status,
			"created_at":   req.CreatedAt,
			"formed_at":    req.FormedAt,
			"completed_at": req.CompletedAt,
		})
	}

	ctx.JSON(http.StatusOK, simplifiedRequests)
}

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

	var simplifiedCalcs []gin.H
	for _, calc := range calcs {
		simplifiedCalcs = append(simplifiedCalcs, gin.H{
			"id":                  calc.ID,
			"particle_id":         calc.ParticleID,
			"particle_name":       calc.Particle.Name,
			"particle_mass":       calc.Particle.Mass,
			"particle_image":      calc.Particle.Image,
			"speed":               calc.Speed,
			"de_broglie_length":   calc.DeBroglieLength,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           request.ID,
		"name":         request.Name,
		"status":       request.Status,
		"created_at":   request.CreatedAt,
		"formed_at":    request.FormedAt,
		"completed_at": request.CompletedAt,
		"calculations": simplifiedCalcs,
	})
}

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
		"status":  "success",
		"message": "Заявка обновлена",
	})
}

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
		"status":  "success",
		"message": "Заявка удалена",
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
		ModeratorID *uint            `json:"moderator_id,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateRequestStatus(uint(id), req.Status, req.ModeratorID)
	if err != nil {
		if err.Error() == "record not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusBadRequest, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Статус заявки обновлен",
	})
}

func (h *Handler) DraftRequestDeBroglieCalculationInfoAPI(ctx *gin.Context) {
	draft, calcs, err := h.Repository.GetDraftRequestDeBroglieCalculationInfo()
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

func (h *Handler) FormRequestDeBroglieCalculationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.Repository.FormDraft(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Заявка сформирована",
	})
}

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

	if err := h.Repository.CompleteRequest(uint(id), req.Approve); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	message := "Заявка отклонена"
	if req.Approve {
		message = "Заявка одобрена и обработана"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
	})
}
