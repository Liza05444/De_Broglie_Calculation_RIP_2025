package handler

import (
	"DeBroglieProject/internal/app/ds"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type errorResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

type successResponse struct {
	Message string `json:"message"`
}

type deBroglieDraftInfoResponse struct {
	DraftID      uint `json:"draft_id"`
	ParticlesCnt int  `json:"particles_cnt"`
}

type deBroglieCalculationDetail struct {
	ID              uint     `json:"id"`
	ParticleID      uint     `json:"particle_id"`
	ParticleName    string   `json:"particle_name"`
	ParticleMass    float64  `json:"particle_mass"`
	ParticleImage   *string  `json:"particle_image"`
	Speed           float64  `json:"speed"`
	DeBroglieLength *float64 `json:"de_broglie_length"`
}

type deBroglieRequestDetailResponse struct {
	ID           uint                         `json:"id"`
	Name         *string                      `json:"name"`
	Status       string                       `json:"status"`
	CreatedAt    string                       `json:"created_at"`
	FormedAt     *string                      `json:"formed_at"`
	CompletedAt  *string                      `json:"completed_at"`
	Calculations []deBroglieCalculationDetail `json:"calculations"`
}

// Login godoc
// @Summary Авторизация пользователя
// @Description Авторизация пользователя по email и паролю
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body loginReq true "Данные для входа"
// @Success 200 {object} loginResp "Успешная авторизация"
// @Failure 400 {object} errorResponse "Неверный формат запроса"
// @Failure 401 {object} errorResponse "Неверный email или пароль"
// @Failure 403 {object} errorResponse "Доступ запрещен"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /profile/login [post]
func (h *Handler) Login(gCtx *gin.Context) {
	cfg := h.Config
	req := &loginReq{}

	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":      "error",
				"description": "invalid email or password",
			})
		} else {
			gCtx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	if req.Email == user.Email && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) == nil {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &ds.JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.ExpiresIn)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "bitop-admin",
			},
			UserUUID:    user.ID,
			IsModerator: user.IsModerator,
		})

		if token == nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
			return
		}

		strToken, err := token.SignedString([]byte(cfg.JWT.Token))
		if err != nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
			return
		}

		gCtx.JSON(http.StatusOK, loginResp{
			ExpiresIn:   int64(cfg.JWT.ExpiresIn.Seconds()),
			AccessToken: strToken,
			TokenType:   "Bearer",
		})
		return
	}

	gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":      "error",
		"description": "invalid email or password",
	})
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создание нового пользователя в системе
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body registerReq true "Данные для регистрации"
// @Success 200 {object} ds.User "Успешная регистрация"
// @Failure 400 {object} errorResponse "Неверный формат запроса или пустые поля"
// @Failure 409 {object} errorResponse "Пользователь с таким email уже существует"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /profile/register [post]
func (h *Handler) Register(gCtx *gin.Context) {
	req := &registerReq{}

	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": "invalid JSON format",
		})
		return
	}

	if req.Password == "" {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": "password is empty",
		})
		return
	}

	if req.Email == "" {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": "email is empty",
		})
		return
	}

	if req.Name == "" {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": "name is empty",
		})
		return
	}

	hashedPassword, err := generateHashString(req.Password)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":      "error",
			"description": "failed to hash password",
		})
		return
	}

	user := &ds.User{
		ID:          uuid.New(),
		IsModerator: false,
		Email:       req.Email,
		Name:        req.Name,
		Password:    hashedPassword,
	}

	err = h.Repository.Register(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			gCtx.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"status":      "error",
				"description": "user with this email already exists",
			})
		} else {
			gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":      "error",
				"description": "failed to create user",
			})
		}
		return
	}

	gCtx.JSON(http.StatusOK, user)
}

// Logout godoc
// @Summary Выход пользователя из системы
// @Description Добавление токена в черный список для выхода из системы
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} successResponse "Успешный выход"
// @Failure 400 {object} errorResponse "Отсутствует или неверный заголовок авторизации"
// @Failure 500 {object} errorResponse "Ошибка при выходе"
// @Router /profile/logout [post]
func (h *Handler) Logout(gCtx *gin.Context) {
	authHeader := gCtx.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		gCtx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": "authorization header missing or invalid format",
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	ctx := context.Background()
	err := h.Redis.GetClient().Set(ctx, "blacklist:"+token, "1", time.Hour*24).Err()
	if err != nil {
		gCtx.JSON(http.StatusInternalServerError, gin.H{
			"status":      "error",
			"description": "failed to logout",
		})
		return
	}

	gCtx.JSON(http.StatusOK, gin.H{
		"message": "logged out",
	})
}

func generateHashString(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
