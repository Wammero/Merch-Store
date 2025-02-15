package repository

import (
	"Merch-Store/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

var ErrUserNotFound = errors.New("user not found")

// Получение ID пользователя по имени
func (repo *PGRepo) GetUserID(ctx context.Context, username string) (int, error) {
	var userID int
	err := repo.pool.QueryRow(ctx, `SELECT user_id FROM users WHERE username = $1`, username).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("пользователь не найден: %w", err)
	}
	return userID, nil
}

func (repo *PGRepo) GetUserBalance(ctx context.Context, tx pgx.Tx, username string) (int, int64, error) {
	var userID int
	var balance int64

	err := tx.QueryRow(ctx, `SELECT user_id, balance FROM users WHERE username = $1 FOR UPDATE`, username).
		Scan(&userID, &balance)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка при получении данных пользователя: %w", err)
	}

	return userID, balance, nil
}

// Получение учетных данных пользователя (пароля и соли)
func (repo *PGRepo) GetUserCredentials(ctx context.Context, username string) (string, string, error) {
	var storedHashedPassword, storedSalt string

	query := `SELECT password_hash, salt FROM users WHERE username = $1`
	err := repo.pool.QueryRow(ctx, query, username).Scan(&storedHashedPassword, &storedSalt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", ErrUserNotFound
		}
		return "", "", err
	}

	return storedHashedPassword, storedSalt, nil
}

// Создание нового пользователя
func (repo *PGRepo) CreateUser(ctx context.Context, username, hashedPassword, salt string) error {
	query := `
		INSERT INTO users (username, password_hash, salt, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := repo.pool.Exec(ctx, query, username, hashedPassword, salt, time.Now(), time.Now())
	return err
}

// Обновление баланса пользователя
func (repo *PGRepo) UpdateUserBalance(ctx context.Context, tx pgx.Tx, userID int, amount int64) error {
	_, err := tx.Exec(ctx, `UPDATE users SET balance = balance + $1 WHERE user_id = $2`, amount, userID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса пользователя: %w", err)
	}
	return nil
}

func (repo *PGRepo) GetUserInventory(ctx context.Context, userID int) ([]model.InventoryItem, error) {
	// Получаем инвентарь пользователя
	rows, err := repo.pool.Query(ctx, `SELECT m.name, COUNT(p.purchase_id) 
	FROM merchandise m 
	LEFT JOIN purchases p ON m.merch_id = p.merch_id 
	WHERE p.user_id = $1 
	GROUP BY m.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []model.InventoryItem
	for rows.Next() {
		var item model.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, err
		}
		inventory = append(inventory, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return inventory, nil
}

func (repo *PGRepo) GetUserCoinHistory(ctx context.Context, userID int) (model.CoinHistory, error) {
	// Получаем историю монет пользователя (полученные и отправленные транзакции)
	var history model.CoinHistory

	// Полученные транзакции
	rows, err := repo.pool.Query(ctx, `SELECT u.username, t.amount 
	FROM transactions t 
	JOIN users u ON u.user_id = t.sender_id 
	WHERE t.receiver_id = $1`, userID)
	if err != nil {
		return history, err
	}
	defer rows.Close()

	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(&tx.FromUser, &tx.Amount); err != nil {
			return history, err
		}
		history.Received = append(history.Received, tx)
	}
	if err := rows.Err(); err != nil {
		return history, err
	}

	// Отправленные транзакции
	rows, err = repo.pool.Query(ctx, `SELECT u.username, t.amount 
	FROM transactions t 
	JOIN users u ON u.user_id = t.receiver_id 
	WHERE t.sender_id = $1`, userID)
	if err != nil {
		return history, err
	}
	defer rows.Close()

	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(&tx.FromUser, &tx.Amount); err != nil {
			return history, err
		}
		history.Sent = append(history.Sent, tx)
	}
	if err := rows.Err(); err != nil {
		return history, err
	}

	return history, nil
}
