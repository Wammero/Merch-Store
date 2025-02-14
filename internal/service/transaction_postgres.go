package service

import (
	"context"
	"errors"
	"fmt"
)

func (s *Service) SendCoin(ctx context.Context, senderUsername, receiverUsername string, amount int64) error {
	tx, err := s.repo.GetPool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	// Получение данных отправителя
	senderID, senderBalance, err := s.repo.GetUserBalance(ctx, senderUsername)
	if err != nil {
		return err
	}

	// Проверка баланса отправителя
	if senderBalance < amount {
		return errors.New("недостаточно средств для перевода")
	}

	// Получение данных получателя
	receiverID, err := s.repo.GetUserID(ctx, receiverUsername)
	if err != nil {
		return err
	}

	// Обновление баланса отправителя и получателя
	err = s.repo.UpdateUserBalance(ctx, tx, senderID, -amount)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUserBalance(ctx, tx, receiverID, amount)
	if err != nil {
		return err
	}

	// Запись транзакции
	err = s.repo.RecordTransaction(ctx, tx, senderID, receiverID, amount)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return nil
}
