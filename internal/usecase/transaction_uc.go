package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tayane/payflow-engine/internal/domain"
)

type transactionUseCase struct {
	repo domain.TransactionRepository
}

// NewTransactionUseCase inicializa o caso de uso com o repositório injetado
func NewTransactionUseCase(repo domain.TransactionRepository) domain.TransactionUseCase {
	return &transactionUseCase{
		repo: repo,
	}
}

func (uc *transactionUseCase) CreateTransaction(ctx context.Context, accountID string, amount float64) (*domain.Transaction, error) {
	// Validações das regras de negócio
	if accountID == "" {
		return nil, errors.New("account_id is required")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Construção do objeto de domínio
	tx := &domain.Transaction{
		ID:        uuid.New().String(), // Gera um UUID único para a transação
		AccountID: accountID,
		Amount:    amount,
		Status:    domain.StatusPending,
		CreatedAt: time.Now().UTC(),
	}

	// Persiste a transação via repositório
	if err := uc.repo.Save(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return tx, nil
}
