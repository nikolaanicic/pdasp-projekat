package dto

type UserLoginDto struct {
	UserID string `json:"user_id" binding:"required" form:"user_id"`
}
