package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/RAF-SI-2025/EXBanka-3-Backend/transfer-service/internal/models"
)

// TransferVerifierInterface is the minimal interface needed for verification.
type TransferVerifierInterface interface {
	VerifyTransfer(transferID uint, verificationCode string) (*models.Transfer, error)
}

// VerifyTransferHTTPHandler handles POST /api/v1/transfers/{id}/verify
type VerifyTransferHTTPHandler struct {
	svc TransferVerifierInterface
}

func NewVerifyTransferHTTPHandler(svc TransferVerifierInterface) *VerifyTransferHTTPHandler {
	return &VerifyTransferHTTPHandler{svc: svc}
}

func (h *VerifyTransferHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract {id} from path: /api/v1/transfers/{id}/verify
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// Expected: ["api", "v1", "transfers", "{id}", "verify"]
	if len(parts) < 5 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	idStr := parts[3]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid transfer id", http.StatusBadRequest)
		return
	}

	var body struct {
		VerificationCode string `json:"verification_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.VerificationCode == "" {
		http.Error(w, "verification_code required", http.StatusBadRequest)
		return
	}

	transfer, err := h.svc.VerifyTransfer(uint(id), body.VerificationCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     transfer.ID,
		"status": transfer.Status,
	})
}

// CombinedTransferHandler routes /transfers/{id}/verify to the verify handler,
// and all other /transfers requests to the gRPC gateway.
type CombinedTransferHandler struct {
	verify   *VerifyTransferHTTPHandler
	fallback http.Handler
}

func NewCombinedTransferHandler(verify *VerifyTransferHTTPHandler, fallback http.Handler) *CombinedTransferHandler {
	return &CombinedTransferHandler{verify: verify, fallback: fallback}
}

func (h *CombinedTransferHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/verify") {
		h.verify.ServeHTTP(w, r)
		return
	}
	h.fallback.ServeHTTP(w, r)
}
