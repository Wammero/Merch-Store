package service

import (
	"context"
	"fmt"

	"Merch-Store/internal/model"
	"Merch-Store/internal/repository"
	ps "Merch-Store/pkg/password"
)

func (s *Service) AuthenticateUser(ctx context.Context, username, password string) error {
	// Получаем хеш пароля и соль из БД
	storedHashedPassword, storedSalt, err := s.repo.GetUserCredentials(ctx, username)
	if err != nil && err != repository.ErrUserNotFound {
		return err // Ошибка при запросе к БД
	}

	// Если пользователь найден, проверяем пароль
	if err == nil {
		if !ps.CheckPassword(password, storedSalt, storedHashedPassword) {
			return repository.ErrInvalidCredentials
		}
		return nil // Успешная аутентификация
	}

	// Если пользователя нет, регистрируем нового
	hashedPassword, salt, err := ps.HashPassword(password)
	if err != nil {
		return err
	}

	err = s.repo.CreateUser(ctx, username, hashedPassword, salt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUserInfo(ctx context.Context, username string) (model.UserInfoResponse, error) {
	// Начинаем транзакцию
	tx, err := s.repo.GetPool().Begin(ctx)
	if err != nil {
		return model.UserInfoResponse{}, fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback(ctx) // Откат при ошибке

	userID, balance, err := s.repo.GetUserBalance(ctx, tx, username)
	if err != nil {
		return model.UserInfoResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.UserInfoResponse{}, fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	inventory, err := s.repo.GetUserInventory(ctx, userID)
	if err != nil {
		return model.UserInfoResponse{}, err
	}

	coinHistory, err := s.repo.GetUserCoinHistory(ctx, userID)
	if err != nil {
		return model.UserInfoResponse{}, err
	}

	// Формируем ответ
	response := model.UserInfoResponse{
		Coins:       balance,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}

	return response, nil
}
