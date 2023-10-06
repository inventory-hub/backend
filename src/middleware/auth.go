package middleware

import (
	"Smart-Machine/backend/src/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateAuthorization(roleList []uint) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := util.ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization is required"})
			context.Abort()
			return
		}
		err = util.VerifyRoleInList(context, roleList)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": err})
			context.Abort()
			return
		}
		context.Next()
	}
}
