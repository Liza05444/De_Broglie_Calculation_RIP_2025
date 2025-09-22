package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"DeBroglieProject/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetParticlesAPI(ctx *gin.Context) {
	particleName := ctx.Query("particle")

	particles, err := h.Repository.GetParticles(particleName)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   particles,
	})
}

func (h *Handler) GetParticleAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	particle, err := h.Repository.GetParticle(uint(id))
	if err != nil {
		if err.Error() == "particle not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   particle,
	})
}

func (h *Handler) CreateParticleAPI(ctx *gin.Context) {
	type particleInput struct {
		Name        string  `json:"name" binding:"required"`
		Mass        float64 `json:"mass" binding:"required"`
		Description string  `json:"description"`
	}

	var input particleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	particle := ds.Particle{
		Name: input.Name,
		Mass: input.Mass,
		Description: sql.NullString{
			String: input.Description,
			Valid:  input.Description != "",
		},
	}

	createdParticle, err := h.Repository.CreateParticle(particle)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   createdParticle,
	})
}

func (h *Handler) UpdateParticleAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	type particleUpdateRequest struct {
		Name        string  `json:"name"`
		Mass        float64 `json:"mass"`
		Description string  `json:"description"`
	}

	var req particleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	particle := ds.Particle{
		Name: req.Name,
		Mass: req.Mass,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
	}

	err = h.Repository.UpdateParticle(uint(id), particle)
	if err != nil {
		if err.Error() == "particle not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Частица обновлена",
	})
}

func (h *Handler) DeleteParticleAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	particle, err := h.Repository.GetParticle(uint(id))
	if err != nil {
		if err.Error() == "particle not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// Удаляем изображение из MinIO, если оно существует
	if particle.Image.Valid {
		fileName := particle.Image.String
		parts := strings.Split(fileName, "/")
		if len(parts) > 0 {
			fileName = parts[len(parts)-1]
		}
		if fileName != "" {
			err = h.Repository.DeleteFileFromMinIO(ctx.Request.Context(), fileName)
			if err != nil {
				log.Printf("Warning: failed to delete image from MinIO: %v", err)
			}
		}
	}

	err = h.Repository.DeleteParticle(uint(id))
	if err != nil {
		if err.Error() == "particle not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Частица удалена",
	})
}

func (h *Handler) UploadParticleImageAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	_, err = h.Repository.GetParticle(uint(id))
	if err != nil {
		if err.Error() == "particle not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	src, err := file.Open()
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	defer src.Close()

	fileName := "particle_" + strconv.FormatUint(uint64(id), 10) + ".png"

	// Загружаем файл в MinIO
	err = h.Repository.UploadFileToMinIO(
		context.Background(),
		fileName,
		src,
		file.Size,
		file.Header.Get("Content-Type"),
	)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Сохраняем путь к файлу в базе данных
	imagePath := "http://127.0.0.1:9000/particles/" + fileName
	err = h.Repository.UpdateParticleImage(uint(id), imagePath)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Изображение загружено",
		"image_path": imagePath,
	})
}

func (h *Handler) AddParticleToRequestDeBroglieCalculationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	draft, _, err := h.Repository.GetDraftRequestDeBroglieCalculationInfo()
	if err != nil {
		created, createErr := h.Repository.CreateRequestDeBroglieCalculationWithParticle(uint(id))
		if createErr != nil {
			if strings.Contains(createErr.Error(), "duplicate key value violates unique constraint") {
				h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("частица уже добавлена в заявку"))
				return
			}
			h.errorHandler(ctx, http.StatusInternalServerError, createErr)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"message":  "Частица добавлена в новую заявку",
			"draft_id": created.ID,
		})
		return
	}

	if draft.Status != ds.RequestStatusDraft {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("заявка должна быть в статусе черновика"))
		return
	}

	if err := h.Repository.AddDeBroglieCalculationToRequest(draft.ID, uint(id)); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("частица уже добавлена в заявку"))
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Частица добавлена в заявку",
		"draft_id": draft.ID,
	})
}
