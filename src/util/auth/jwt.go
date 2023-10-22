package auth

import (
	"Smart-Machine/backend/src/model"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// retrieve JWT key from .env file
var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

// generate JWT token
func GenerateBasicJWT(user model.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.RoleID,
		"iat":  time.Now().Unix(),
		"eat":  time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

// generate JWT refresh or invite token
func GenerateTimeoutJWT(user model.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"eat": time.Now().Add(time.Minute * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

// validate JWT token
func ValidateJWT(param interface{}) error {
	token, err := getToken(param)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func ValidateJWTWithRole(context *gin.Context, roles []uint) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	role := uint(claims["role"].(float64))
	if ok && token.Valid {
		for i := 0; i < len(roles); i++ {
			if role == roles[i] {
				return nil
			}
		}
	}
	return errors.New("invalid token provided")
}

// fetch user details from the token
func CurrentUser(param interface{}) model.User {
	var err error
	var token *jwt.Token

	switch param.(type) {
	case *gin.Context:
		err = ValidateJWT(param.(*gin.Context))
		if err != nil {
			return model.User{}
		}
		token, _ = getToken(param.(*gin.Context))
	case model.RefreshPayload:
		err = ValidateJWT(param.(model.RefreshPayload))
		if err != nil {
			return model.User{}
		}
		token, _ = getToken(param.(model.RefreshPayload))
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	user, err := model.GetUserById(userId)
	if err != nil {
		return model.User{}
	}
	return user
}

// check token validity
func getToken(param interface{}) (*jwt.Token, error) {
	var tokenString string

	switch param.(type) {
	case *gin.Context:
		tokenString = getTokenFromRequest(param.(*gin.Context))
	case model.RefreshPayload:
		tokenString = param.(model.RefreshPayload).RefreshToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

// extract token from request Authorization header
func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
