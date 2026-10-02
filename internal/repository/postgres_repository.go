package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tayane/payflow-engine/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository cria uma nova instância do repositório PostgreSQL.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Save insere uma nova transação no banco de dados.
func (r *PostgresRepository) Save(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, account_id, amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		tx.ID,
		tx.AccountID,
		tx.Amount,
		tx.Status,
		tx.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save transaction in postgres: %w", err)
	}

	return nil
}

// FindByID busca uma transação existente pelo seu ID.
func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `
		SELECT id, account_id, amount, status, created_at
		FROM transactions
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var tx domain.Transaction
	err := row.Scan(
		&tx.ID,
		&tx.AccountID,
		&tx.Amount,
		&tx.Status,
		&tx.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Retorna nil se a transação não for encontrada
		}
		return nil, fmt.Errorf("failed to find transaction by id: %w", err)
	}

	return &tx, nil
}

// UpdateStatus atualiza apenas o status de uma transação existente.
func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, status domain.TransactionStatus) error {
	query := `
		UPDATE transactions
		SET status = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no transaction found with id %s to update", id)
	}

	return nil
}
