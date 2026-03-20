package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RAF-SI-2025/EXBanka-3-Backend/account-service/internal/config"
	"github.com/RAF-SI-2025/EXBanka-3-Backend/account-service/internal/models"
	"github.com/RAF-SI-2025/EXBanka-3-Backend/account-service/internal/repository"
	"github.com/RAF-SI-2025/EXBanka-3-Backend/account-service/internal/service"
	"github.com/RAF-SI-2025/EXBanka-3-Backend/account-service/internal/util"
	"gorm.io/gorm"
)

type CreateAccountHTTPHandler struct {
	svc *service.AccountService
	db  *gorm.DB
	cfg *config.Config
}

func NewCreateAccountHTTPHandler(db *gorm.DB, cfg *config.Config) *CreateAccountHTTPHandler {
	accountRepo := repository.NewAccountRepository(db)
	currencyRepo := repository.NewCurrencyRepository(db)
	notifSvc := service.NewNotificationService(cfg)
	svc := service.NewAccountServiceWithRepos(accountRepo, currencyRepo, notifSvc)
	return &CreateAccountHTTPHandler{svc: svc, db: db, cfg: cfg}
}

type createAccountHTTPRequest struct {
	ClientID      uint    `json:"clientId"`
	FirmaID       uint    `json:"firmaId"`
	CurrencyID    uint    `json:"currencyId"`
	Tip           string  `json:"tip"`
	Vrsta         string  `json:"vrsta"`
	Naziv         string  `json:"naziv"`
	PocetnoStanje float64 `json:"pocetnoStanje"`
}

func (h *CreateAccountHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var employeeID uint
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := util.ParseToken(tokenStr, h.cfg.JWTSecret)
		if err == nil && claims.EmployeeID != 0 {
			employeeID = claims.EmployeeID
		}
	}

	var req createAccountHTTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	input := service.CreateAccountInput{
		CurrencyID:    req.CurrencyID,
		Tip:           req.Tip,
		Vrsta:         req.Vrsta,
		Naziv:         req.Naziv,
		PocetnoStanje: req.PocetnoStanje,
	}
	if req.ClientID != 0 {
		input.ClientID = &req.ClientID
		var client models.Client
		if err := h.db.First(&client, req.ClientID).Error; err == nil {
			input.ClientEmail = client.Email
			input.ClientName = client.Ime + " " + client.Prezime
		}
	}
	if req.FirmaID != 0 {
		input.FirmaID = &req.FirmaID
	}
	if employeeID != 0 {
		input.ZaposleniID = &employeeID
	}

	acc, err := h.svc.CreateAccount(input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"account": acc,
		"message": "Račun uspešno kreiran",
	})
}
