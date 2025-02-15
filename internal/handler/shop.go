package handler

import (
	"Merch-Store/pkg/jwt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (api *API) BuyItem(w http.ResponseWriter, r *http.Request) {
	itemName := chi.URLParam(r, "item")

	if itemName == "" {
		writeJSONError(w, "Имя товара отсутствует в запросе", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(jwt.UserContextKey).(string)
	if !ok || user == "" {
		writeJSONError(w, "Не удалось определить отправителя", http.StatusUnauthorized)
		return
	}

	if err := api.service.BuyMerchandise(r.Context(), user, itemName, 1); err != nil {
		writeJSONError(w, "Ошибка при покупке товара: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
