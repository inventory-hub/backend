package model

type RefreshPayload struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
