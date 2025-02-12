package repository

import (
	"context"
	"errors"
	"fmt"
)

func (repo *PGRepo) SendCoin(ctx context.Context, senderUsername, receiverUsername string, amount int64) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	var senderID, receiverID int
	var senderBalance int64

	// Проверка баланса отправителя
	err = tx.QueryRow(ctx, `
		SELECT user_id, balance FROM users WHERE username = $1
	`, senderUsername).Scan(&senderID, &senderBalance)

	if err != nil {
		return fmt.Errorf("ошибка при получении данных отправителя: %w", err)
	}

	if senderBalance < amount {
		return errors.New("недостаточно средств для перевода")
	}

	// Получение ID получателя
	err = tx.QueryRow(ctx, `SELECT user_id FROM users WHERE username = $1`, receiverUsername).Scan(&receiverID)
	if err != nil {
		return fmt.Errorf("получатель не найден: %w", err)
	}

	// Обновление баланса отправителя и получателя
	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance - $1 WHERE user_id = $2`, amount, senderID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса отправителя: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance + $1 WHERE user_id = $2`, amount, receiverID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса получателя: %w", err)
	}

	// Запись транзакции
	_, err = tx.Exec(ctx, `
		INSERT INTO transactions (sender_id, receiver_id, amount)
		VALUES ($1, $2, $3)
	`, senderID, receiverID, amount)
	if err != nil {
		return fmt.Errorf("ошибка при записи транзакции: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return nil
}
