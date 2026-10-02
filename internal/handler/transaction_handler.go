package handler

import (
	"encoding/json"
	"net/http"

	"github.com/tayane/payflow-engine/internal/domain"
)

type TransactionHandler struct {
	useCase domain.TransactionUseCase
}

func NewTransactionHandler(uc domain.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		useCase: uc,
	}
}

// CreateTransactionInput define a DTO de entrada do JSON
type CreateTransactionInput struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateTransactionInput

	// 1. Decodifica o JSON do corpo da requisição
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid JSON payload",
		})
		return
	}

	// 2. Chama a regra de negócio (Caso de Uso)
	tx, err := h.useCase.CreateTransaction(r.Context(), input.AccountID, input.Amount)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 3. Retorna a transação criada com status 201 Created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(tx)
}
