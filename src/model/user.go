package model

import (
	"Smart-Machine/backend/src/database"
	"fmt"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primary_key"`
	RoleID    uint   `gorm:"not null;DEFAULT:3" json:"roleId"`
	FirstName string `gorm:"size:255;not null;" json:"firstName"`
	LastName  string `gorm:"size:255;not null;" json:"lastName"`
	Username  string `gorm:"size:255;not null;" json:"username"`
	Email     string `gorm:"size:255;not null;unique" json:"email"`
	Password  string `gorm:"size:255;not null" json:"-"`
	Role      Role   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
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

func UpdateUser(user User) error {
	err := database.DB.Omit("password").Updates(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(user User) error {
	err := database.DB.Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
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

func GetUsersWithParam(search string) ([]User, error) {
	var users []User
	search = fmt.Sprintf("%%%s%%", search)
	err := database.DB.Where("first_name LIKE ?", search).Or("last_name LIKE ?", search).Or("username LIKE ?", search).Or("email LIKE ?", search).Find(&users).Error
	if err != nil {
		return []User{}, err
	}
	return users, nil
}

func FilterUsersByRole(users []User, role uint) []User {
	var filteredUsers []User
	for i := 0; i < len(users); i++ {
		if users[i].RoleID >= role {
			filteredUsers = append(filteredUsers, users[i])
		}
	}
	return filteredUsers
}
