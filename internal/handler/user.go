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
	err := api.db.AuthenticateUser(context.Background(), user.Username, user.Password)
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

func GetInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetInfo successful"))
}
