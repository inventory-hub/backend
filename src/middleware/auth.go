package middleware

import (
	"Smart-Machine/backend/src/util/auth"
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func ValidateAuthorization(roleList []uint) gin.HandlerFunc {
	return func(context *gin.Context) {
		var err error

		if reflect.DeepEqual(roleList, AdminRole) {
			err = auth.ValidateAdminRoleJWT(context)
		} else if reflect.DeepEqual(roleList, AuthorizedRoles) {
			adminErr := auth.ValidateAdminRoleJWT(context)
			managerErr := auth.ValidateManagerRoleJWT(context)
			if adminErr != nil && managerErr != nil {
				err = errors.New("invalid authorization token provided")
			}
		} else if reflect.DeepEqual(roleList, AllRoles) {
			err = auth.ValidateJWT(context)
		}

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization is required"})
			context.Abort()
			return
		}

		context.Next()
	}
}
