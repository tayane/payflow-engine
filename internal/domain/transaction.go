package domain

import (
	"context"
	"time"
)

// TransactionStatus representa os possíveis estados de uma transação
type TransactionStatus string

const (
	StatusPending   TransactionStatus = "PENDING"
	StatusCompleted TransactionStatus = "COMPLETED"
	StatusFailed    TransactionStatus = "FAILED"
)

// Transaction é a entidade principal do domínio
type Transaction struct {
	ID        string            `json:"id"`
	AccountID string            `json:"account_id"`
	Amount    float64           `json:"amount"`
	Status    TransactionStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}

// TransactionRepository define o contrato (Port) que qualquer repositório (Postgres, Mock, etc) deve seguir
type TransactionRepository interface {
	Save(ctx context.Context, tx *Transaction) error
	FindByID(ctx context.Context, id string) (*Transaction, error)
	UpdateStatus(ctx context.Context, id string, status TransactionStatus) error
}

// TransactionUseCase define as regras de negócio
type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, accountID string, amount float64) (*Transaction, error)
}
