package model

type RegisterPayload struct {
	InviteToken string `json:"token" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
}
