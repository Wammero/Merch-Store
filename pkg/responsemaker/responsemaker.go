package responsemaker

import (
	"Merch-Store/internal/model"
	"encoding/json"
	"net/http"
)

func WriteJSONError(w http.ResponseWriter, message string, statusCode int) {
	WriteJSONResponse(w, model.ErrorResponse{Error: message}, statusCode)
}

// writeJSONResponse отправляет JSON-ответ с отступами
func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, `{
  "errors": "Ошибка при обработке данных"
}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
