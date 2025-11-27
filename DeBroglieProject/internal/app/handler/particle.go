package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"DeBroglieProject/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetParticlesAPI godoc
// @Summary Получение списка частиц
// @Description Возвращает список всех частиц с возможностью фильтрации по названию
// @Tags Particles
// @Produce json
// @Param particle query string false "Название частицы для фильтрации"
// @Success 200 {array} ds.Particle "Список частиц"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /particles [get]
func (h *Handler) GetParticlesAPI(ctx *gin.Context) {
	particleName := ctx.Query("particle")

	particles, err := h.Repository.GetParticles(particleName)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, particles)
}

// GetParticleAPI godoc
// @Summary Получение частицы по ID
// @Description Возвращает информацию о частице по её ID
// @Tags Particles
// @Produce json
// @Param id path int true "ID частицы"
// @Success 200 {object} ds.Particle "Информация о частице"
// @Failure 400 {object} errorResponse "Неверный ID"
// @Failure 404 {object} errorResponse "Частица не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /particles/{id} [get]
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

	ctx.JSON(http.StatusOK, particle)
}

// CreateParticleAPI godoc
// @Summary Создание новой частицы
// @Description Создает новую частицу (только для модераторов)
// @Tags Particles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{name=string,mass=number,description=string} true "Данные частицы"
// @Success 201 {object} ds.Particle "Созданная частица"
// @Failure 400 {object} errorResponse "Неверный формат запроса"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 403 {object} errorResponse "Недостаточно прав"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /particles [post]
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
	}
	if input.Description != "" {
		particle.Description = &input.Description
	}

	createdParticle, err := h.Repository.CreateParticle(particle)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, createdParticle)
}

// UpdateParticleAPI godoc
// @Summary Обновление частицы
// @Description Обновляет информацию о частице (только для модераторов)
// @Tags Particles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID частицы"
// @Param request body object{name=string,mass=number,description=string} true "Обновленные данные частицы"
// @Success 200 {object} successResponse "Успешное обновление"
// @Failure 400 {object} errorResponse "Неверный формат запроса"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 403 {object} errorResponse "Недостаточно прав"
// @Failure 404 {object} errorResponse "Частица не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /particles/{id} [put]
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
	}
	if req.Description != "" {
		particle.Description = &req.Description
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
		"message": "частица обновлена",
	})
}

// DeleteParticleAPI godoc
// @Summary Удаление частицы
// @Description Удаляет частицу и связанное с ней изображение (только для модераторов)
// @Tags Particles
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID частицы"
// @Success 200 {object} successResponse "Успешное удаление"
// @Failure 400 {object} errorResponse "Неверный ID"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 403 {object} errorResponse "Недостаточно прав"
// @Failure 404 {object} errorResponse "Частица не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /particles/{id} [delete]
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
	if particle.Image != nil {
		fileName := *particle.Image
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
		"message": "частица удалена",
	})
}

// UploadParticleImageAPI godoc
// @Summary Загрузка изображения частицы
// @Description Загружает изображение для частицы в MinIO (только для модераторов)
// @Tags Particles
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID частицы"
// @Param image formData file true "Файл изображения"
// @Success 200 {object} successResponse "Успешная загрузка"
// @Failure 400 {object} errorResponse "Неверный формат запроса"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 403 {object} errorResponse "Недостаточно прав"
// @Failure 404 {object} errorResponse "Частица не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /particles/{id}/image [post]
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
		"message": "изображение загружено",
	})
}

// AddParticleToRequestDeBroglieCalculationAPI godoc
// @Summary Добавление частицы в заявку
// @Description Добавляет частицу в черновик заявки на расчет де Бройля
// @Tags Particles
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID частицы"
// @Success 200 {object} successResponse "Успешное добавление"
// @Failure 400 {object} errorResponse "Неверный ID или частица уже добавлена"
// @Failure 401 {object} errorResponse "Требуется авторизация"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /particles/{id}/addParticle [post]
func (h *Handler) AddParticleToRequestDeBroglieCalculationAPI(ctx *gin.Context) {
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

	researcherID, ok := userUUID.(uuid.UUID)
	if !ok {
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("invalid user UUID in context"))
		return
	}

	draft, _, err := h.Repository.GetDraftRequestDeBroglieCalculationInfo(researcherID)
	if err != nil {
		_, createErr := h.Repository.CreateRequestDeBroglieCalculationWithParticle(uint(id), researcherID)
		if createErr != nil {
			if strings.Contains(createErr.Error(), "duplicate key value violates unique constraint") {
				h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("частица уже добавлена в заявку"))
				return
			}
			h.errorHandler(ctx, http.StatusInternalServerError, createErr)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "частица добавлена в новую заявку",
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
		"message": "частица добавлена в заявку",
	})
}
