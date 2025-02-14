package handler

import (
	"Merch-Store/pkg/jwt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Новый обработчик покупки товара
func (api *API) BuyItem(w http.ResponseWriter, r *http.Request) {
	itemName := chi.URLParam(r, "item")
	log.Printf("Извлечён параметр item: '%s'", itemName)

	if itemName == "" {
		http.Error(w, `{"errors":"Имя товара отсутствует в запросе"}`, http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(jwt.UserContextKey).(string)
	if !ok || user == "" {
		http.Error(w, `{"errors":"Не удалось определить отправителя"}`, http.StatusUnauthorized)
		return
	}

	log.Printf("Покупка товара: user=%s, item=%s", user, itemName)

	if err := api.service.BuyMerchandise(r.Context(), user, itemName, 1); err != nil {
		log.Printf("Ошибка при покупке: %v", err)
		http.Error(w, fmt.Sprintf(`{"errors":"Ошибка при покупке товара: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": fmt.Sprintf("Товар '%s' успешно куплен", itemName)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
