package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Welcome(context *gin.Context) {
	context.JSON{http.StatusOK, gin.H{"message": "Welcome to Inventory Hub API"}}
}
