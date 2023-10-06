package controller

import (
	"Smart-Machine/inventory-hub-2/model"
	"Smart-Machine/inventory-hub-2/util"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Login(context *gin.Context) {
	var input model.Login
	var err error

	if err = context.ShouldBindJSON(&input); err != nil {
		var errorMessage string
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			validationError := validationErrors[0]
			log.Print(validationError)
			if validationError.Tag() == "required" {
				errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
			} else if validationError.Tag() == "email" {
				errorMessage = fmt.Sprintf("%s wrong email format", validationError.Field())
			}
		}

		context.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	user, err := model.GetUserByEmail(input.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidateUserPassword(input.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := util.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := util.GenerateRefreshJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"accessToken": jwtToken, "refreshToken": refreshToken})
}

func Refresh(context *gin.Context) {
	var input model.Refresh
	var user model.User

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user = util.CurrentUser(input); user == (model.User{}) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No such user with the provided token"})
		return
	}

	jwtToken, err := util.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := util.GenerateRefreshJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"accessToken": jwtToken, "refreshToken": refreshToken})
}

func GetListOfUsers(context *gin.Context) {
	var users *[]model.User
	page := context.Query("page")
	// pageSize := context.Query("pageSize")
	search := context.Query("search")
	log.Print(search)

	users, err := model.GetUsersWithParam(search)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	log.Print(users)

	user := util.CurrentUser(context)
	if user == (model.User{}) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorized user not found."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"users": users, "totalPages": page})
}
