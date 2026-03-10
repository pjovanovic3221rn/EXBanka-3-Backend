package models

type ActivateRequest struct {
	ActivationToken  string `json:"activation_token"`
	Password         string `json:"password"`
	ConfirmPassword  string `json:"confirm_password"`
}