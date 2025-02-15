package jwt

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Устанавливаем тестовый секрет
func init() {
	SetSecret("test_secret")
}

func TestGenerateJWT(t *testing.T) {
	username := "testuser"
	token, err := GenerateJWT(username)
	if err != nil {
		t.Fatalf("Ошибка генерации токена: %v", err)
	}

	if token == "" {
		t.Error("Ожидался непустой токен")
	}
}

func TestJWTValidator_ValidToken(t *testing.T) {
	token, _ := GenerateJWT("testuser")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := JWTValidator(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, но получен %d", rr.Code)
	}
}

func TestJWTValidator_InvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	rr := httptest.NewRecorder()
	handler := JWTValidator(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус 401, но получен %d", rr.Code)
	}
}

func TestJWTValidator_MissingToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	handler := JWTValidator(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус 401, но получен %d", rr.Code)
	}
}
