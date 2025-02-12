package repository

import (
	"context"
	"time"

	ps "Merch-Store/pkg/password"

	"github.com/jackc/pgx/v4"
)

// AuthenticateUser проверяет пользователя или регистрирует нового
func (repo *PGRepo) AuthenticateUser(ctx context.Context, username, password string) error {
	var storedHashedPassword, storedSalt string

	// Проверка, существует ли пользователь и получение хеша пароля и соли
	checkQuery := `
		SELECT password_hash, salt FROM users WHERE username = $1
	`
	err := repo.pool.QueryRow(ctx, checkQuery, username).Scan(&storedHashedPassword, &storedSalt)

	if err == nil {
		// Пользователь найден, проверим совпадение пароля
		if !ps.CheckPassword(password, storedSalt, storedHashedPassword) {
			return err
		}
		// Пользователь авторизован, нет необходимости обновлять данные
		return nil
	} else if err != pgx.ErrNoRows {
		// Ошибка запроса, не связанная с отсутствием пользователя
		return err
	}

	// Регистрация нового пользователя
	hashedPassword, salt, err := ps.HashPassword(password)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO users (username, password_hash, salt, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = repo.pool.Exec(ctx, insertQuery, username, hashedPassword, salt, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

// // GetUserInfo получает информацию о пользователе по ID
// func (repo *PGRepo) GetUserInfo(ctx context.Context, userID int64) (*models.User, error) {

// }
