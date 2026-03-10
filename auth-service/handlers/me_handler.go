package handlers

import (
	"encoding/json"
	"net/http"

	"auth-service/middleware"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	credentialID, _ := r.Context().Value(middleware.ContextCredentialIDKey).(int64)
	employeeID, _ := r.Context().Value(middleware.ContextEmployeeIDKey).(int64)
	email, _ := r.Context().Value(middleware.ContextEmailKey).(string)

	response := map[string]any{
		"message":       "authorized request",
		"credential_id": credentialID,
		"employee_id":   employeeID,
		"email":         email,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}