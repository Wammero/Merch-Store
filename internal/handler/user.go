package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"Merch-Store/internal/model"
	"Merch-Store/pkg/jwt"
)

type AuthResponse struct {
	Token string `json:"token"`
}

func (api *API) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Вызов метода из репозитория через api.db
	err := api.service.AuthenticateUser(context.Background(), user.Username, user.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Генерация JWT токена
	tokenString, err := jwt.GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Token: tokenString}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *API) GetInfo(w http.ResponseWriter, r *http.Request) {
	// Извлекаем имя пользователя из контекста (например, после успешной аутентификации через JWT)
	username, ok := r.Context().Value(jwt.UserContextKey).(string)
	if !ok || username == "" {
		http.Error(w, `{"errors":"Не удалось определить отправителя"}`, http.StatusUnauthorized)
		return
	}

	// Получаем информацию о пользователе
	response, err := api.service.GetUserInfo(r.Context(), username)
	if err != nil {
		http.Error(w, `{"errors":"Ошибка получения данных"}`, http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки для ответа (Content-Type: application/json)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Создаем отформатированный JSON-ответ
	indentedResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, `{"errors":"Ошибка при отправке данных"}`, http.StatusInternalServerError)
		return
	}

	// Отправляем отформатированный JSON с данными пользователя
	w.Write(indentedResponse)
}
