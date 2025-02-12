package handler

import (
	"encoding/json"
	"net/http"

	"Merch-Store/pkg/jwt"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int64  `json:"amount"`
}

func (api *API) SendCoin(w http.ResponseWriter, r *http.Request) {
	var req SendCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"errors":"Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	if req.ToUser == "" || req.Amount <= 0 {
		http.Error(w, `{"errors":"Имя пользователя и количество монет обязательны"}`, http.StatusBadRequest)
		return
	}

	// Получаем имя отправителя из контекста
	sender, ok := r.Context().Value(jwt.UserContextKey).(string)
	if !ok || sender == "" {
		http.Error(w, `{"errors":"Не удалось определить отправителя"}`, http.StatusUnauthorized)
		return
	}

	// Отправка монет
	if err := api.db.SendCoin(r.Context(), sender, req.ToUser, req.Amount); err != nil {
		http.Error(w, `{"errors":"Ошибка при отправке монет"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Монеты успешно отправлены"}`))
	w.Write([]byte(sender))
}
