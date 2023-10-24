package controller

import (
	"Smart-Machine/backend/src/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCategories(context *gin.Context) {
	categories, err := model.GetCategories()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"categories": categories})
}

func CreateCategory(context *gin.Context) {
	var input model.CategoryPayload

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := model.CreateCategory(input.Name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": category.ID, "name": category.Name})
}

func DeleteCategory(context *gin.Context) {
	name := context.Param("name")

	err := model.DeleteCategoryByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}
