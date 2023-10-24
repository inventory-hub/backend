package model

import (
	"Smart-Machine/backend/src/database"

	"gorm.io/gorm"
)

type CategoryPayload struct {
	Name string `json:"name" binding:"required"`
}

type Category struct {
	gorm.Model
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"size:255;not null;" json:"name"`
	Quantity uint   `gorm:"DEFAULT:0" json:"itemQuantity"`
}

func (category *Category) Save() (*Category, error) {
	err := database.DB.Create(&category).Error
	if err != nil {
		return &Category{}, err
	}
	return category, nil
}

func GetCategories() (categories []Category, err error) {
	err = database.DB.Find(&categories).Error
	if err != nil {
		return []Category{}, err
	}
	return categories, nil
}

func CreateCategory(name string) (category Category, err error) {
	category = Category{Name: name}
	err = database.DB.Save(&category).Error
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func DeleteCategoryByName(name string) (err error) {
	err = database.DB.Where("name LIKE ?", name).Delete(&Category{}).Error
	if err != nil {
		return err
	}
	return nil
}
