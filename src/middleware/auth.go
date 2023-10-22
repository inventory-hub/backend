package middleware

import (
	"Smart-Machine/backend/src/util/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateAuthorization(roleList []uint) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := auth.ValidateJWTWithRole(context, roleList)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization is required"})
			context.Abort()
			return
		}

		context.Next()
	}
}
