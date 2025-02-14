package service

import (
	"context"

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

func (s *Service) GetUserInfo(ctx context.Context, username string) (repository.UserInfoResponse, error) {
	// Шаг 1: Получаем ID пользователя и баланс
	userID, balance, err := s.repo.GetUserBalance(ctx, username)
	if err != nil {
		return repository.UserInfoResponse{}, err
	}

	// Шаг 2: Получаем инвентарь пользователя
	inventory, err := s.repo.GetUserInventory(ctx, userID)
	if err != nil {
		return repository.UserInfoResponse{}, err
	}

	// Шаг 3: Получаем историю монет (полученные и отправленные транзакции)
	coinHistory, err := s.repo.GetUserCoinHistory(ctx, userID)
	if err != nil {
		return repository.UserInfoResponse{}, err
	}

	// Формируем ответ
	response := repository.UserInfoResponse{
		Coins:       balance,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}

	return response, nil
}
