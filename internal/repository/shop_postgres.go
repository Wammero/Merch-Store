package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (repo *PGRepo) BuyMerchandise(ctx context.Context, username string, itemName string, amount int64) error {
	// Начинаем транзакцию
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID int
	var userBalance int64
	var merchPrice int64

	// Получаем ID пользователя и его баланс
	err = tx.QueryRow(ctx, `SELECT user_id, balance FROM users WHERE username = $1`, username).Scan(&userID, &userBalance)
	if err != nil {
		return fmt.Errorf("ошибка при получении данных пользователя: %w", err)
	}

	// Получаем цену товара по имени
	err = tx.QueryRow(ctx, `SELECT price FROM merchandise WHERE name = $1`, itemName).Scan(&merchPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("товар с именем '%s' не найден", itemName)
		}
		return fmt.Errorf("ошибка при получении данных товара: %w", err)
	}

	// Проверяем, достаточно ли средств на балансе
	if userBalance < merchPrice*amount {
		return errors.New("недостаточно средств для покупки")
	}

	// Обновляем баланс пользователя
	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance - $1 WHERE user_id = $2`, merchPrice*amount, userID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса пользователя: %w", err)
	}

	// Получаем ID товара по имени
	var merchID int
	err = tx.QueryRow(ctx, `SELECT merch_id FROM merchandise WHERE name = $1`, itemName).Scan(&merchID)
	if err != nil {
		return fmt.Errorf("ошибка при получении ID товара: %w", err)
	}

	// Создаем запись о покупке
	_, err = tx.Exec(ctx, `INSERT INTO purchases (user_id, merch_id, timestamp) VALUES ($1, $2, CURRENT_TIMESTAMP)`, userID, merchID)
	if err != nil {
		return fmt.Errorf("ошибка при записи покупки: %w", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return nil
}
