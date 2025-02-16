package service

import (
	"context"
	"fmt"
)

func (s *Service) BuyMerchandise(ctx context.Context, username, itemName string, amount int64) error {
	// Начинаем транзакцию
	tx, err := s.repo.GetPool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	// Записываем покупку
	if err := s.repo.InsertPurchase(ctx, tx, username, itemName); err != nil {
		return err
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return nil
}
