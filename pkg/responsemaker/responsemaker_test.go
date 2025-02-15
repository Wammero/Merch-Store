package responsemaker

import (
	"Merch-Store/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSONError(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteJSONError(rr, "Ошибка сервера", http.StatusInternalServerError)

	// Проверяем статус код
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("ожидался статус %d, но получен %d", http.StatusInternalServerError, rr.Code)
	}

	// Проверяем заголовок Content-Type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидался заголовок Content-Type application/json, но получен %s", contentType)
	}

	// Проверяем тело ответа
	expectedResponse := `{
  "errors": "Ошибка сервера"
}`
	if strings.TrimSpace(rr.Body.String()) != expectedResponse {
		t.Errorf("ожидался ответ %s, но получен %s", expectedResponse, rr.Body.String())
	}
}

func TestWriteJSONResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	responseData := model.ErrorResponse{Error: "Другая ошибка"}

	WriteJSONResponse(rr, responseData, http.StatusBadRequest)

	// Проверяем статус код
	if rr.Code != http.StatusBadRequest {
		t.Errorf("ожидался статус %d, но получен %d", http.StatusBadRequest, rr.Code)
	}

	// Проверяем заголовок Content-Type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидался заголовок Content-Type application/json, но получен %s", contentType)
	}

	// Проверяем тело ответа
	expectedResponse := `{
  "errors": "Другая ошибка"
}`
	if strings.TrimSpace(rr.Body.String()) != expectedResponse {
		t.Errorf("ожидался ответ %s, но получен %s", expectedResponse, rr.Body.String())
	}
}
