package model

import (
	"Smart-Machine/backend/src/database"
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DraftUserPayload struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

type DraftUser struct {
	ID          uint           `gorm:"primary_key" json:"id"`
	Username    string         `gorm:"size:255;not null;" json:"username"`
	FirstName   string         `gorm:"size:255;not null;" json:"firstName"`
	LastName    string         `gorm:"size:255;not null;" json:"lastName"`
	Email       string         `gorm:"size:255;not null;unique" json:"email"`
	RoleName    string         `gorm:"size:255;not null;" json:"role"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	InviteToken string         `gorm:"size:255;not null;unique" json:"-"`
	RoleID      uint           `gorm:"not null;DEFAULT:3" json:"-"`
	Role        Role           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (user *DraftUser) Save() (*DraftUser, error) {
	err := database.DB.Create(&user).Error
	if err != nil {
		return &DraftUser{}, err
	}
	return user, nil
}

// GORM Hook on the object creation
func (user *DraftUser) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))

	return nil
}

func GetDraftUserByInvitationToken(inviteToken string) (DraftUser, error) {
	var draftUser DraftUser

	err := database.DB.Where("invite_token = ?", inviteToken).Find(&draftUser).Error
	if err != nil {
		return DraftUser{}, err
	}

	return draftUser, nil
}
