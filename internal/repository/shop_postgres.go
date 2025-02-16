package repository

import (
	"context"
	"errors"
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
func (repo *PGRepo) InsertPurchase(ctx context.Context, tx pgx.Tx, username, merchname string) error {
	query := `
       WITH merch_take AS (
           SELECT merch_id, price FROM merchandise
           WHERE name = $1
       )
       , sender_update AS (
           UPDATE users
           SET balance = balance - (SELECT price FROM merch_take)
           WHERE username = $2 AND balance >= (SELECT price FROM merch_take)
           RETURNING user_id
       )
       , purchase_insert AS (
           INSERT INTO purchases (user_id, merch_id)
           SELECT su.user_id, mt.merch_id
           FROM sender_update su, merch_take mt
           RETURNING purchase_id
       )
       SELECT purchase_id FROM purchase_insert;
  	`
	var purchaseID int64
	err := tx.QueryRow(ctx, query, merchname, username).Scan(&purchaseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("покупка не выполнена: недостаточно средств или товар не найден")
		}
		return fmt.Errorf("ошибка при записи покупки: %w", err)
	}
	return nil
}
