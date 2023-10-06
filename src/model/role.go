package model

import (
	"Smart-Machine/backend/src/database"

	"gorm.io/gorm"
)

// Role model
type Role struct {
	gorm.Model
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"size:50;not null;unique" json:"name"`
	Description string `gorm:"size:255;not null" json:"description"`
}

// Create a role
func CreateRole(Role *Role) (err error) {
	err = database.DB.Create(Role).Error
	if err != nil {
		return err
	}
	return nil
}
