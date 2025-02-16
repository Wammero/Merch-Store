package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// Запись транзакции
func (repo *PGRepo) RecordTransaction(ctx context.Context, tx pgx.Tx, amount int64, senderUsername string, receiverUsername string) error {
	query := `
       WITH sender_update AS (
           UPDATE users
           SET balance = balance - $1
           WHERE username = $2 AND balance >= $1
           RETURNING user_id, balance
       )
       , receiver_update AS (
           UPDATE users
           SET balance = balance + $1
           WHERE username = $3 
           AND EXISTS (SELECT 1 FROM sender_update)
           RETURNING user_id, username
       )
       , transaction_insert AS (
           INSERT INTO transactions (sender_id, receiver_id, amount)
           SELECT sender_update.user_id, receiver_update.user_id, $1
           FROM sender_update, receiver_update
       )
       SELECT balance FROM sender_update;
   `
	var balance int64
	err := tx.QueryRow(ctx, query, amount, senderUsername, receiverUsername).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("транзакция не выполнена: недостаточно средств или неверные пользователи")
		}
		return fmt.Errorf("ошибка при записи транзакции: %w", err)
	}
	return nil
}
