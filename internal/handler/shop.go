package handler

import (
	"fmt"
	"net/http"

	"Merch-Store/pkg/jwt"

	"github.com/gorilla/mux"
)

// Новый обработчик покупки товара
func (api *API) BuyItem(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметр item из URL
	vars := mux.Vars(r)
	itemName := vars["item"]

	// Получаем имя отправителя из контекста
	user, ok := r.Context().Value(jwt.UserContextKey).(string)
	if !ok || user == "" {
		http.Error(w, `{"errors":"Не удалось определить отправителя"}`, http.StatusUnauthorized)
		return
	}

	// Вызов метода покупки товара с количеством 1
	err := api.db.BuyMerchandise(r.Context(), user, itemName, 1)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"errors":"Ошибка при покупке товара: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message":"Товар '%s' успешно куплен"}`, itemName)))
}
