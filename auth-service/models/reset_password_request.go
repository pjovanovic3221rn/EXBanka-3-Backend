package models

type ResetPasswordRequest struct {
	ResetToken      string `json:"reset_token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}