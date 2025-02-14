package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// Запись транзакции
func (repo *PGRepo) RecordTransaction(ctx context.Context, tx pgx.Tx, senderID, receiverID int, amount int64) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO transactions (sender_id, receiver_id, amount)
		VALUES ($1, $2, $3)
	`, senderID, receiverID, amount)
	if err != nil {
		return fmt.Errorf("ошибка при записи транзакции: %w", err)
	}
	return nil
}
