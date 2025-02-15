package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"Merch-Store/internal/model"
	"Merch-Store/pkg/jwt"
	"Merch-Store/pkg/validators"
)

type AuthResponse struct {
	Token string `json:"token"`
}

func (api *API) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJSONError(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверка корректности username и password
	if !validators.IsValidUsername(user.Username) {
		writeJSONError(w, "Некорректное имя пользователя", http.StatusBadRequest)
		return
	}
	if !validators.IsValidPassword(user.Password) {
		writeJSONError(w, "Некорректный пароль", http.StatusBadRequest)
		return
	}

	// Вызов метода из репозитория через api.service
	err := api.service.AuthenticateUser(context.Background(), user.Username, user.Password)
	if err != nil {
		writeJSONError(w, "Ошибка аутентификации", http.StatusUnauthorized)
		return
	}

	// Генерация JWT токена
	tokenString, err := jwt.GenerateJWT(user.Username)
	if err != nil {
		writeJSONError(w, "Ошибка при генерации токена", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Token: tokenString}

	writeJSONResponse(w, response, http.StatusOK)
}

func (api *API) GetInfo(w http.ResponseWriter, r *http.Request) {
	// Извлекаем имя пользователя из контекста
	username, ok := r.Context().Value(jwt.UserContextKey).(string)
	if !ok || username == "" {
		writeJSONError(w, "Не удалось определить отправителя", http.StatusUnauthorized)
		return
	}

	// Получаем информацию о пользователе
	response, err := api.service.GetUserInfo(r.Context(), username)
	if err != nil {
		writeJSONError(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ с отступами
	writeJSONResponse(w, response, http.StatusOK)
}
