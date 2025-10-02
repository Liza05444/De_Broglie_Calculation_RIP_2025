package pkg

import (
	"DeBroglieProject/internal/app/ds"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const jwtPrefix = "Bearer "

func (a *Application) WithAuthCheck(requireModerator bool) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":      "error",
				"description": "authorization header missing or invalid format",
			})
			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]

		ctx := context.Background()
		_, err := a.Redis.GetClient().Get(ctx, "blacklist:"+jwtStr).Result()
		if err == nil {
			gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":      "error",
				"description": "token has been revoked",
			})
			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.Config.JWT.Token), nil
		})
		if err != nil {
			gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":      "error",
				"description": "invalid or expired token",
			})
			log.Println(err)
			return
		}

		myClaims := token.Claims.(*ds.JWTClaims)

		gCtx.Set("user_uuid", myClaims.UserUUID)
		gCtx.Set("is_moderator", myClaims.IsModerator)

		if requireModerator && !myClaims.IsModerator {
			gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":      "error",
				"description": "moderator permissions required for this operation",
			})
			return
		}
	}
}
