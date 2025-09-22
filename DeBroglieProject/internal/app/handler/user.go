package handler

import (
	"net/http"

	"DeBroglieProject/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUserAPI(ctx *gin.Context) {
	var user ds.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	createdUser, err := h.Repository.CreateUser(user)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   createdUser,
	})
}

func (h *Handler) LoginAPI(ctx *gin.Context) {
	var user ds.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	authenticatedUser, err := h.Repository.CheckCredentials(user.Email, user.Password)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   authenticatedUser,
	})
}

func (h *Handler) LogoutAPI(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "logged out",
	})
}

func (h *Handler) GetMeAPI(ctx *gin.Context) {
	meID := ds.GetCreatorID()
	user, err := h.Repository.GetUserByID(meID)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

func (h *Handler) UpdateMeAPI(ctx *gin.Context) {
	meID := ds.GetCreatorID()
	var user ds.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.Repository.UpdateUser(meID, user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	updatedUser, _ := h.Repository.GetUserByID(meID)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   updatedUser,
	})
}
