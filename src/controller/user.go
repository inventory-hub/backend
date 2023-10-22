package controller

import (
	"Smart-Machine/backend/src/model"
	"Smart-Machine/backend/src/util/auth"
	"Smart-Machine/backend/src/util/msgqueue"
	"Smart-Machine/backend/src/util/random"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Login(context *gin.Context) {
	var input model.LoginPayload
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

	jwtToken, err := auth.GenerateBasicJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := auth.GenerateTimeoutJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"accessToken": jwtToken, "refreshToken": refreshToken})
}

func Refresh(context *gin.Context) {
	var input model.RefreshPayload
	var user model.User

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user = auth.CurrentUser(input); user == (model.User{}) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No such user with the provided token"})
		return
	}

	jwtToken, err := auth.GenerateBasicJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := auth.GenerateTimeoutJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"accessToken": jwtToken, "refreshToken": refreshToken})
}

func Invite(context *gin.Context) {
	var input model.DraftUserPayload
	var user model.DraftUser

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleId, err := strconv.Atoi(input.RoleID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	user = model.DraftUser{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		RoleID:    uint(roleId),
	}

	inviteToken, err := random.PseudoUUID()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.InviteToken = inviteToken

	_, err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := msgqueue.Message{
		To:          input.Email,
		FullName:    fmt.Sprintf("%s %s", input.FirstName, input.LastName),
		InviteToken: inviteToken,
		CallbackURL: "https://inventory-hub.space/sign-up",
	}
	encodedMessage, err := json.Marshal(message)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msgqueue.QueueClientEnqueueMessage(string(encodedMessage))

	context.JSON(http.StatusCreated, gin.H{})
}

func Register(context *gin.Context) {
	var input model.RegisterPayload

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draftUser, err := model.GetDraftUserByInvitationToken(input.InviteToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	user := model.User{
		RoleID:    draftUser.RoleID,
		FirstName: draftUser.FirstName,
		LastName:  draftUser.LastName,
		Username:  input.Username,
		Email:     draftUser.Email,
		Password:  input.Password,
	}
	_, err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	jwtToken, err := auth.GenerateBasicJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := auth.GenerateTimeoutJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"accessToken": jwtToken, "refreshToken": refreshToken})
}
