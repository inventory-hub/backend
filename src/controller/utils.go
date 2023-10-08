package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Welcome(context *gin.Context) {
	context.JSON(http.StatusOK, "Welcome to the Inventory Hub API")
}
