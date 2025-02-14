package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// getMerchInfo получает ID и цену товара по имени
func (repo *PGRepo) GetMerchInfo(ctx context.Context, tx pgx.Tx, itemName string) (int, int64, error) {
	var merchID int
	var price int64

	err := tx.QueryRow(ctx, `SELECT merch_id, price FROM merchandise WHERE name = $1`, itemName).Scan(&merchID, &price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, 0, fmt.Errorf("товар с именем '%s' не найден", itemName)
		}
		return 0, 0, fmt.Errorf("ошибка при получении данных товара: %w", err)
	}

	return merchID, price, nil
}

// insertPurchase записывает покупку в базу
func (repo *PGRepo) InsertPurchase(ctx context.Context, tx pgx.Tx, userID, merchID int) error {
	_, err := tx.Exec(ctx, `INSERT INTO purchases (user_id, merch_id, timestamp) VALUES ($1, $2, CURRENT_TIMESTAMP)`, userID, merchID)
	if err != nil {
		return fmt.Errorf("ошибка при записи покупки: %w", err)
	}
	return nil
}
