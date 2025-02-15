package handler

import (
	"encoding/json"
	"net/http"

	"Merch-Store/pkg/jwt"
	"Merch-Store/pkg/validators"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int64  `json:"amount"`
}

func (api *API) SendCoin(w http.ResponseWriter, r *http.Request) {
	var req SendCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if !validators.IsValidUsername(req.ToUser) || req.Amount <= 0 {
		writeJSONError(w, "Имя пользователя и количество монет обязательны", http.StatusBadRequest)
		return
	}

	// Получаем имя отправителя из контекста
	sender, ok := r.Context().Value(jwt.UserContextKey).(string)
	if !ok || sender == "" {
		writeJSONError(w, "Не удалось определить отправителя", http.StatusUnauthorized)
		return
	}

	// Отправка монет
	if err := api.service.SendCoin(r.Context(), sender, req.ToUser, req.Amount); err != nil {
		writeJSONError(w, "Ошибка при отправке монет: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
