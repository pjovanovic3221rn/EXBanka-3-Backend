package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"auth-service/config"
	"auth-service/models"
	"auth-service/repository"
	"auth-service/services"
)

type RefreshHandler struct {
	CredentialService *services.CredentialService
}

func NewRefreshHandler(db *sql.DB) *RefreshHandler {
	repo := repository.NewCredentialRepository(db)
	cfg := config.LoadConfig()
	jwtService := services.NewJWTService(
		cfg.JWTSecret,
		cfg.JWTAccessExpirationMinutes,
		cfg.JWTRefreshExpirationHours,
	)
	service := services.NewCredentialService(repo, jwtService)

	return &RefreshHandler{
		CredentialService: service,
	}
}

func (h *RefreshHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RefreshRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.CredentialService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}