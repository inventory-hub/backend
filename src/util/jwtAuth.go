package util

import (
	"Smart-Machine/backend/src/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// check for valid admin token
func JWTAuthAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		error := ValidateAdminRoleJWT(context)
		if error != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Only Administrator is allowed to perform this action"})
			context.Abort()
			return
		}
		context.Next()
	}
}

// check for valid customer token
func JWTAuthManager() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		error := ValidateManagerRoleJWT(context)
		if error != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Only registered Managers are allowed to perform this action"})
			context.Abort()
			return
		}
		context.Next()
	}
}

// check for valid customer token
func JWTAuthCustomer() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		error := ValidateCustomerRoleJWT(context)
		if error != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Only registered Customers are allowed to perform this action"})
			context.Abort()
			return
		}
		context.Next()
	}
}

func VerifyRoleInList(context *gin.Context, roleList []uint) error {
	var user model.User

	if user = CurrentUser(context); user == (model.User{}) {
		return fmt.Errorf("no such user with the provided token")
	}

	for i := 0; i < len(roleList); i++ {
		if user.ID == roleList[i] {
			return nil
		}
	}
	return fmt.Errorf("current action is not allowed for this role")
}
