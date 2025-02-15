package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"Merch-Store/internal/model"
	"Merch-Store/pkg/jwt"
	"Merch-Store/pkg/responsemaker"
	"Merch-Store/pkg/validators"
)

type AuthResponse struct {
	Token string `json:"token"`
}

func (api *API) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responsemaker.WriteJSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Проверка корректности username и password
	if !validators.IsValidUsername(user.Username) {
		responsemaker.WriteJSONError(w, "Invalid username", http.StatusBadRequest)
		return
	}
	if !validators.IsValidPassword(user.Password) {
		responsemaker.WriteJSONError(w, "Invalid password", http.StatusBadRequest)
		return
	}

	// Вызов метода из репозитория через api.service
	err := api.service.AuthenticateUser(context.Background(), user.Username, user.Password)
	if err != nil {
		responsemaker.WriteJSONError(w, "Authentication error", http.StatusUnauthorized)
		return
	}

	// Генерация JWT токена
	tokenString, err := jwt.GenerateJWT(user.Username)
	if err != nil {
		responsemaker.WriteJSONError(w, "Error during token generation", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Token: tokenString}

	responsemaker.WriteJSONResponse(w, response, http.StatusOK)
}

func (api *API) GetInfo(w http.ResponseWriter, r *http.Request) {
	// Извлекаем имя пользователя из контекста
	username, ok := r.Context().Value(jwt.UserContextKey).(string)
	if !ok || username == "" {
		responsemaker.WriteJSONError(w, "Не удалось определить отправителя", http.StatusUnauthorized)
		return
	}

	// Получаем информацию о пользователе
	response, err := api.service.GetUserInfo(r.Context(), username)
	if err != nil {
		responsemaker.WriteJSONError(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ с отступами
	responsemaker.WriteJSONResponse(w, response, http.StatusOK)
}
