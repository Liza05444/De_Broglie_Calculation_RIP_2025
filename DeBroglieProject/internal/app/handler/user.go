package handler

import (
	"errors"
	"net/http"

	"DeBroglieProject/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetMeAPI godoc
// @Summary Получение информации о текущем пользователе
// @Description Возвращает информацию о авторизованном пользователе
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ds.User "Информация о пользователе"
// @Failure 401 {object} errorResponse "Пользователь не авторизован"
// @Failure 404 {object} errorResponse "Пользователь не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /profile/me [get]
func (h *Handler) GetMeAPI(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	meID, ok := userUUID.(uuid.UUID)
	if !ok {
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("invalid user UUID in context"))
		return
	}

	user, err := h.Repository.GetUserByUUID(meID)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// UpdateMeAPI godoc
// @Summary Обновление информации о текущем пользователе
// @Description Обновляет информацию о авторизованном пользователе
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ds.User true "Данные для обновления"
// @Success 200 {object} ds.User "Обновленная информация о пользователе"
// @Failure 400 {object} errorResponse "Неверный формат запроса"
// @Failure 401 {object} errorResponse "Пользователь не авторизован"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /profile/me [put]
func (h *Handler) UpdateMeAPI(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	meID, ok := userUUID.(uuid.UUID)
	if !ok {
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("invalid user UUID in context"))
		return
	}

	var user ds.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.Repository.UpdateUserByUUID(meID, user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updatedUser, err := h.Repository.GetUserByUUID(meID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("failed to retrieve updated user data"))
		return
	}
	ctx.JSON(http.StatusOK, updatedUser)
}
