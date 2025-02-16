package handler

import (
	"encoding/json"
	"net/http"

	"Merch-Store/pkg/jwt"
	"Merch-Store/pkg/responsemaker"
	"Merch-Store/pkg/validators"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int64  `json:"amount"`
}

func (api *API) SendCoin(w http.ResponseWriter, r *http.Request) {
	var req SendCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responsemaker.WriteJSONError(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if !validators.IsValidUsername(req.ToUser) || req.Amount <= 0 {
		responsemaker.WriteJSONError(w, "Имя пользователя и количество монет обязательны", http.StatusBadRequest)
		return
	}

	// Получаем имя отправителя из контекста
	sender, ok := r.Context().Value(jwt.UserContextKey).(string)
	if !ok || sender == "" {
		responsemaker.WriteJSONError(w, "Не удалось определить отправителя", http.StatusUnauthorized)
		return
	}

	// Отправка монет
	if err := api.service.SendCoin(r.Context(), sender, req.ToUser, req.Amount); err != nil {
		responsemaker.WriteJSONError(w, "Ошибка при отправке монет: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
