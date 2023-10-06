package model

type Refresh struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
