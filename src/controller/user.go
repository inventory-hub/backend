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
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

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

	user = model.DraftUser{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		RoleName:  input.Role,
		RoleID:    model.RoleMap[input.Role],
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		RoleID:    draftUser.RoleID,
		RoleName:  draftUser.RoleName,
		FirstName: draftUser.FirstName,
		LastName:  draftUser.LastName,
		Username:  input.Username,
		Email:     draftUser.Email,
		Password:  input.Password,
	}
	_, err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func GetListOfUsers(context *gin.Context) {
	searchParam := context.Query("search")
	pageParam := context.DefaultQuery("page", "1")
	pageSizeParam := context.DefaultQuery("pageSize", "2")

	users, err := model.GetUsersWithParam(searchParam)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleId, err := auth.GetRoleFromToken(context)
	if err != nil || roleId == 0 {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed at receiving the role"})
		return
	}

	users = model.FilterUsersByRoleId(users, roleId)
	if users == nil {
		context.JSON(http.StatusOK, gin.H{"users": []model.User{}})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil || pageSize > len(users) {
		context.JSON(http.StatusOK, gin.H{"users": users})
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page*pageSize > len(users) {
		context.JSON(http.StatusOK, gin.H{"users": users})
		return
	}

	context.JSON(http.StatusOK, gin.H{"users": users[(page-1)*pageSize : page*pageSize], "totalPages": uint(len(users) / pageSize)})
}

func GetUserById(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := model.GetUserById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roleId, err := auth.GetRoleFromToken(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.RoleID < roleId {
		context.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("The user is inaccessible with the id: %d", user.ID)})
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": user.ID, "firstName": user.FirstName, "lastName": user.LastName, "role": user.RoleName, "email": user.Email})
}

func UpdateUser(context *gin.Context) {
	var input model.DraftUserPayload

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := model.GetUserById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	role, err := auth.GetRoleFromToken(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.RoleID < role {
		context.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("The user is inaccessible with the id: %d", user.ID)})
		return
	}

	user = model.User{
		ID:        user.ID,
		RoleID:    model.RoleMap[input.Role],
		RoleName:  input.Role,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Username:  user.Username,
		Password:  user.Password,
	}
	err = model.UpdateUser(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

func DeleteUser(context *gin.Context) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := model.GetUserById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	role, err := auth.GetRoleFromToken(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.RoleID < role {
		context.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("The user is inaccessible with the id: %d", user.ID)})
		return
	}

	err = model.DeleteUser(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}
