package service

import (
	"context"
	"fmt"
)

func (s *Service) SendCoin(ctx context.Context, senderUsername, receiverUsername string, amount int64) error {
	tx, err := s.repo.GetPool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	err = s.repo.RecordTransaction(ctx, tx, amount, senderUsername, receiverUsername)

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return nil
}
