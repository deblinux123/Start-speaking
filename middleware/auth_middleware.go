package middleware

import (
	"fmt"
	"net/http"
	"start-speek/db"
	"start-speek/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {

		authHeader := context.GetHeader("Authorization")

		if authHeader == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required.",
			})
			return
		}

		tokenString := strings.TrimSpace(
			strings.TrimPrefix(authHeader, "Bearer "),
		)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return utils.JWTKey, nil
		})

		if err != nil || !token.Valid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token.",
			})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		userId := uint(claims["user_id"].(float64))

		// redis session check
		sessionKey := fmt.Sprintf("session:user:%d", userId)

		val, err := db.RedisClient.Get(db.Ctx, sessionKey).Result()

		if err != nil || val != "active" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session expire or logged out"})
			return
		}

		context.Set("user_id", userId)

		context.Next()
	}
}
