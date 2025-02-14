package service

import (
	"context"
	"errors"
	"fmt"
)

func (s *Service) BuyMerchandise(ctx context.Context, username, itemName string, amount int64) error {
	// Начинаем транзакцию
	tx, err := s.repo.GetPool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	userID, userBalance, err := s.repo.GetUserBalance(ctx, username)
	if err != nil {
		return err
	}

	merchID, merchPrice, err := s.repo.GetMerchInfo(ctx, tx, itemName)
	if err != nil {
		return err
	}

	// Проверяем баланс
	if userBalance < merchPrice*amount {
		return errors.New("недостаточно средств для покупки")
	}

	// Обновляем баланс пользователя
	if err := s.repo.UpdateUserBalance(ctx, tx, userID, -merchPrice*amount); err != nil {
		return err
	}

	// Записываем покупку
	if err := s.repo.InsertPurchase(ctx, tx, userID, merchID); err != nil {
		return err
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return nil
}
