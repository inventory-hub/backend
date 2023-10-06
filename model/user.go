package model

import (
	"Smart-Machine/inventory-hub-2/database"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primary_key"`
	RoleID   uint   `gorm:"not null;DEFAULT:3" json:"roleId"`
	Username string `gorm:"size:255;not null;" json:"username"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     Role   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (user *User) Save() (*User, error) {
	err := database.DB.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// GORM Hook on the object creation
func (user *User) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))

	return nil
}

func (user *User) ValidateUserPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := database.DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := database.DB.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserById(id uint) (User, error) {
	var user User
	err := database.DB.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUsersWithParam(search string) (*[]User, error) {
	var users *[]User
	err := database.DB.Where("username = ?", "calin").Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return users, nil
}
