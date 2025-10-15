package ds

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserUUID    uuid.UUID `json:"user_uuid"`
	IsProfessor bool      `json:"is_professor"`
}
